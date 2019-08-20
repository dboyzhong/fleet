package service

import (
	"encoding/json"
	"time"
	kitlog "github.com/go-kit/kit/log"
	"github.com/kolide/fleet/server/kolide"
	"github.com/DeanThompson/jpush-api-go-client"
	"github.com/DeanThompson/jpush-api-go-client/push"
	"context"
	"net/http"
	"errors"
	"encoding/hex"
	"encoding/binary"
	"crypto/md5"
	"strconv"
)

type rulesEngine struct {
	ctx context.Context
	eventChan chan *kolide.Alarm
	msgChan chan []json.RawMessage
	startSeq int64
}

func newRulesEngine(ctx context.Context) *rulesEngine {
	ret := &rulesEngine {
		ctx: ctx,
		eventChan : make(chan *kolide.Alarm, 100),
		msgChan   : make(chan []json.RawMessage),
	}

	go func() {
		defer close(ret.eventChan)
		for {
			msg, err := ret.subscribeAlarm(ret.ctx)
			if(err != nil) {
				return
			}
			ret.eventChan <- msg
		}
	}()
	return ret;
}

func (re *rulesEngine) AlarmChannel() <-chan *kolide.Alarm {
	return re.eventChan
}

func (re *rulesEngine) sendEvent(msg []json.RawMessage) error {
	re.msgChan <- msg
	return nil
}

func (re *rulesEngine) subscribeAlarm(ctx context.Context) (*kolide.Alarm, error) {
	msg := <- re.msgChan

	a := &kolide.Alarm{
		Platform:"Alarm_ebi",
		Hostname:"virtual-1",
		Data: make([]*kolide.AlarmData, 0),
	}

	event := Event{}
	for _, v := range msg {

		err := json.Unmarshal(v, &event)
		if err != nil {
			continue
		}
	    a.Uid = event.Uid

	    re.startSeq += 1
	    h := md5.New()
	    h.Write(v)

	    var buf = make([]byte, 8)
	    binary.BigEndian.PutUint64(buf, uint64(re.startSeq))
	    h.Write(buf)
	    eventId := hex.EncodeToString(h.Sum(nil))

	    ad := &kolide.AlarmData{}
	    ad.Level = 0
	    ad.EventId = eventId 
	    ad.Title = "time changed"
	    ad.Type = 1
	    ad.CreateTime = time.Unix(event.Ts, 0)
	    ad.Ip = "192.168.0.2"
	    ad.IOC = "ioc"
	    ad.Details = "Alarm for test"
	    a.Data = append(a.Data, ad)
	}
	return  a, nil
}

// event middleware logs the service actions
type eventMiddleware struct {
	kolide.Service
	ds      kolide.Datastore
	logger  kitlog.Logger
	jclient *jpush.JPushClient
	pf      *push.Platform
	re      *rulesEngine
	startSeq int64
}

func (ew eventMiddleware) push(a *kolide.Alarm) error {

	msg, _ := json.Marshal(a)
	audience := push.NewAudience() 
	audience.SetAlias([]string{a.Uid})

	iosNotification := push.NewIosNotification("Ebi Alert")
    iosNotification.Badge = 1
    iosNotification.ContentAvailable = true
	iosNotification.AddExtra("Alarm", string(msg))

	notification := push.NewNotification("Ebi Alert")
    notification.Ios = iosNotification
	
	options := push.NewOptions()
    options.TimeToLive = 10000000
    options.ApnsProduction = true
	options.BigPushDuration = 1500
	
	payload := push.NewPushObject()
    payload.Platform = ew.pf
    payload.Audience = audience
    payload.Notification = notification
	payload.Options = options
	
	result, err := ew.jclient.PushValidate(payload)
    if err != nil {
		ew.logger.Log("err", "push failed: ", err)
	} else {
		if _, ok := result.MsgId.(int); !ok {
			ew.logger.Log("err", "push falied: ", result)
			err = errors.New(result.Error.String())
		} else {
			ew.logger.Log("err", "push success: ", result)
		}
	}
	return err
}

// NewLoggingService takes an existing service and adds a logging wrapper
func NewEventService(svc kolide.Service, ds kolide.Datastore, logger kitlog.Logger, jpushID, jpushKey string) kolide.Service {

	s := eventMiddleware{
		Service: svc,
		ds: ds,
		logger: logger,
		jclient: jpush.NewJPushClient(jpushID, jpushKey),
		pf: &push.Platform{},
		re: newRulesEngine(context.Background()),
		startSeq: time.Now().Unix(),
	}
	s.pf.Add("ios", "android")

	go func(){
		s.AlarmRoutine()
	}()
	return s
}

func decodeRiskMetricRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	uid := r.URL.Query().Get("uid")
	if(uid == "") {
		return nil, errors.New("no uid field found")
	}
	return riskMetricRequest{Uid: uid}, nil
}

func decodeSetEventStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req setEventStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return req, nil
}

func decodeEventHistoryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	uid   := r.URL.Query().Get("uid")
	sort  := r.URL.Query().Get("sort")
	var start, end int64
	var err error


	if start, err = strconv.ParseInt(r.URL.Query().Get("start"),10,64); err != nil {
		return nil, errors.New("param error")
	}

	if end, err = strconv.ParseInt(r.URL.Query().Get("end"),10,64); err != nil {
		return nil, errors.New("param error")
	}

	return eventHistoryRequest{Uid: uid, Sort: sort, Start: start, End: end}, nil
}

func (ew eventMiddleware) AlarmRoutine() {
	for {
		select {
		case msg, ok := <-ew.re.AlarmChannel():
			alarms, err := ew.getAlarm(0)
			if err == nil {
				for _, v := range alarms {
					ew.logger.Log("info", "repush failed event: ", v.Uid, v.EventId)
					if nil == ew.push(v) {
						ew.logger.Log("info", "update alarm with staus 1", v.Uid, v.EventId)
						ew.update(v, 1)
					}
				}
			}
			if ok {
				if(nil == ew.push(msg)) {
					ew.logger.Log("info", "save alarm with status 1", msg.Uid, msg.EventId)
					ew.save(msg, 1)
				} else {
					ew.logger.Log("info", "save alarm with status 0", msg.Uid, msg.EventId)
					ew.save(msg, 0)
				}
			} else {
				return
			}
		}
	}
}

func (ew eventMiddleware) update(a *kolide.Alarm, status int) error {
	for _, v := range a.Data {
		ew.ds.SetEventStatus(a.Uid, v.EventId, status)	
	}
	return nil
}

func (ew eventMiddleware) save(a *kolide.Alarm, status int) error {
	for _, v := range a.Data {
		AlarmString, _:= json.Marshal(v)
		ew.ds.NewEvent(a.Uid, v.EventId, a.Platform, a.Hostname, a.Content, string(AlarmString), status)	
	}
	return nil
}

func (ew eventMiddleware) getAlarm(status int) (map[string]*kolide.Alarm, error) {

	alarmMap := make(map[string]*kolide.Alarm)

	msg, err := ew.ds.GetAlarm(status)
	if err == nil {
		for _, v := range msg {
			if _, ok := alarmMap[v.Uid]; !ok {
				alarmMap[v.Uid] = &kolide.Alarm{
					Uid: v.Uid,
					Platform: v.Platform,
					Hostname: v.Hostname,
					Data: make([]*kolide.AlarmData, 0),
				}
			}

			alarmData := &kolide.AlarmData{}
			if err = json.Unmarshal([]byte(v.DataDB) , alarmData); err != nil {
				ew.logger.Log("err", "getAlarm json unmarshal error:%v:", err)
				return nil, err
			} else {
				alarmMap[v.Uid].Data = append(alarmMap[v.Uid].Data, alarmData)
			}
		}
	} else {
		ew.logger.Log("err", "getAlarm.ds.GetAlarm error:%v", err)
	}
	return alarmMap, nil
}


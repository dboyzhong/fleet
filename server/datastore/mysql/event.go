package mysql

import (
	_"database/sql"
	"time"
	"github.com/kolide/fleet/server/kolide"
	"fmt"
	_"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"encoding/json"
)

var (
	riskMetricDesc  = map[int]string {
		0 : "目前风险系统较优",
		1 : "风险指数较差，需要尽一步排查更改",
		2 : "安全威胁风险指数非常大，请及时排查安全隐患",
	}

	riskMetricScoreTotal = 100
	riskMetricScoreLimit = 10 

	riskMetricLevel2First = true
	riskMetricLevel2Minus = 10
	riskMetricLevel2MinusP= 20
	riskMetricLevel2Limit = 50 

	riskMetricLevel1Minus = 5
	riskMetricLevel1Limit = 30

	riskMetricLevel0Minus = 3
	riskMetricLevel0Limit = 20
)

func (d *Datastore) NewEvent(uid, eventId, platform, hostname string, content, alarm string, level, status int) (error) {
	sqlStatement := `
	INSERT INTO event (
		uid,
		event_id,
		platform,
		hostname,
		content,
		alarm,
		status,
		level
	)
	VALUES( ?,?,?,?,?,?,?,? )
	`
	_, err := d.db.Exec(sqlStatement, uid, eventId, platform, hostname, content, alarm, status, level);
	if err != nil {
		time.Sleep(time.Second)
		_, err = d.db.Exec(sqlStatement, uid, eventId, platform, hostname, content, alarm, status);
	}
	return err;
}

func (d* Datastore) GetRiskMetric(uid string) (*kolide.RiskMetric, error) {

	sqlStatement := `
		SELECT level FROM event 
		WHERE uid = ? and (status = 0 or status = 1)
	`
	var content []int
	err := d.db.Select(&content, sqlStatement, uid)
	if err != nil {
		fmt.Println("get risk metric error: ", err)
		time.Sleep(time.Second)
		err = d.db.Select(&content, sqlStatement, uid)
		if err != nil {
			return nil, errors.Wrap(err, "get risk metric")
		}
	}

	level0Minus := 0
	level1Minus := 0
	level2Minus := 0
	score := 0
	totalMinus := 0

	for _, level := range content {
		if level == 1 {
			level1Minus += riskMetricLevel1Minus
			if level1Minus > riskMetricLevel1Limit {
				level1Minus = riskMetricLevel1Limit
			}
		} else if level == 2 {
			riskMetricScoreTotal = 60
			if riskMetricLevel2First {
				level2Minus = (riskMetricLevel2Minus + riskMetricLevel2MinusP)
			} else {
				level2Minus += riskMetricLevel2Minus
			}
			if level2Minus > riskMetricLevel2Limit {
				level2Minus = riskMetricLevel2Limit
			}
		} else if level == 0 {
			level0Minus += riskMetricLevel0Minus
			if level0Minus > riskMetricLevel0Limit {
				level0Minus = riskMetricLevel0Limit
			}
		}
	}
	totalMinus = level0Minus + level1Minus + level2Minus
	if (totalMinus + riskMetricScoreLimit) > riskMetricScoreTotal {
		totalMinus = (riskMetricScoreTotal - riskMetricScoreLimit)
	}

	score = (riskMetricScoreTotal - totalMinus)
	var desc string
	if score >= 80 {
		desc = riskMetricDesc[0]
	} else if score >= 60 && score < 80 {
		desc = riskMetricDesc[1]
	} else {
		desc = riskMetricDesc[2]
	}
	
	return &kolide.RiskMetric{
		Uid: uid,
		Score: score,
		Desc : desc,
	}, nil
}

func (d* Datastore) SetEventStatus(uid, eventId string, status int) (string, error) {

	sqlStatement := `
	UPDATE event SET status=? WHERE event_id=?
	`
	res, err := d.db.Exec(sqlStatement, status, eventId);

	if err != nil {
		fmt.Println("set event error: ", err)
		time.Sleep(time.Second)
		res, err = d.db.Exec(sqlStatement, status, eventId);
		if err != nil {
			return "failed", err
		}
	}
	rowsaffected, err := res.RowsAffected()
   	if err != nil {
   	    return "failed", err
	} else if rowsaffected > 0 {
		return "success", nil
	} else {
		return "failed", errors.New("no such event_id")
	}
}

func (d* Datastore) GetAlarm(status int) ([]*kolide.Alarm, error) {

	sqlStatement := `
		SELECT uid, platform, hostname, event_id, alarm FROM event 
		WHERE status = ? LIMIT 2
	`
	var content []*kolide.Alarm
	err := d.db.Select(&content, sqlStatement, status)
	if err != nil {
		fmt.Println("get alarm error: ", err)
		time.Sleep(time.Second)
		err = d.db.Select(&content, sqlStatement, status)
		if err != nil {
			return nil, errors.Wrap(err, "get alarm")
		}
	}
	return content, nil
}

func (d* Datastore) EventHistory(uid, sort string, start, end, level, status int64) ([]*kolide.EventHistory, error) {

	var sqlStatement string
	var history []*kolide.EventHistory

	if sort == "dec" {
		sort = "desc"
	} else if sort == "inc" {
		sort = "asc"
	}

	sqlStatementFormat := `SELECT uid, platform, hostname, level, alarm, status FROM event 
					WHERE uid = ? %s order by id %s limit ?,?`

	args := make([]interface{}, 0, 10)
	args = append(args, uid)

	if 3 == level {
		if 1 == status {
			sqlStatement = fmt.Sprintf(sqlStatementFormat, "and (status = 0 or status = 1)", sort)
		} else if 2 == status {
			sqlStatement = fmt.Sprintf(sqlStatementFormat, "and (status = 2)", sort)
		} else {
			sqlStatement = fmt.Sprintf(sqlStatementFormat, "", sort)
		}
	} else {
		if 1== status {
			sqlStatement = fmt.Sprintf(sqlStatementFormat, "and level = ? and (status = 0 or status = 1)", sort)
			args = append(args, level)
		} else if 2 == status {
			sqlStatement = fmt.Sprintf(sqlStatementFormat, "and level = ? and status = 2", sort)
			args = append(args, level)
		} else {
			sqlStatement = fmt.Sprintf(sqlStatementFormat, "and level = ?", sort)
			args = append(args, level)
		}
	}

	args = append(args, start)
	args = append(args, end - start + 1)

	err := d.db.Select(&history, sqlStatement, args...)
	if err != nil {
		fmt.Println("get event history error: ", err)
		time.Sleep(time.Second)
		err = d.db.Select(&history, sqlStatement, args...)
		if err != nil {
			return nil, errors.Wrap(err, "event history")
		}
	}

	for _, v := range history {
		if err := json.Unmarshal([]byte(v.DataDB), v); err != nil {
			return nil, errors.Wrap(err, "event history json error")
		}
	}

	return history, nil
}

func (d* Datastore) EventDetails(uid, event_id string) (*kolide.EventDetails, error) {

	var content []*kolide.EventDetails
	sqlStatement := `
		SELECT e.uid, e.platform, e.hostname, e.event_id, e.level, e.alarm, e.status, IFNULL(i.ioc, '') AS ioc FROM event e LEFT JOIN ioc i 
		ON e.event_id = i.event_id WHERE e.uid = ? and e.event_id = ? LIMIT 1
	`

	err := d.db.Select(&content, sqlStatement, uid, event_id)
	if err != nil {
		fmt.Println("get details error: ", err)
		time.Sleep(time.Second)
		err = d.db.Select(&content, sqlStatement, uid, event_id)
		if err != nil {
			return nil, errors.Wrap(err, "get event details")
		}
	}

	if content == nil {
		return nil, fmt.Errorf("event not found")
	}

	ioc := content[0].IOC

	if err := json.Unmarshal([]byte(content[0].DataDB), content[0]); err != nil {
		fmt.Println("get details json error: ", err)
		return nil, errors.Wrap(err, "event details json error")
	}
	content[0].IOC = ioc
	return content[0], nil
}

func (d *Datastore) BannerInf(uid, host_uuid string) (*kolide.BannerInf, error) {
	return nil, nil
}
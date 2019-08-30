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
	return &kolide.RiskMetric{
		Uid: uid,
		Score: 80,
	}, nil
}

func (d* Datastore) SetEventStatus(uid, eventId string, status int) (string, error) {

	sqlStatement := `
	UPDATE event SET status=? WHERE event_id=?
	`
	res, err := d.db.Exec(sqlStatement, status, eventId);

	if err != nil {
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
		time.Sleep(time.Second)
		err = d.db.Select(&content, sqlStatement, status)
		if err != nil {
			return nil, errors.Wrap(err, "get alarm")
		}
	}
	return content, nil
}

func (d* Datastore) EventHistory(uid, sort string, start, end, level int64) ([]*kolide.EventHistory, error) {

	var sqlStatement string
	var history []*kolide.EventHistory

	if 3 == level {
		if sort == "desc" {
		    sqlStatement = `
		    	SELECT uid, platform, hostname, alarm FROM event 
		    	WHERE uid = ? order by id desc limit ?,?
		    `
		} else {
		    sqlStatement = `
		    	SELECT uid, platform, hostname, level, alarm, status FROM event 
				WHERE uid = ? limit ?,?
			`
		}
		err := d.db.Select(&history, sqlStatement, uid, start, end - start + 1)
		if err != nil {
			time.Sleep(time.Second)
			err = d.db.Select(&history, sqlStatement, uid, start, end - start + 1)
			if err != nil {
				return nil, errors.Wrap(err, "event history")
			}
		}
	} else {
		if sort == "desc" {
		    sqlStatement = `
		    	SELECT uid, platform, hostname, alarm FROM event 
		    	WHERE uid = ? and level = ? order by id desc limit ?,?
		    `
		} else {
		    sqlStatement = `
		    	SELECT uid, platform, hostname, level, alarm, status FROM event 
				WHERE uid = ? and level = ? limit ?,?
			`
		}
		err := d.db.Select(&history, sqlStatement, uid, level, start, end - start + 1)
		if err != nil {
			time.Sleep(time.Second)
			err = d.db.Select(&history, sqlStatement, uid, level, start, end - start + 1)
			if err != nil {
				return nil, errors.Wrap(err, "event history")
			}
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
		SELECT uid, platform, hostname, event_id, level, alarm, status FROM event 
		WHERE uid = ? and event_id = ? LIMIT 1
	`

	err := d.db.Select(&content, sqlStatement, uid, event_id)
	if err != nil {
		time.Sleep(time.Second)
		err = d.db.Select(&content, sqlStatement, uid, event_id)
		if err != nil {
			return nil, errors.Wrap(err, "get event details")
		}
	}

	if content == nil {
		return nil, fmt.Errorf("event not found")
	}

	if err := json.Unmarshal([]byte(content[0].DataDB), content[0]); err != nil {
		return nil, errors.Wrap(err, "event details json error")
	}
	return content[0], nil
}
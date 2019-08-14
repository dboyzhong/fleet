package mysql

import (
	_"database/sql"
	_"time"
	"github.com/kolide/fleet/server/kolide"
	_"fmt"
	_"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (d *Datastore) NewEvent(uid, eventId, platform, hostname string, content, alarm string, status int) (error) {
	sqlStatement := `
	INSERT INTO event (
		uid,
		event_id,
		platform,
		hostname,
		content,
		alarm,
		status
	)
	VALUES( ?,?,?,?,?,?,? )
	`
	_, err := d.db.Exec(sqlStatement, uid, eventId, platform, hostname, content, alarm, status);
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
		return "failed", err
	} else {
		rowsaffected, err := res.RowsAffected()
    	if err != nil {
    	    return "failed", err
		} else if rowsaffected > 0 {
			return "success", nil
		} else {
			return "failed", errors.New("no such event_id")
		}
	}
}

func (d* Datastore) GetAlarm(status int) ([]*kolide.Alarm, error) {

	sqlStatement := `
		SELECT uid, event_id, alarm FROM event 
		WHERE status = ? LIMIT 2
	`
	var content []*kolide.Alarm
	err := d.db.Select(&content, sqlStatement, status)
	if err != nil {
		return nil, errors.Wrap(err, "get alarm")
	}
	return content, nil
}
package kolide
import(
	"context"
	"time"
)	

type RiskMetric struct {
	Uid   string `json:"uid"`
	Score int    `json:"score"`
}

type AlarmData struct {
	Level int            `json:"level"`
	EventId string       `json:"event_id"`
	Title string         `json:"title"`
	Type int             `json:"type"`
	CreateTime time.Time `json:"create_time"`
	Ip string            `json:"ip"`
	IOC string           `json:"ioc"`
	Details string       `json:"details"`
}

type Alarm struct {
	Uid      string       `json:"uid" db:"uid"`
	Platform string       `json:"platform" db:"platform"`
	Hostname string       `json:"hostname" db:"hostname"`
	Data     []*AlarmData `json:"data" db:"-"`
	EventId  string       `json:"-" db:"event_id"`
	Content  string       `json:"-" db:"content"`
	DataDB   string       `json:"-" db:"alarm"`
}

type EventHistory struct {
	Uid      string       `json:"uid"      db:"uid"`
	Platform string       `json:"platform" db:"platform"`
	Hostname string       `json:"hostname" db:"hostname"`
	Level int             `json:"level"    db:"-"`
	EventId string        `json:"event_id" db:"-"`
	Title string          `json:"title"    db:"-"`
	Type int              `json:"type"     db:"-"`
	CreateTime time.Time  `json:"create_time" db:"-"`
	Ip string             `json:"ip"       db:"-"`
	IOC string            `json:"ioc"      db:"-"`
	Details string        `json:"details"  db:"-"`
	DataDB   string       `json:"-"        db:"alarm"`
}


type EventStore interface {
	NewEvent(uid, eventId, platform, hostname string, content, alarm string, status int) error
	GetRiskMetric(uid string) (*RiskMetric, error)
	SetEventStatus(uid, eventId string, status int) (string, error)
	GetAlarm(status int) ([]*Alarm, error)
	EventHistory(uid, sort string, start, end int64) ([]*EventHistory, error)
}

type EventService interface {
	GetRiskMetric(ctx context.Context, uid string) (*RiskMetric, error)
	SetEventStatus(ctx context.Context, uid, eventId string, status int) (string, error)
	EventHistory(ctx context.Context, uid, sort string, start, end int64) ([]*EventHistory, error)
}
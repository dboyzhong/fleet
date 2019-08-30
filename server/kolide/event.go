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
	RemoteIp string      `json:"remote_address"`
	AttackIp string      `json:"attack_address"`
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
	Level int             `json:"level"    db:"level"`
	EventId string        `json:"event_id" db:"event_id"`
	Title string          `json:"title"    db:"-"`
	Type int              `json:"type"     db:"-"`
	CreateTime time.Time  `json:"create_time" db:"-"`
	RemoteIp string       `json:"remote_address"   db:"-"`
	AttackIp string       `json:"attack_address"   db:"-"`
	IOC string            `json:"ioc"      db:"-"`
	Details string        `json:"details"  db:"-"`
	DataDB   string       `json:"-"        db:"alarm"`
	Status  int           `json:"status"   db:"status"`
}

type EventDetails EventHistory

type EventStore interface {
	NewEvent(uid, eventId, platform, hostname string, content, alarm string, level, status int) error
	GetRiskMetric(uid string) (*RiskMetric, error)
	SetEventStatus(uid, eventId string, status int) (string, error)
	GetAlarm(status int) ([]*Alarm, error)
	EventHistory(uid, sort string, start, end, level int64) ([]*EventHistory, error)
	EventDetails(uid, event_id string) (*EventDetails, error)
}

type EventService interface {
	GetRiskMetric(ctx context.Context, uid string) (*RiskMetric, error)
	SetEventStatus(ctx context.Context, uid, eventId string, status int) (string, error)
	EventHistory(ctx context.Context, uid, sort string, start, end, level int64) ([]*EventHistory, error)
	EventDetails(ctx context.Context, uid, event_id string) (*EventDetails, error)
}

package kolide
import(
	"context"
	"time"
)	

type RiskMetric struct {
	Uid   string `json:"uid"`
	Score int    `json:"score"`
	Desc  string `json:"description"`
}

type AlarmData struct {
	Level int            `json:"level"`
	EventId string       `json:"event_id"`
	Title string         `json:"title"`
	Type int             `json:"type"`
	CreateTime time.Time `json:"create_time"`
	RemoteIp string      `json:"remote_address"`
	AttackIp string      `json:"attack_address"`
	AttackRegion string  `json:"attack_region"`
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
	AttackRegion string   `json:"attack_region"    db:"-"`
	IOC string            `json:"ioc"      db:"ioc"`
	Details string        `json:"details"  db:"-"`
	DataDB   string       `json:"-"        db:"alarm"`
	Status  int           `json:"status"   db:"status"`
}

type SmtpConfig struct {
	Uid        string   `db:"uid"`
	ServerAddr string   `db:"smtp_server"`
	ServerPort int      `db:"smtp_server_port"`
	User       string   `db:"smtp_user"`
	Passwd     string   `db:"smtp_passwd"`
	Emails     []string
}

type BannerInf struct {
	Uid        string  `json:"uid"              db:"uid"`
	HostUUID   string  `json:"host_uuid"        db:"host_uuid"`
	Time       time.Time  `json:"time"          db:"time"`
	Address    string  `json:"address"         db:"address"`
	Service    string  `json:"service"         db:"service"`
	Status     string  `json:"state"           db:"state"`
	Banner     string  `json:"banner"          db:"banner"`
	Version    string  `json:"version"         db:"version"`
	ScriptRes  string  `json:"script_results"  db:"script_res"`
}

type PropertyCfg struct {
	Targets []string `json:"targets"       db:"targets"`
	Ports   string   `json:"port"          db:"ports"`
	Args    []string `json:"args"          db:"args"`
}

type PropertyResult struct {
	Code int `json:"code"`
}

type EventDetails EventHistory

type EventStore interface {
	NewEvent(uid, eventId, platform, hostname string, content, alarm string, level, status int) error
	GetRiskMetric(uid string) (*RiskMetric, error)
	SetEventStatus(uid, eventId string, status int) (string, error)
	GetAlarm(status int) ([]*Alarm, error)
	EventHistory(uid, sort string, start, end, level, status int64) ([]*EventHistory, error)
	EventDetails(uid, event_id string) (*EventDetails, error)
	BannerInf(uid, host_uuid string) (*BannerInf, error)
	PropertyCfg(uid string) (*PropertyCfg, error)
	PropertyResult(uid, host_uuid, result string, ts time.Time) (*PropertyResult, error)
	GetEventEmailCfg(uid string) (*SmtpConfig, error)
}

type EventService interface {
	GetRiskMetric(ctx context.Context, uid string) (*RiskMetric, error)
	SetEventStatus(ctx context.Context, uid, eventId string, status int) (string, error)
	EventHistory(ctx context.Context, uid, sort string, start, end, level, status int64) ([]*EventHistory, error)
	EventDetails(ctx context.Context, uid, event_id string) (*EventDetails, error)
	BannerInf(ctx context.Context, uid, host_uuid string) (*BannerInf, error)
	PropertyCfg(ctx context.Context, uid string) (*PropertyCfg, error)
	PropertyResult(ctx context.Context, uid, host_uuid, result string, ts time.Time) (*PropertyResult, error)
}

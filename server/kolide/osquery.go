package kolide

import (
	"context"
	"encoding/json"
)

type OsqueryService interface {
	EnrollAgent(ctx context.Context, enrollSecret, hostIdentifier string, hostDetails map[string](map[string]string)) (nodeKey string, err error)
	AuthenticateHost(ctx context.Context, nodeKey string) (host *Host, err error)
	GetClientConfig(ctx context.Context) (config map[string]interface{}, err error)
	// GetDistributedQueries retrieves the distributed queries to run for
	// the host in the provided context. These may be detail queries, label
	// queries, or user-initiated distributed queries. A map from query
	// name to query is returned. To enable the osquery "accelerated
	// checkins" feature, a positive integer (number of seconds to activate
	// for) should be returned. Returning 0 for this will not activate the
	// feature.
	GetDistributedQueries(ctx context.Context) (queries map[string]string, accelerate uint, err error)
	SubmitDistributedQueryResults(ctx context.Context, results OsqueryDistributedQueryResults, statuses map[string]OsqueryStatus) (err error)
	SubmitStatusLogs(ctx context.Context, logs []json.RawMessage) (err error)
	SubmitResultLogs(ctx context.Context, logs []json.RawMessage) (err error)
	SubmitResultCampaigns(ctx context.Context, logs []json.RawMessage) (err error)
}

// OsqueryDistributedQueryResults represents the format of the results of an
// osquery distributed query.
type OsqueryDistributedQueryResults map[string][]map[string]string

// OsqueryStatus represents osquery status codes (0 = success, nonzero =
// failure)
type OsqueryStatus int

const (
	// StatusOK is the success code returned by osquery
	StatusOK OsqueryStatus = 0
)

// QueryContent is the format of a query stanza in an osquery configuration.
type QueryContent struct {
	Query       string  `json:"query"`
	Description string  `json:"description,omitempty"`
	Interval    uint    `json:"interval"`
	Platform    *string `json:"platform,omitempty"`
	Version     *string `json:"version,omitempty"`
	Snapshot    *bool   `json:"snapshot,omitempty"`
	Removed     *bool   `json:"removed,omitempty"`
	Shard       *uint   `json:"shard,omitempty"`
}

type PermissiveQueryContent struct {
	QueryContent
	Interval interface{} `json:"interval"`
}

// Queries is a helper which represents the format of a set of queries in a pack.
type Queries map[string]QueryContent

type PermissiveQueries map[string]PermissiveQueryContent

// PackContent is the format of an osquery query pack.
type PackContent struct {
	Platform  string   `json:"platform,omitempty"`
	Version   string   `json:"version,omitempty"`
	Shard     uint     `json:"shard,omitempty"`
	Discovery []string `json:"discovery,omitempty"`
	Queries   Queries  `json:"queries"`
}

type PermissivePackContent struct {
	Platform  string            `json:"platform,omitempty"`
	Version   string            `json:"version,omitempty"`
	Shard     uint              `json:"shard,omitempty"`
	Discovery []string          `json:"discovery,omitempty"`
	Queries   PermissiveQueries `json:"queries"`
}

// Packs is a helper which represents the format of a list of osquery query packs.
type Packs map[string]PackContent

type PermissivePacks map[string]PermissivePackContent

// Decorators is the format of the decorator configuration in an osquery config.
type Decorators struct {
	Load     []string            `json:"load,omitempty"`
	Always   []string            `json:"always,omitempty"`
	Interval map[string][]string `json:"interval,omitempty"`
}

// OsqueryConfig is a struct that can be serialized into a valid osquery config
// using Go's JSON tooling.
type OsqueryConfig struct {
	Schedule   map[string]QueryContent `json:"schedule,omitempty"`
	Options    map[string]interface{}  `json:"options"`
	Decorators Decorators              `json:"decorators,omitempty"`
	Packs      Packs                   `json:"packs,omitempty"`
	// FilePaths contains named collections of file paths used for
	// FIM (File Integrity Monitoring)
	FilePaths FIMSections `json:"file_paths,omitempty"`
}

type PermissiveOsqueryConfig struct {
	OsqueryConfig
	Packs PermissivePacks `jsoon:"packs,omitempty"`
}

package service

import (
	"context"
	"encoding/json"
	_"encoding/hex"
	_"encoding/binary"
	_"time"
	_"crypto/md5"

	"github.com/kolide/fleet/server/kolide"
)

type Event struct {
	Uid string `json:"uid"`
	Ts  int64  `json:"unixTime"`
	Content json.RawMessage
}

func (ew eventMiddleware) EnrollAgent(ctx context.Context, enrollSecret string, hostIdentifier string, hostDetails map[string](map[string]string)) (string, error) {
	nodeKey, err := ew.Service.EnrollAgent(ctx, enrollSecret, hostIdentifier, hostDetails)
	return nodeKey, err
}

func (ew eventMiddleware) AuthenticateHost(ctx context.Context, nodeKey string) (*kolide.Host, error) {
	host, err := ew.Service.AuthenticateHost(ctx, nodeKey)
	return host, err
}

func (ew eventMiddleware) GetClientConfig(ctx context.Context) (map[string]interface{}, error) {
	config, err := ew.Service.GetClientConfig(ctx)
	return config, err
}

func (ew eventMiddleware) GetDistributedQueries(ctx context.Context) (map[string]string, uint, error) {
	queries, accelerate, err := ew.Service.GetDistributedQueries(ctx)
	return queries, accelerate, err
}

func (ew eventMiddleware) SubmitDistributedQueryResults(ctx context.Context, results kolide.OsqueryDistributedQueryResults, statuses map[string]kolide.OsqueryStatus) error {
	err := ew.Service.SubmitDistributedQueryResults(ctx, results, statuses)
	return err
}

func (ew eventMiddleware) SubmitStatusLogs(ctx context.Context, logs []json.RawMessage) error {
	err := ew.Service.SubmitStatusLogs(ctx, logs)
	return err
}

func (ew eventMiddleware) SubmitResultLogs(ctx context.Context, logs []json.RawMessage) error {
	var (
		err error
	)
	ew.re.sendEvent(logs)
	err = ew.Service.SubmitResultLogs(ctx, logs)
	return err
}

func (ew eventMiddleware) SubmitResultCampaigns(ctx context.Context, logs []json.RawMessage) error {
	var (
		err error
	)

	//TODO: real event handle
	err = ew.Service.SubmitResultCampaigns(ctx, logs)
	return err
}

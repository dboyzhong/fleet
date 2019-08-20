package service

import (
	"context"
	"github.com/kolide/fleet/server/kolide"
)
func (svc service) GetRiskMetric(ctx context.Context, uid string) (*kolide.RiskMetric, error) {
	return svc.ds.GetRiskMetric(uid)
}

func (svc service) SetEventStatus(ctx context.Context, uid, eventId string, status int) (string, error) {
	return svc.ds.SetEventStatus(uid, eventId, status)
}

func (svc service) EventHistory(ctx context.Context, uid, sort string, start, end int64) ([]*kolide.EventHistory, error) {
	return svc.ds.EventHistory(uid, sort, start, end)
}
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

func (svc service) EventHistory(ctx context.Context, uid, sort string, start, end, level, status int64) ([]*kolide.EventHistory, error) {
	return svc.ds.EventHistory(uid, sort, start, end, level, status)
}

func (svc service) EventDetails(ctx context.Context, uid, event_id string) (*kolide.EventDetails, error) {
	return svc.ds.EventDetails(uid, event_id)
}

func (svc service) BannerInf(ctx context.Context, uid, host_uuid string) (*kolide.BannerInf, error) {
	return svc.ds.BannerInf(uid, host_uuid)
}
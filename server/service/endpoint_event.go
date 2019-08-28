package service

import (
	"context"
	_"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
)

type riskMetricResponse struct {
	Metric      *kolide.RiskMetric `json:"metric"`
	Err         error              `json:"error,omitempty"`
}

type riskMetricRequest struct {
	Uid string	`json:"uid"`
}

type setEventStatusRequest struct {
	Uid     string	`json:"uid"`
	EventId string  `json:"event_id"`
	Status  int     `json:"status"`
}

type setEventStatusResponse struct {
	Result string `json:"result"`
	Err    error  `json:"error,omitempty"`
}

type eventHistoryRequest struct {
	Uid   string `json:"uid"`
	Sort  string `json:"sort"`
	Start int64  `json:"start"`
	End   int64  `json:"end"` 
	Level int64  `json:"level"`
}

type eventHistoryResponse struct {
	History []*kolide.EventHistory   `json:"event_history,omitempty"`
	Err         error              `json:"error,omitempty"`
}

func (r riskMetricResponse) error() error { return r.Err }

func makeRiskMetricEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(riskMetricRequest)
		metric, err := svc.GetRiskMetric(ctx, req.Uid)
		if err != nil {
			return riskMetricResponse{Err: err}, nil
		}
		return riskMetricResponse{Metric:metric, Err:nil}, nil
	}
}

func makeSetEventStatusEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(setEventStatusRequest)
		result, err := svc.SetEventStatus(ctx, req.Uid, req.EventId, req.Status)
		if err != nil {
			return setEventStatusResponse{Err: err}, nil
		}
		return setEventStatusResponse{Result: result, Err:nil}, nil
	}
}

func makeEventHistoryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(eventHistoryRequest)
		result, err := svc.EventHistory(ctx, req.Uid, req.Sort, req.Start, req.End,req.Level)
		if err != nil {
			return eventHistoryResponse{Err: err}, nil
		}
		return eventHistoryResponse{History: result, Err:nil}, nil
	}
}

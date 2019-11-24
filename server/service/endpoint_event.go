package service

import (
	"context"
	"time"
	"encoding/json"
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
	Status int64 `json:"status"`
}

type eventHistoryResponse struct {
	History []*kolide.EventHistory   `json:"event_history,omitempty"`
	Err         error              `json:"error,omitempty"`
}

type eventDetailsRequest struct {
	Uid string     `json:"uid"`
	EventId string `json:"event_id"`
}

type eventDetailsResponse struct {
	EventDetails *kolide.EventDetails `json:"event_details,omitempty"`
	Err         error                 `json:"error,omitempty"`
}

type eventBannerInfRequest struct {
	Uid string      `json:"uid"`
	HostUUID string `json:"host_uuid"`
}

type eventBannerInfResponse struct {
	BannerInf *kolide.BannerInf `json:"banner_inf ,omitempty"`
	Err    error    `json:"error,omitempty"`
}

type eventPropertyCfgRequest struct {
	Uid string      `json:"uid"`
}

type eventPropertyCfgResponse struct {
	PropertyCfg *kolide.PropertyCfg `json:"config,omitempty"`
	Err    error    `json:"error,omitempty"`
}

type eventPropertyResultRequest struct {
	Uid string      `json:"uid"`
	HostUUID string `json:"host_uuid"`
	Results  json.RawMessage `json:"results"`
	Ts   time.Time  `json:"ts"`
}

type eventPropertyResultResponse struct {
	PropertyResult *kolide.PropertyResult `json:"result,omitempty"`
	Err    error    `json:"error,omitempty"`
}

func (r eventPropertyCfgResponse) error() error { return r.Err }

func (r eventBannerInfResponse) error() error { return r.Err }

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
		result, err := svc.EventHistory(ctx, req.Uid, req.Sort, req.Start, req.End,req.Level, req.Status)
		if err != nil {
			return eventHistoryResponse{Err: err}, nil
		}
		return eventHistoryResponse{History: result, Err:nil}, nil
	}
}

func makeEventDetailsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(eventDetailsRequest)
		result, err := svc.EventDetails(ctx, req.Uid, req.EventId)
		if err != nil {
			return eventDetailsResponse{Err: err}, nil
		}
		return eventDetailsResponse{EventDetails: result, Err:nil}, nil
	}
}

func makeEventBannerInfEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(eventBannerInfRequest)
		result, err := svc.BannerInf(ctx, req.Uid, req.HostUUID)
		if err != nil {
			return eventBannerInfResponse{Err: err}, nil
		}
		return eventBannerInfResponse{BannerInf: result, Err:nil}, nil
	}
}

func makePropertyCfgEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(eventPropertyCfgRequest)
		result, err := svc.PropertyCfg(ctx, req.Uid)
		if err != nil {
			return eventPropertyCfgResponse{Err: err}, nil
		}
		return eventPropertyCfgResponse{PropertyCfg: result, Err:nil}, nil
	}
}

func makePropertyResultEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(eventPropertyResultRequest)
		result, err := svc.PropertyResult(ctx, req.Uid, req.HostUUID, string(req.Results), req.Ts)
		if err != nil {
			return eventPropertyResultResponse{Err: err}, nil
		}
		return eventPropertyResultResponse{PropertyResult: result, Err:nil}, nil
	}
}
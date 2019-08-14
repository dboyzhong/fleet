package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
)

type hostEbiResponse struct {
	kolide.Host
	Status      string `json:"status"`
	DisplayText string `json:"display_text"`
}

func hostEbiResponseForHost(ctx context.Context, svc kolide.Service, host *kolide.Host) (*hostEbiResponse, error) {
	return &hostEbiResponse{
		Host:        *host,
		Status:      host.Status(time.Now()),
		DisplayText: host.HostName,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////
// List Hosts
////////////////////////////////////////////////////////////////////////////////
type listHostsEbiRequest struct {
	Uid string	
}

type listHostsEbiResponse struct {
	Hosts []hostEbiResponse `json:"hosts"`
	Err   error          `json:"error,omitempty"`
}

func (r listHostsEbiResponse) error() error { return r.Err }

func makeListHostsEbiEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listHostsEbiRequest)
		hosts, err := svc.ListEbiHosts(ctx, req.Uid)
		if err != nil {
			return listHostsEbiResponse{Err: err}, nil
		}

		hostEbiResponses := make([]hostEbiResponse, len(hosts))
		for i, host := range hosts {
			h, err := hostEbiResponseForHost(ctx, svc, host)
			if err != nil {
				return listHostsEbiResponse{Err: err}, nil
			}

			hostEbiResponses[i] = *h
		}
		return listHostsEbiResponse{Hosts: hostEbiResponses}, nil
	}
}

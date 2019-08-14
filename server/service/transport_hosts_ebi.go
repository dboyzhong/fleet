package service

import (
	"context"
	"net/http"
	"errors"
)

func decodeListHostsEbiRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	uid := r.URL.Query().Get("uid")
	if(uid == "") {
		return nil, errors.New("no uid field found")
	}
	return listHostsEbiRequest{Uid: uid}, nil
}

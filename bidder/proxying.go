package main

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
)

func proxyingMiddleware(ctx context.Context, instances string) ServiceMiddleware {
	if instances == "" {
		return func(next BidderService) BidderService { return next }
	}

	var (
		//qps         = 100
		maxAttempts = 1
		maxTime     = 200 * time.Millisecond
	)

	var (
		instanceList = split(instances)
		endpointer   sd.FixedEndpointer
	)
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeBidProxy(ctx, instance)
		endpointer = append(endpointer, e)
	}

	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	return func(next BidderService) BidderService {
		return proxymw{ctx, next, retry}
	}
}

type proxymw struct {
	ctx  context.Context
	next BidderService
	bid  endpoint.Endpoint
}

func (mw proxymw) Bid(adReq adRequest) (adObject, error) {
	response, err := mw.bid(mw.ctx, adReq)
	if err != nil {
		return adObject{}, err
	}

	resp := response.(bidResponse)

	if resp.Err != "" {
		return resp.Bid, errors.New(resp.Err)
	}
	return resp.Bid, nil

}

func makeBidProxy(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}

	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}

	if u.Path == "" {
		u.Path = "/bid"
	}

	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeBidResponse,
	).Endpoint()
}

func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}

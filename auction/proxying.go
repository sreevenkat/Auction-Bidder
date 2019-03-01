package main

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sony/gobreaker"
)

func proxyingMiddleware(ctx context.Context, instances string) ServiceMiddleware {
	if instances == "" {
		return func(next AuctionService) AuctionService { return next }
	}

	var (
		maxAttempts = 1
		maxTime     = 300 * time.Millisecond
	)

	var (
		instanceList = split(instances)
		endpointer   sd.FixedEndpointer
	)
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeAuctionProxy(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		endpointer = append(endpointer, e)
	}

	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	return func(next AuctionService) AuctionService {
		return proxymw{ctx, next, retry}
	}
}

type proxymw struct {
	ctx     context.Context
	next    AuctionService
	auction endpoint.Endpoint
}

func (mw proxymw) Auction(adReq adRequest) (adObject, error) {
	response, err := mw.auction(mw.ctx, adReq)
	if err != nil {
		return adObject{}, err
	}

	resp := response.(auctionResponse)

	if resp.Err != "" {
		return resp.Bid, errors.New(resp.Err)
	}
	return resp.Bid, nil

}

func makeAuctionProxy(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}

	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}

	if u.Path == "" {
		u.Path = "/auction"
	}

	return httptransport.NewClient(
		"GET",
		u,
		encodeAuctionRequest,
		decodeAuctionResponse,
	).Endpoint()
}

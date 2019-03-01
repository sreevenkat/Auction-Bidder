package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type auctionResponse struct {
	Bid adObject `json:"bid"`
	Err string   `json:"err,omitempty"`
}

func makeAuctionEndpoint(as AuctionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(adRequest)
		adObj, err := as.Auction(req)
		if err != nil {
			return auctionResponse{Bid: adObj, Err: err.Error()}, nil

		}
		return auctionResponse{Bid: adObj, Err: ""}, nil

	}
}
func decodeAuctionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request adRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeAuctionResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	// Checking to see if the adspace got a bid and returning 200 or 204 appropriately
	if response.(auctionResponse).Bid.AdID != "" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeAuctionRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeAuctionResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response auctionResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

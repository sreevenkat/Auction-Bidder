package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type bidResponse struct {
	Bid adObject `json:"bid"`
	Err string   `json:"err,omitempty"`
}

func makeBidEndpoint(bs BidderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(adRequest)
		adObj, err := bs.Bid(req)

		if err != nil {
			return bidResponse{Bid: adObj, Err: err.Error()}, nil
		}
		return bidResponse{Bid: adObj, Err: ""}, nil
	}
}

func decodeBidRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request adRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	// Checking to see if bidder placed a bid or not and returning a 200 or 204 appropriately
	if response.(bidResponse).Bid.BidPlaced {
		w.WriteHeader(http.StatusOK)
	}else{
		w.WriteHeader(http.StatusNoContent)
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeBidResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response bidResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

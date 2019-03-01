package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"
)

type AuctionService interface {
	Auction(adRequest) (adObject, error)
}

type adObject struct {
	AdID          string  `json:"adID"`
	Price         float64 `json:"price"`
	AdPlacementID string  `json:"adPlacementID"`
	Name          string  `json:"name"`
}

type bid struct {
	Bid adObject `json:"bid"`
}

type adRequest struct {
	AdPlacementID string `json:"adPlacementID"`
}

type auctionService struct{}

func GetHighestBidder(bids []bid) adObject {
	var adObj adObject

	if len(bids) > 0 {

		sort.Slice(bids, func(i, j int) bool {
			return bids[i].Bid.Price > bids[j].Bid.Price
		})
		adObj = bids[0].Bid
	}

	return adObj
}

func makeRequestTobidder(adrequest adRequest, url string, c chan bid){

		payloadBuffer := new(bytes.Buffer)
		json.NewEncoder(payloadBuffer).Encode(adrequest)

		res, errReq := http.Post(url, "application/json; charset=utf-8", payloadBuffer)

		if errReq != nil {
			log.Fatal(errReq)
		}

		defer res.Body.Close()

		var b bid

		if res.StatusCode == http.StatusNoContent {
			c <- b

			return
		}

		err := json.NewDecoder(res.Body).Decode(&b)
		if err != nil {
			panic(err)
		}
		c <- b
}

func (auctionService) Auction(adrequest adRequest) (adObject, error) {
	if adrequest.AdPlacementID == "" {
		emptyAdObject := adObject{AdPlacementID: adrequest.AdPlacementID}
		return emptyAdObject, ErrEmpty
	}

	// Channel for returning bids
	c := make(chan bid, len(bidders))

	for i := 0; i < len(bidders); i++ {
		go makeRequestTobidder(adrequest ,bidders[i] + "/bid", c)
	}

	timeout := time.Duration(180) * time.Millisecond

	var bids []bid

	bidderLoop:
		for {
			select {
			// Waiting for bids to return
			case b, ok := <-c:
				if !ok {
					break
				} else {
					bids = append(bids, b)
				}

				if len(bidders) == len(bids) {
					break bidderLoop
				}
			// Timeout when bidder takes too long
			case <-time.After(timeout):
				fmt.Printf("Timed out waiting for bidder \n")
				break bidderLoop
			}
		}

	newadObject := GetHighestBidder(bids)
	return newadObject, nil
}

var ErrEmpty = errors.New("Empty AdPlacementId")

type ServiceMiddleware func(AuctionService) AuctionService

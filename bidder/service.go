package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/satori/go.uuid"
)

type BidderService interface {
	Bid(adRequest) (adObject, error)
}

type adObject struct {
	AdID          string  `json:"adID"`
	Price         float64 `json:"price"`
	AdPlacementID string  `json:"adPlacementID"`
	BidPlaced     bool    `json:"-"`
	Name          string  `json:"name"`
}

type adRequest struct {
	AdPlacementID string `json:"adPlacementID"`
}

type bidderService struct{}

func (bidderService) Bid(adrequest adRequest) (adObject, error) {
	if adrequest.AdPlacementID == "" {
		emptyAdObject := adObject{AdPlacementID: adrequest.AdPlacementID}
		return emptyAdObject, ErrEmpty
	}
	newadObject := GetBidOnAd(adrequest.AdPlacementID)
	return newadObject, nil
}

func GetBidOnAd(adPlacementID string) adObject {
	sendValue := rand.Intn(11)
	adID, _ := uuid.NewV4()
	newAdObject := adObject{AdPlacementID: adPlacementID, Name: *bidName}
	if sendValue <= 10 && *bidBehaviour != "nobuy" {
		newAdObject.AdID = adID.String()
		newAdObject.Price = rand.Float64()
	}
	if *bidBehaviour == "sleep" || sendValue == 7{
		time.Sleep(200 * time.Millisecond)
	}
	newAdObject.BidPlaced = newAdObject.Price != 0

	return newAdObject
}

var ErrEmpty = errors.New("Empty AdPlacementId")

type ServiceMiddleware func(BidderService) BidderService

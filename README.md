# Auction-Bidder

There are three docker-compose files:

## Normal Auction Bidder

  Command: `docker-compose up`
  File: docker-compose.yml
  This compose file demonstrates a normal usecase where any of the bidders may take more than 200ms to respond or may not buy   the adspace randomly. 
  
  It Consists of:
  
   1. **auctionproxy**: This is a proxy service that takes a HTTP request at port 8080 and forwards it to two
                          auction services using round-robin method
                          
   2. **auction1** & **auction2**: This is an auction service that takes a HTTP request forwarded from auctionproxy and 
                                    gets bids for the request from all the bidders simultaneously. Returns an AdObject
                                    for the highest bidding value or returns a 204 if there is no bid for the ad space or the request times out.
                                    
  3. **bid1**, **bid2** .. **bid5**: These are the bidder services which either bids for the adspace or skips bidding and returns either an AdObject or a 204 respectively.
  
  Curl Request to be sent to **auctionproxy**: 
  ```
  curl -X POST \
  http://localhost:8080/auction \
  -H 'Content-Type: application/json' \
  -d '{
    "adPlacementId": "xxxx-xxxx"
}'
  ```
  Response:
  ``` javascript
  { "bid":
    {
    "adID":"ba665727-fd38-4a1b-b5c5-1683b72a0a97",
    "price":0.9405090880450124,
    "adPlacementID":"xxxx-xxxx",
    "name":"bid1"
    }
  }
  ```
                               

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
  
  ## Slow Bidder
  Command: `docker-compose -f slow-bidder.yml up`
  
  File: [slow-bidder.yml](/slow-bidder.yml)
  
  This is the same as the normal auction bidder except that the bidder service **bid1** is forced to take more than 200ms by passing the `--behaviour` option with `sleep` as its value.
  
  ## NoBuy Bidder
  Command: `docker-compose -f nobuy-bidder.yml up`
  
  File: [nobuy-bidder.yml](/nobuy-bidder.yml)
  
  This is the same as the normal auction bidder except that the bidder service **bid1** is forced to not buy any ad space by passing the `--behaviour` option with `nobuy` as its value.
  
  
### What went well

- Once I got to know go-kit it was fairly easy to get around the application logic 
- Proxy service made it convenient to spin up additional auction services

### What went wrong

- Couldn't write tests since familiarising with go-kit took some time. Did tests manually

### What I'd do for the next iteration

- Explore how to improve concurrency further
- Would setup auto scaling
- Add service discovery such that if a new bidder went up the auction services would be made aware of

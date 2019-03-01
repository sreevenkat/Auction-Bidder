package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

var bidName, bidBehaviour *string

func main() {

	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy bid requests")
		name   = flag.String("name", "", "Name of the bidder")
		seed   = flag.Int64("seed", 0, "Random seed")
		behaviour = flag.String("behaviour", "", "Optional string to simulate sleep/nobuy scenario. Options are 'sleep' and 'nobuy'.")
	)

	flag.Parse()

	rand.Seed(*seed)
	bidName = name
	bidBehaviour = behaviour

	var bs BidderService
	bs = bidderService{}
	bs = proxyingMiddleware(context.Background(), *proxy)(bs)

	BidHandler := httptransport.NewServer(
		makeBidEndpoint(bs),
		decodeBidRequest,
		encodeResponse,
	)

	http.Handle("/bid", BidHandler)
	log.Println("msg", "HTTP", "addr", *listen, "name:", *name, "bidName", *bidName)
	log.Println("err", http.ListenAndServe(*listen, nil))
	log.Fatal(http.ListenAndServe(":8080", nil))
}


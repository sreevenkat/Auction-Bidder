package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
)

var bidders []string

func main() {

	var (
		listen      = flag.String("listen", ":8080", "HTTP listen address")
		biddersFlag = flag.String("biddersFlag", "", "list of bidder ports")
		proxy       = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy auction requests")
	)

	flag.Parse()

	bidders = split(*biddersFlag)

	var as AuctionService
	as = auctionService{}
	as = proxyingMiddleware(context.Background(), *proxy)(as)

	AuctionHandler := httptransport.NewServer(
		makeAuctionEndpoint(as),
		decodeAuctionRequest,
		encodeAuctionResponse,
	)

	http.Handle("/auction", AuctionHandler)
	log.Println("msg", "HTTP", "addr", *listen)
	log.Println("err", http.ListenAndServe(*listen, nil))
}

func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}

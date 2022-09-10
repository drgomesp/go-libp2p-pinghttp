package main

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peerstore"

	libp2pping "github.com/drgomesp/go-libp2p-pinghttp"
)

func main() {
	h1, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}
	h2, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}

	h1.Peerstore().AddAddrs(h2.ID(), h2.Addrs(), peerstore.PermanentAddrTTL)

	ctx := context.Background()
	svc, err := libp2pping.NewHttpPingService(ctx, h1, libp2pping.WithHttpAddr(":4000"))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		_ = svc.ListenAndServe(ctx)
	}()

	log.Println(fmt.Sprintf("visit: http://localhost:4000/v1/ping?peerId=%s", h2.ID().String()))
	<-ctx.Done()
}
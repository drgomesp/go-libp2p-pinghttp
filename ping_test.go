package libp2pping_test

import (
	"context"
	"crypto/rand"
	"io"
	mrand "math/rand"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
)

func newHost(t *testing.T, addrStr string, randseed int64, opts ...libp2p.Option) host.Host {
	ma, err := multiaddr.NewMultiaddr(addrStr)
	assert.NoError(t, err)

	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	assert.NoError(t, err)

	h, err := libp2p.New(
		append(opts,
			libp2p.ListenAddrs(ma),
			libp2p.Identity(priv),
			libp2p.DisableRelay(),
		)...,
	)
	assert.NoError(t, err)

	return h
}

func TestPing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ha := newHost(t, "/ip4/127.0.0.1/tcp/10000", 1)
	defer ha.Close()

	hb := newHost(t, "/ip4/127.0.0.1/tcp/10001", 2)
	defer hb.Close()

	ha.Peerstore().AddAddrs(hb.ID(), hb.Addrs(), peerstore.PermanentAddrTTL)

	pingService := ping.NewPingService(ha)

	res := <-pingService.Ping(ctx, hb.ID())
	spew.Dump(res)
}

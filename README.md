# go-libp2p-pinghttp

[![madeby](https://img.shields.io/badge/made%20by-%40drgomesp-blue)](https://github.com/drgomesp/)
[![Go Report Card](https://goreportcard.com/badge/github.com/drgomesp/go-libp2p-pinghttp)](https://goreportcard.com/report/github.com/drgomesp/go-libp2p-pinghttp)
[![build](https://github.com/drgomesp/go-libp2p-pinghttp/actions/workflows/go-test.yml/badge.svg?style=squared)](https://github.com/drgomesp/go-libp2p-pinghttp/actions)
[![codecov](https://codecov.io/gh/drgomesp/go-libp2p-pinghttp/branch/main/graph/badge.svg?token=BRMFJRJV2X)](https://codecov.io/gh/drgomesp/go-libp2p-pinghttp)

> Expose a Libp2p host's Ping service through HTTP.

## Table of Contents

- [Install](#install)
- [Features](#features)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Install

```bash
go get github.com/drgomesp/go-libp2p-pinghttp
```
g
## Features



## Usage

Initialize the http ping service with the libp2p host and start listening:

```go
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
    _ = svc.ListenAndServe()
}()

log.Println(fmt.Sprintf("visit: http://localhost:4000/v1/ping?peerId=%s", h2.ID().String()))
<-ctx.Done()
```

Visiting the suggested URL will give you a response like this: 

```json
$ curl http://localhost:4000/v1/ping\?peerId\=12D3KooWNNU1i2FqRbekHLMFssUJwAxmDyJvKs4D7VXU8F3rBpYq | jq 
{
  "duration": "231.905µs",
  "error": ""
}

```

## Contributing

PRs accepted.

## License

MIT © [Daniel Ribeiro](https://github.com/drgomesp)
package libp2pping

import (
	"context"
	"net/http"

	libp2pgrpc "github.com/drgomesp/go-libp2p-grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	libp2p "github.com/libp2p/go-libp2p/p2p/protocol/ping"

	v1 "github.com/drgomesp/go-libp2p-pinghttp/proto/v1"
)

type PingHttpService struct {
	v1.UnimplementedPingServiceServer

	host host.Host

	httpAddr string
	mux      *runtime.ServeMux
	grpcsrv  *libp2pgrpc.Server
	pingsvc  *libp2p.PingService
}

type ServiceOption func(service *PingHttpService)

func WithHttpAddr(httpAddr string) ServiceOption {
	return func(g *PingHttpService) {
		g.httpAddr = httpAddr
	}
}

func WithServeMux(mux *runtime.ServeMux) ServiceOption {
	return func(s *PingHttpService) {
		s.mux = mux
	}
}

func NewPingHttpService(ctx context.Context, h host.Host, opts ...ServiceOption) (*PingHttpService, error) {
	svc := &PingHttpService{host: h}

	grpcServer, err := libp2pgrpc.NewGrpcServer(ctx, svc.host)
	svc.grpcsrv = grpcServer
	svc.pingsvc = libp2p.NewPingService(h)

	for _, opt := range opts {
		opt(svc)
	}

	v1.RegisterPingServiceServer(svc.grpcsrv, svc)
	err = v1.RegisterPingServiceHandlerServer(ctx, svc.mux, svc)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *PingHttpService) Ping(ctx context.Context, req *v1.PingRequest) (*v1.PingResponse, error) {
	peerId, err := peer.Decode(req.PeerId)
	if err != nil {
		return &v1.PingResponse{Error: err.Error()}, nil
	}

	res := <-s.pingsvc.Ping(ctx, peerId)

	if res.Error != nil {
		return &v1.PingResponse{Error: res.Error.Error()}, nil
	}

	return &v1.PingResponse{Duration: res.RTT.String()}, nil
}

func (s *PingHttpService) ListenAndServe() error {
	return http.ListenAndServe(s.httpAddr, s.mux)
}

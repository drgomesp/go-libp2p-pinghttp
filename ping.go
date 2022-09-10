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

type HttpPingService struct {
	v1.UnimplementedPingServiceServer

	host     host.Host
	mux      *runtime.ServeMux
	grpcsrv  *libp2pgrpc.Server
	pingsvc  *libp2p.PingService
	httpAddr string
}

type ServiceOption func(service *HttpPingService)

func WithHttpAddr(httpAddr string) ServiceOption {
	return func(g *HttpPingService) {
		g.httpAddr = httpAddr
	}
}

func WithServeMux(mux *runtime.ServeMux) ServiceOption {
	return func(s *HttpPingService) {
		s.mux = mux
	}
}

func NewHttpPingService(ctx context.Context, h host.Host, opts ...ServiceOption) (*HttpPingService, error) {
	s := &HttpPingService{host: h, mux: runtime.NewServeMux()}
	for _, opt := range opts {
		opt(s)
	}
	grpcServer, err := libp2pgrpc.NewGrpcServer(ctx, s.host)
	s.grpcsrv = grpcServer
	s.pingsvc = libp2p.NewPingService(s.host)

	v1.RegisterPingServiceServer(s.grpcsrv, s)
	err = v1.RegisterPingServiceHandlerServer(ctx, s.mux, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *HttpPingService) Ping(ctx context.Context, req *v1.PingRequest) (*v1.PingResponse, error) {
	peerId, err := peer.Decode(req.PeerId)
	if err != nil {
		return &v1.PingResponse{Error: err.Error()}, err
	}

	res := <-s.pingsvc.Ping(ctx, peerId)

	if res.Error != nil {
		return &v1.PingResponse{Error: res.Error.Error()}, res.Error
	}

	return &v1.PingResponse{Duration: res.RTT.String()}, nil
}

func (s *HttpPingService) ListenAndServe(ctx context.Context) error {
	return http.ListenAndServe(s.httpAddr, s.mux)
}
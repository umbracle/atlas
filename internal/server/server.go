package server

import (
	"context"
	"net"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/umbracle/atlas/internal/proto"
	"github.com/umbracle/atlas/internal/state"
	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedAtlasServiceServer

	logger     hclog.Logger
	grpcServer *grpc.Server

	state *state.State
}

func NewServer(logger hclog.Logger) (*Server, error) {
	state, err := state.NewState("")
	if err != nil {
		return nil, err
	}
	s := &Server{
		logger: logger,
		state:  state,
	}

	s.grpcServer = grpc.NewServer(s.withLoggingUnaryInterceptor())
	proto.RegisterAtlasServiceServer(s.grpcServer, s)

	// grpc address
	if err := s.setupGRPCServer("localhost:3030"); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) setupGRPCServer(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Error("failed to serve grpc server", "err", err)
		}
	}()

	s.logger.Info("Server started", "addr", addr)
	return nil
}

func (s *Server) Close() {
	s.grpcServer.Stop()
}

func (s *Server) withLoggingUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(s.loggingServerInterceptor)
}

func (s *Server) loggingServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	s.logger.Trace("Request", "method", info.FullMethod, "duration", time.Since(start), "error", err)
	return h, err
}

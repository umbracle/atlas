package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/umbracle/atlas/internal/framework"
	"github.com/umbracle/atlas/internal/proto"
	"github.com/umbracle/atlas/internal/providers/aws"
	"github.com/umbracle/atlas/internal/server/state"
	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedAtlasServiceServer

	logger     hclog.Logger
	grpcServer *grpc.Server

	evalQueue *evalQueue
	state     *state.State

	nodesWatcher *nodesWatcher
}

func NewServer(logger hclog.Logger) (*Server, error) {
	state, err := state.NewState("my.db")
	if err != nil {
		return nil, err
	}
	s := &Server{
		logger:       logger,
		state:        state,
		evalQueue:    newEvalQueue(),
		nodesWatcher: newNodesWatcher(logger),
	}

	s.grpcServer = grpc.NewServer(s.withLoggingUnaryInterceptor())
	proto.RegisterAtlasServiceServer(s.grpcServer, s)

	// grpc address
	if err := s.setupGRPCServer("localhost:3030"); err != nil {
		return nil, err
	}

	// create aws provider if not exists
	providers, err := s.state.ListProviders()
	if err != nil {
		panic(err)
	}
	if len(providers) == 0 {
		if err := s.state.CreateProvider(&proto.Provider{Id: UUID(), Name: "aws", Provider: "aws"}); err != nil {
			panic(err)
		}
	}

	// re-attach all the running nodes
	// for now we just call the function that for sure starts the grpc
	nodes, err := s.state.ListNodes()
	if err != nil {
		panic(err)
	}
	for _, node := range nodes {
		s.nodesWatcher.handleNodeUpdate(node)
	}

	go s.runScheduler()

	s.logger.Info("Up and running")
	return s, nil
}

func (s *Server) instanceProvider(id string) (framework.Provider, error) {
	providers, err := s.state.ListProviders()
	if err != nil {
		return nil, err
	}
	var provider *proto.Provider
	for _, elem := range providers {
		if elem.Id == id {
			provider = elem
		}
	}
	if provider == nil {
		return nil, fmt.Errorf("provider '%s' not found", id)
	}
	prov, ok := Providers[provider.Provider]
	if !ok {
		return nil, fmt.Errorf("provider backend '%s' not found", provider.Provider)
	}
	return prov, nil
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

func (s *Server) runScheduler() {
	for {
		eval := s.evalQueue.pop(context.Background()) // TODO with a stop

		// do some eval stuff
		s.logger.Info("Evaluation", "node", eval.Node)

		if err := s.handleEval(eval); err != nil {
			s.logger.Error("failed to eval", "err", err)
		}
	}
}

func (s *Server) handleEval(eval *proto.Evaluation) error {
	node, err := s.state.LoadNode(eval.Node)
	if err != nil {
		return err
	}
	provider, err := s.instanceProvider(node.ProviderId)
	if err != nil {
		return err
	}

	s.logger.Info("Update node", "node", eval.Node)

	if err := provider.Update(context.Background(), node); err != nil {
		return err
	}
	if err := s.upsertNode(node); err != nil {
		return err
	}

	return nil
}

func (s *Server) upsertNode(node *proto.Node) error {
	if err := s.state.UpsertNode(node); err != nil {
		return err
	}
	s.nodesWatcher.handleNodeUpdate(node)
	return nil
}

func (s *Server) Close() {
	s.grpcServer.Stop()
	s.nodesWatcher.Close()
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

var Providers = map[string]framework.Provider{
	"aws": &aws.AwsProvider{},
}

package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/umbracle/atlas/internal/framework"
	"github.com/umbracle/atlas/internal/proto"
	"github.com/umbracle/atlas/plugins"
)

func (s *Server) Deploy(ctx context.Context, req *proto.DeployRequest) (*proto.DeployResponse, error) {
	fmt.Println("- deploy")

	if req.Id != "" {
		// its an update
		node, err := s.state.LoadNode(req.Id)
		if err != nil {
			return nil, err
		}
		node.ExpectedConfig = req.Args
		if err := s.state.UpsertNode(node); err != nil {
			return nil, err
		}

		// add an evaluation to start the scheduling
		s.evalQueue.add(&proto.Evaluation{
			Node: node.Id,
		})

		return &proto.DeployResponse{Node: node}, nil
	}

	plugin, ok := plugins.Plugins[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", req.Plugin)
	}

	if _, err := s.instanceProvider(req.ProviderId); err != nil {
		return nil, err
	}

	// check if chain exists
	existsChain := false
	for _, chain := range plugin.Chains() {
		if chain == req.Chain {
			existsChain = true
		}
	}
	if !existsChain {
		return nil, fmt.Errorf("chain %s is not available", req.Chain)
	}

	if req.Config != "" {
		config := plugin.Config()
		if err := json.Unmarshal([]byte(req.Config), &config); err != nil {
			return nil, err
		}
	}

	input := &framework.Input{
		Chain: req.Chain,
	}
	nodeSpec := plugin.Build(input)

	nodeSpec.Expected = proto.NodeSpec_Running // important
	nodeSpec.Volume = &proto.NodeSpec_Volume{
		Size: 100,
	}
	node := &proto.Node{
		Id:             UUID(),
		Chain:          req.Chain,
		Spec:           nodeSpec,
		ProviderId:     req.ProviderId,
		ExpectedConfig: req.Args,
	}
	if err := s.state.UpsertNode(node); err != nil {
		return nil, err
	}

	// add an evaluation to start the scheduling
	s.evalQueue.add(&proto.Evaluation{
		Node: node.Id,
	})

	resp := &proto.DeployResponse{
		Node: node,
	}
	return resp, nil
}

func (s *Server) ListNodes(ctx context.Context, req *proto.ListNodesRequest) (*proto.ListNodesResponse, error) {
	nodes, err := s.state.ListNodes()
	if err != nil {
		return nil, err
	}
	resp := &proto.ListNodesResponse{
		Nodes: nodes,
	}
	return resp, nil
}

func (s *Server) GetProviderByName(ctx context.Context, req *proto.GetProviderByNameRequest) (*proto.Provider, error) {
	providers, err := s.state.ListProviders()
	if err != nil {
		return nil, err
	}
	for _, provider := range providers {
		if provider.Name == req.Name {
			return provider, nil
		}
	}
	return nil, fmt.Errorf("provider by name not found")
}

func (s *Server) ListProviders(ctx context.Context, req *proto.ListProvidersRequest) (*proto.ListProvidersResponse, error) {
	providers, err := s.state.ListProviders()
	if err != nil {
		return nil, err
	}
	resp := &proto.ListProvidersResponse{
		Providers: providers,
	}
	return resp, nil
}

func UUID() string {
	return uuid.New().String()
}

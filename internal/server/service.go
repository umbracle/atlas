package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/umbracle/atlas/internal/framework"
	"github.com/umbracle/atlas/internal/proto"
	"github.com/umbracle/atlas/internal/runtime"
	"github.com/umbracle/atlas/plugins"
)

func (s *Server) Deploy(ctx context.Context, req *proto.DeployRequest) (*proto.DeployResponse, error) {
	d, err := runtime.NewDocker()
	if err != nil {
		return nil, err
	}

	plugin, ok := plugins.Plugins[req.Plugin]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", req.Plugin)
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

	nodeSpec.Volume = &proto.NodeSpec_Volume{
		Size: 100,
	}
	node := &proto.Node{
		Id:    UUID(),
		Chain: req.Chain,
		Spec:  nodeSpec,
	}
	handle, err := d.Run(context.Background(), nodeSpec)
	if err != nil {
		return nil, err
	}
	node.Handle = handle

	if err := s.state.UpsertNode(node); err != nil {
		return nil, err
	}
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

func UUID() string {
	return uuid.New().String()
}

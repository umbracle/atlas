package agent

import (
	"context"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/umbracle/atlas/internal/proto"
)

func (a *Agent) CreateService(ctx context.Context, req *proto.CreateServiceRequest) (*proto.CreateRequestResponse, error) {
	// Create service only updates and notifies which is the current spec expected from this node
	go func() {
		a.specUpdate <- req.Spec
	}()
	return &proto.CreateRequestResponse{}, nil
}

func (a *Agent) Do(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (a *Agent) Stream(req *empty.Empty, stream proto.AgentService_StreamServer) error {
	for {
		event := <-a.emitCh

		err := stream.Send(event)
		if err != nil {
			a.logger.Error(err.Error())
			return io.EOF
		}
	}
	return nil
}

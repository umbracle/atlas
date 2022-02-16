package agent

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/umbracle/atlas/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

/*
func TestService_Stream(t *testing.T) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "agent",
		Level: hclog.LevelFromString("info"),
	})
	agent, err := NewAgent(logger, &Config{})
	assert.NoError(t, err)
	defer agent.Close()

	conn, err := grpc.Dial("0.0.0.0:5454", grpc.WithInsecure())
	assert.NoError(t, err)

	clt := proto.NewAgentServiceClient(conn)

	stream, err := clt.Stream(context.Background(), &emptypb.Empty{})
	assert.NoError(t, err)

	msgCh := make(chan *proto.StreamResponse_Event)
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				panic(err)
			}
			msgCh <- msg.Event
		}
	}()

	agent.emitMsg("a")
	msg := <-msgCh
	fmt.Println(msg.Text)

	agent.emitMsg("b")
	msg = <-msgCh
	fmt.Println(msg.Text)
}
*/

func TestService_Create(t *testing.T) {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "agent",
		Level: hclog.LevelFromString("info"),
	})
	agent, err := NewAgent(logger, &Config{})
	assert.NoError(t, err)
	defer agent.Close()

	conn, err := grpc.Dial("0.0.0.0:5454", grpc.WithInsecure())
	assert.NoError(t, err)

	clt := proto.NewAgentServiceClient(conn)

	stream, err := clt.Stream(context.Background(), &emptypb.Empty{})
	assert.NoError(t, err)

	msgCh := make(chan *proto.StreamResponse)
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				panic(err)
			}
			msgCh <- msg
		}
	}()

	clt.CreateService(context.Background(), &proto.CreateServiceRequest{
		Spec: &proto.NodeSpec{
			Image: &proto.NodeSpec_Image{
				Image: "ethereum/client-go",
				Ref:   "v1.9.25",
			},
		},
	})

	for i := 0; i < 2; i++ {
		msg := <-msgCh
		fmt.Println("-- msg --")
		fmt.Println(msg)
	}
}

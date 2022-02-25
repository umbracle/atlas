package agent

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/umbracle/atlas/internal/agent/docker"
	"github.com/umbracle/atlas/internal/proto"
	"google.golang.org/grpc"

	mountutils "k8s.io/mount-utils"
	utilexec "k8s.io/utils/exec"
)

type Config struct {
}

type Agent struct {
	proto.UnimplementedAgentServiceServer

	logger hclog.Logger
	config *Config
	driver *docker.Docker

	grpcServer *grpc.Server
	emitCh     chan *proto.StreamResponse

	// current docker container running
	specUpdate chan *proto.NodeSpec
}

func NewAgent(logger hclog.Logger, config *Config) (*Agent, error) {
	agent := &Agent{
		logger:     logger,
		config:     config,
		emitCh:     make(chan *proto.StreamResponse, 10),
		specUpdate: make(chan *proto.NodeSpec, 5),
	}

	agent.grpcServer = grpc.NewServer(agent.withLoggingUnaryInterceptor())
	proto.RegisterAgentServiceServer(agent.grpcServer, agent)

	// grpc address
	if err := agent.setupGRPCServer("0.0.0.0:5454"); err != nil {
		return nil, err
	}

	if err := agent.handleVolume(); err != nil {
		return nil, err
	}

	driver, err := docker.NewDocker()
	if err != nil {
		return nil, err
	}
	agent.driver = driver

	go agent.reconcile()

	return agent, nil
}

func newSafeMounter() *mountutils.SafeFormatAndMount {
	return &mountutils.SafeFormatAndMount{
		Interface: mountutils.New(""),
		Exec:      utilexec.New(),
	}
}

func (a *Agent) handleVolume() error {

	fmt.Println(os.MkdirAll("/data", 0700))

	mounter := newSafeMounter()

	mounts, err := mounter.List()
	if err != nil {
		return err
	}
	for _, mount := range mounts {
		fmt.Println(mount)
	}

	resize := mountutils.NewResizeFs(utilexec.New())

	// check if its needs resize
	fmt.Println(resize.NeedResize("/dev/xvdh", "/data"))

	// do the resize
	fmt.Println(mounter.FormatAndMount("/dev/xvdh", "/data", "ext4", []string{}))

	return nil
}

func (a *Agent) emitMsg(resp *proto.StreamResponse) {
	select {
	case a.emitCh <- resp:
	default:
	}
}

func (a *Agent) setupGRPCServer(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		if err := a.grpcServer.Serve(lis); err != nil {
			a.logger.Error("failed to serve grpc server", "err", err)
		}
	}()

	a.logger.Info("Agent started", "addr", addr)
	return nil
}

func (a *Agent) Close() {
	a.grpcServer.Stop()
}

func (a *Agent) withLoggingUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(a.loggingServerInterceptor)
}

func (a *Agent) loggingServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	a.logger.Trace("Request", "method", info.FullMethod, "duration", time.Since(start), "error", err)
	return h, err
}

func (a *Agent) reconcile() {
	// wait for a new update
WAIT:
	a.logger.Info("waiting for updates")
	newSpec := <-a.specUpdate
	a.logger.Info("new spec update", "spec", newSpec)

	// start to reconcile
	for {
		containers, err := a.driver.ListContainers()
		if err != nil {
			panic(err)
		}

		r := &reconciler{
			expected: newSpec,
		}
		if len(containers) > 1 {
			panic("many containers found")
		} else if len(containers) == 1 {
			r.found = containers[0]
		}

		plan := r.reconcile()
		if plan.empty() {
			// wait for the next update
			goto WAIT
		}

		// apply updates
		if plan.stop != nil {
			if err := a.driver.StopID(*plan.stop); err != nil {
				a.emitRawMsg("failed to stop container: " + err.Error())
				a.logger.Error("failed to stop container: %v", err)
			} else {
				a.emitNodeStatus(proto.StreamResponse_Failed)
			}
		}
		if plan.start != nil {
			downloaded, err := a.driver.PullImage(context.Background(), plan.start)
			if downloaded {
				a.emitRawMsg("downloaded image")
			}
			if err != nil {
				a.emitRawMsg("failed to download image " + err.Error())
				a.logger.Error(err.Error())
			} else {
				_, err := a.driver.Run(context.Background(), plan.start)
				if err != nil {
					a.emitRawMsg("failed to start container: " + err.Error())
					a.logger.Error(err.Error())
				} else {
					a.emitNodeStatus(proto.StreamResponse_Started)
				}
			}
		}
	}
}

func (a *Agent) emitNodeStatus(typ proto.StreamResponse_NodeStatus) {
	a.emitMsg(&proto.StreamResponse{
		Resp: &proto.StreamResponse_NodeStatus_{
			NodeStatus: typ,
		},
	})
}

func (a *Agent) emitRawMsg(msg string) {
	a.emitMsg(&proto.StreamResponse{
		Resp: &proto.StreamResponse_Event_{
			Event: &proto.StreamResponse_Event{
				Message: msg,
			},
		},
	})
}

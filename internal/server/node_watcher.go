package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/umbracle/atlas/internal/proto"
	"google.golang.org/grpc"
)

type nodesWatcher struct {
	server *Server
	logger hclog.Logger
	nodes  map[string]*nodeWatcher
	lock   sync.Mutex
}

func newNodesWatcher(server *Server, logger hclog.Logger) *nodesWatcher {
	w := &nodesWatcher{
		server: server,
		logger: logger.Named("watcher"),
		nodes:  map[string]*nodeWatcher{},
	}
	return w
}

func (n *nodesWatcher) Close() {
	for _, w := range n.nodes {
		w.Stop()
	}
}

func (n *nodesWatcher) handleNodeUpdate(node *proto.Node) {
	if node.Handle != nil {
		n.update(node)
	} else {
		n.remove(node)
	}
}

func (n *nodesWatcher) update(node *proto.Node) {
	n.logger.Info("update", "node", node.Id)

	n.lock.Lock()
	defer n.lock.Unlock()

	h, ok := n.nodes[node.Id]
	if ok {
		h.notifyUpdate(node)
		return
	}

	watcher := newNodeWatcher(n.server, n.logger.Named(node.Id), node)
	n.nodes[node.Id] = watcher
}

func (n *nodesWatcher) remove(node *proto.Node) {
	n.logger.Info("remove", "node", node.Id)

	n.lock.Lock()
	defer n.lock.Unlock()

	if watcher, ok := n.nodes[node.Id]; ok {
		watcher.Stop()
		delete(n.nodes, node.Id)
	}
}

type nodeWatcher struct {
	server *Server
	logger hclog.Logger

	ctx      context.Context
	cancelFn context.CancelFunc

	node     *proto.Node
	lock     sync.Mutex
	updateCh chan struct{}
}

func newNodeWatcher(server *Server, logger hclog.Logger, node *proto.Node) *nodeWatcher {
	ctx, cancelFn := context.WithCancel(context.Background())

	n := &nodeWatcher{
		server:   server,
		logger:   logger,
		ctx:      ctx,
		cancelFn: cancelFn,
		node:     node,
		updateCh: make(chan struct{}),
	}

	go n.watch()
	return n
}

func (n *nodeWatcher) Stop() {
	n.cancelFn()
}

func (n *nodeWatcher) watch() {
	n.logger.Info("watch")

	var clt proto.AgentServiceClient

	// wait to reach Grpc connection to be established
BACK:
	// try to connect with it
	conn, err := grpc.Dial(n.node.Handle.Ip+":5454", grpc.WithInsecure())
	if err != nil {
		fmt.Println("- err ", err)
		goto BACK
	} else {
		clt = proto.NewAgentServiceClient(conn)
		if _, err := clt.Do(context.Background(), &empty.Empty{}); err != nil {

			time.Sleep(1 * time.Second)

			n.logger.Trace("failed to connect to node", "err", err)

			goto BACK
		}
	}

	// send the request to create the stuff
	n.logger.Info("connected")

	// open the stream
	stream, err := clt.Stream(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}

	streamStop := make(chan struct{})

	go func() {
		defer close(streamStop)

		for {
			msg, err := stream.Recv()
			if err != nil {
				n.logger.Error("update stream stopped")
				return
			} else {
				if nodeStatus, ok := msg.GetResp().(*proto.StreamResponse_NodeStatus_); ok {
					if nodeStatus.NodeStatus == proto.StreamResponse_Done {
						// update the node
						nn, err := n.server.state.LoadNode(n.node.Id)
						if err != nil {
							panic(err)
						}
						nn.NodeStatus = proto.Node_Done
						if err := n.server.upsertNode(nn); err != nil {
							panic(err)
						}
						n.server.handleEval(&proto.Evaluation{
							Node: n.node.Id,
						})

						if err := n.server.state.AddNodeEvent(n.node.Id, "done"); err != nil {
							n.logger.Error("failed to insert event", "err", err)
						}
					} else {
						panic("this should not happen yet")
					}
				} else if nodeMsg, ok := msg.GetResp().(*proto.StreamResponse_Event_); ok {
					if err := n.server.state.AddNodeEvent(n.node.Id, nodeMsg.Event.Message); err != nil {
						n.logger.Error("failed to insert event", "err", err)
					}
				}

				n.logger.Info("msg", "text", msg)
			}
		}
	}()

	// just send the update of the service
	for {
		n.lock.Lock()
		node := n.node
		n.lock.Unlock()

		n.logger.Info("send update")

		req := &proto.CreateServiceRequest{
			Spec: node.Spec,
		}
		if _, err := clt.CreateService(context.Background(), req); err != nil {
			fmt.Println(err)
		}

		select {
		case <-streamStop:
			goto BACK
		case <-n.updateCh:
		case <-n.ctx.Done():
			return
		}
	}
}

func (n *nodeWatcher) notifyUpdate(node *proto.Node) {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.node = node
	select {
	case n.updateCh <- struct{}{}:
	default:
	}
}

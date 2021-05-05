package node

import (
	"context"
	"sync"

	"github.com/spf13/viper"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	"github.com/thetatoken/theta-eth-rpc-adaptor/rpc/ethrpc"
)

type Node struct {
	RPC *ethrpc.RPCAdaptorServer

	// Life cycle
	wg      *sync.WaitGroup
	quit    chan struct{}
	ctx     context.Context
	cancel  context.CancelFunc
	stopped bool
}

func NewNode() *Node {
	node := &Node{}

	if viper.GetBool(common.CfgRPCEnabled) {
		node.RPC = ethrpc.NewRPCAdaptorServer()
	}
	return node
}

// Start starts sub components and kick off the main loop.
func (n *Node) Start(ctx context.Context) {
	c, cancel := context.WithCancel(ctx)
	n.ctx = c
	n.cancel = cancel

	if viper.GetBool(common.CfgRPCEnabled) {
		n.RPC.Start(n.ctx)
	}
}

// Stop notifies all sub components to stop without blocking.
func (n *Node) Stop() {
	n.cancel()
}

// Wait blocks until all sub components stop.
func (n *Node) Wait() {
	if n.RPC != nil {
		n.RPC.Wait()
	}
}

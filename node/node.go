package node

import (
	"context"
	"sync"

	"github.com/spf13/viper"
	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	"github.com/thetatoken/theta-eth-rpc-adaptor/rpc"

	erpclib "github.com/ethereum/go-ethereum/rpc"
)

type Node struct {
	// Life cycle
	wg      *sync.WaitGroup
	quit    chan struct{}
	ctx     context.Context
	cancel  context.CancelFunc
	stopped bool
}

func NewNode() *Node {
	node := &Node{
		wg: &sync.WaitGroup{},
	}

	return node
}

// Start starts sub components and kick off the main loop.
func (n *Node) Start(ctx context.Context) {
	c, cancel := context.WithCancel(ctx)
	n.ctx = c
	n.cancel = cancel

	if viper.GetBool(common.CfgRPCEnabled) {
		rpc.StartServers([]erpclib.API{})
	}

	n.wg.Add(1)
	go n.mainLoop()
}

// Stop notifies all sub components to stop without blocking.
func (n *Node) Stop() {
	n.cancel()

	rpc.StopServers()
}

// Wait blocks until all sub components stop.
func (n *Node) Wait() {
	n.wg.Wait()
}

func (n *Node) mainLoop() {
	defer n.wg.Done()

	<-n.ctx.Done()
	n.stopped = true
}

package ethrpc

import (
	"sync"
	"time"

	erpclib "github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/thetatoken/theta/common/timer"
)

var blockInterval time.Duration = 6 * time.Second

var logger *log.Entry = log.WithFields(log.Fields{"prefix": "ethrpc"})

// EthRPCService provides an API to access to the Eth endpoints.
type EthRPCService struct {
	pendingHeavyEthLogQueryCounter           uint64
	pendingHeavyEthLogQueryCounterLock       *sync.Mutex
	pendingHeavyEthLogQueryCounterResetTimer *timer.RepeatTimer
}

// NewEthRPCService creates a new API for the Ethereum RPC interface
func NewEthRPCService(namespace string) erpclib.API {
	if namespace == "" {
		namespace = "eth"
	}

	ethRPCService := &EthRPCService{
		pendingHeavyEthLogQueryCounter:           0,
		pendingHeavyEthLogQueryCounterLock:       &sync.Mutex{},
		pendingHeavyEthLogQueryCounterResetTimer: timer.NewRepeatTimer("pendingHeavyEthLogQueryCounterReset", 30*time.Minute),
	}
	ethRPCService.pendingHeavyEthLogQueryCounterResetTimer.Reset()
	ethRPCService.mainLoop()

	return erpclib.API{
		Namespace: namespace,
		Version:   "1.0",
		Service:   ethRPCService,
		Public:    true,
	}
}

func (e *EthRPCService) mainLoop() {
	for {
		select {
		case <-e.pendingHeavyEthLogQueryCounterResetTimer.Ch:
			e.pendingHeavyEthLogQueryCounterLock.Lock()
			e.pendingHeavyEthLogQueryCounter = 0 // reset the counter to zero at a fixed time interval, otherwise the counter could get stuck if some pending queries never return
			e.pendingHeavyEthLogQueryCounterLock.Unlock()
		}
	}
}

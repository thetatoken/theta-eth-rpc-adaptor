package ethrpc

import (
	"context"
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
	pendingHeavyGetLogsCounter           uint64
	pendingHeavyGetLogsCounterLock       *sync.Mutex
	pendingHeavyGetLogsCounterResetTimer *timer.RepeatTimer

	// Life cycle
	wg      *sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	stopped bool
}

// NewEthRPCService creates a new API for the Ethereum RPC interface
func NewEthRPCService(namespace string) erpclib.API {
	if namespace == "" {
		namespace = "eth"
	}

	serv := &EthRPCService{
		wg: &sync.WaitGroup{},

		pendingHeavyGetLogsCounter:           0,
		pendingHeavyGetLogsCounterLock:       &sync.Mutex{},
		pendingHeavyGetLogsCounterResetTimer: timer.NewRepeatTimer("pendingHeavyGetLogsCounterReset", 5*time.Minute),
	}
	serv.pendingHeavyGetLogsCounterResetTimer.Reset()

	serv.wg.Add(1)
	go serv.heavyQueryCounterLoop()

	return erpclib.API{
		Namespace: namespace,
		Version:   "1.0",
		Service:   serv,
		Public:    true,
	}
}

func (serv *EthRPCService) heavyQueryCounterLoop() {
	defer serv.wg.Done()

	for {
		select {
		case <-serv.pendingHeavyGetLogsCounterResetTimer.Ch:
			serv.pendingHeavyGetLogsCounterLock.Lock()
			serv.pendingHeavyGetLogsCounter = 0 // reset the counter to zero at a fixed time interval, otherwise the counter could get stuck if some pending queries never return
			serv.pendingHeavyGetLogsCounterLock.Unlock()
		}
	}
}

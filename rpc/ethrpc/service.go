package ethrpc

import (
	"time"

	erpclib "github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

var blockInterval time.Duration = 6 * time.Second

var logger *log.Entry = log.WithFields(log.Fields{"prefix": "ethrpc"})

// EthRPCService provides an API to access to the Eth endpoints.
type EthRPCService struct {
}

// NewEthRPCService creates a new API for the Ethereum RPC interface
func NewEthRPCService(namespace string) erpclib.API {
	if namespace == "" {
		namespace = "eth"
	}

	return erpclib.API{
		Namespace: namespace,
		Version:   "1.0",
		Service:   &EthRPCService{},
		Public:    true,
	}
}

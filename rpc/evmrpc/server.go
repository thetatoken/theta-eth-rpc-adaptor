package evmrpc

import (
	erpclib "github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

var logger *log.Entry = log.WithFields(log.Fields{"prefix": "evmrpc"})

// EvmRPCService provides an API to access to the evm endpoints.
type EvmRPCService struct {
}

// NewEvmRPCService creates a new API for the evm RPC interface
func NewEvmRPCService(namespace string) erpclib.API {
	if namespace == "" {
		namespace = "evm"
	}

	return erpclib.API{
		Namespace: namespace,
		Version:   "1.0",
		Service:   &EvmRPCService{},
		Public:    true,
	}
}

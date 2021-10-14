package web3rpc

import (
	erpclib "github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

var logger *log.Entry = log.WithFields(log.Fields{"prefix": "web3rpc"})

// Web3RPCService provides an API to access to the Web3 endpoints.
type Web3RPCService struct {
}

// NewWeb3RPCService creates a new API for the Web3ereum RPC interface
func NewWeb3RPCService(namespace string) erpclib.API {
	if namespace == "" {
		namespace = "web3"
	}

	return erpclib.API{
		Namespace: namespace,
		Version:   "1.0",
		Service:   &Web3RPCService{},
		Public:    true,
	}
}

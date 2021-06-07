package netrpc

import (
	erpclib "github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

var logger *log.Entry = log.WithFields(log.Fields{"prefix": "netrpc"})

// NetRPCService provides an API to access to the Net endpoints.
type NetRPCService struct {
}

// NewNetRPCService creates a new API for the Netereum RPC interface
func NewNetRPCService(namespace string) erpclib.API {
	if namespace == "" {
		namespace = "net"
	}

	return erpclib.API{
		Namespace: namespace,
		Version:   "1.0",
		Service:   &NetRPCService{},
		Public:    true,
	}
}

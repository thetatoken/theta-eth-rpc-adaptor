package ethrpc

import (
	"context"
)

// ------------------------------- eth_chainId -----------------------------------

func (e *EthRPCService) ChainId(ctx context.Context) (result string, err error) {
	logger.Infof("eth_chainId called")

	//TODO: change to the correct chain ID with theta RPC
	result = "0x18888"
	return result, nil
}

package evmrpc

import (
	"context"
)

// ------------------------------- evm_mine -----------------------------------

func (e *EvmRPCService) Mine(ctx context.Context) (result string, err error) {
	logger.Infof("evm_mine called")
	result = "0x0"
	return result, nil
}

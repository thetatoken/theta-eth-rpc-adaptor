package evmrpc

import (
	"context"
)

// ------------------------------- evm_revert -----------------------------------

func (e *EvmRPCService) Revert(ctx context.Context) (result bool, err error) {
	logger.Infof("evm_revert called")
	result = true
	return result, nil
}

package evmrpc

import (
	"context"
)

// ------------------------------- evm_snapshot -----------------------------------

func (e *EvmRPCService) Snapshot(ctx context.Context) (result string, err error) {
	logger.Infof("evm_snapshot called")
	result = "0x1"
	return result, nil
}

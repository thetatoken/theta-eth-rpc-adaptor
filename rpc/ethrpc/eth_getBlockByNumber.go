package ethrpc

import (
	"context"
)

// ------------------------------- eth_getBlockByNumber -----------------------------------

func (e *EthRPCService) GetBlockByNumber(ctx context.Context, tag string, b bool) (result string, err error) {
	logger.Infof("eth_getBlockByNumber called")

	return result, nil
}

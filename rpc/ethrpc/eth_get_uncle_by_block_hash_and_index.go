package ethrpc

import (
	"context"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
)

// ------------------------------- eth_getUncleByBlockHashAndIndex -----------------------------------
func (e *EthRPCService) GetUncleByBlockHashAndIndex(ctx context.Context, hashStr string, indexStr string) (result common.EthGetBlockResult, err error) {
	logger.Infof("eth_getUncleByBlockHashAndIndex called")
	//This is a place holder
	return common.EthGetBlockResult{}, nil
}

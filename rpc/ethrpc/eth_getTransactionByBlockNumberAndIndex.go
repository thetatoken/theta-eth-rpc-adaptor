package ethrpc

import (
	"context"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	trpc "github.com/thetatoken/theta/rpc"

	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getTransactionByBlockNumberAndIndex -----------------------------------
func (e *EthRPCService) GetTransactionByBlockNumberAndIndex(ctx context.Context, numberStr string, txIndexStr string) (result common.EthGetTransactionResult, err error) {
	logger.Infof("GetTransactionByBlockNumberAndIndex called")
	height := GetHeightByTag(numberStr)
	txIndex := GetHeightByTag(txIndexStr) //TODO: use common
	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetBlockByHeight", trpc.GetBlockByHeightArgs{Height: height})
	return GetIndexedTransactionFromBlock(rpcRes, rpcErr, txIndex)
}

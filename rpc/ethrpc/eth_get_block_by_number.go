package ethrpc

import (
	"context"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getBlockByNumber -----------------------------------
func (e *EthRPCService) GetBlockByNumber(ctx context.Context, numberStr string, txDetails bool) (result common.EthGetBlockResult, err error) {
	logger.Infof("eth_getBlockByNumber called")
	height := common.GetHeightByTag(numberStr)

	if height == 0 {
		height, err = common.GetCurrentHeight()

		if err != nil {
			return result, err
		}
	}

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetBlockByHeight", trpc.GetBlockByHeightArgs{
		Height: height})
	return GetBlockFromTRPCResult(rpcRes, rpcErr, txDetails)
}

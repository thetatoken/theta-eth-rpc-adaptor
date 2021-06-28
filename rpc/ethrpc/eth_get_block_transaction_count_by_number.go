package ethrpc

import (
	"context"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	hexutil "github.com/thetatoken/theta/common/hexutil"
	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getBlockTransactionCountByNumber -----------------------------------
func (e *EthRPCService) GetBlockTransactionCountByNumber(ctx context.Context, numberStr string) (result hexutil.Uint64, err error) {
	logger.Infof("eth_getBlockTransactionCountByNumber called")
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
	block, err := GetBlockFromTRPCResult(rpcRes, rpcErr, false)
	return hexutil.Uint64(len(block.Transactions)), err
}

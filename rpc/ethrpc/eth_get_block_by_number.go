package ethrpc

import (
	"context"
	"math"
	"math/big"
	"time"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getBlockByNumber -----------------------------------
func (e *EthRPCService) GetBlockByNumber(ctx context.Context, numberStr string, txDetails bool) (result common.EthGetBlockResult, err error) {
	logger.Infof("eth_getBlockByNumber called, blockHeight: %v", numberStr)
	height := common.GetHeightByTag(numberStr)

	if height == math.MaxUint64 {
		height, err = common.GetCurrentHeight()
		if err != nil {
			return result, err
		}
	}

	chainIDStr, err := e.ChainId(ctx)
	if err != nil {
		logger.Errorf("Failed to get chainID\n")
		return result, nil
	}
	chainID := new(big.Int)
	chainID.SetString(chainIDStr, 16)

	maxRetry := 5
	for i := 0; i < maxRetry; i++ { // It might take some time for a block to be finalized, retry a few times
		client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
		rpcRes, rpcErr := client.Call("theta.GetBlockByHeight", trpc.GetBlockByHeightArgs{
			Height: height})

		//logger.Infof("eth_getBlockByNumber, rpcRes: %v, rpcRes.Rsult: %v", rpcRes, rpcRes.Result)

		result, err = GetBlockFromTRPCResult(chainID, rpcRes, rpcErr, txDetails)
		if err == nil {
			return result, err
		}

		time.Sleep(blockInterval) // one block duration
	}

	return result, err
}

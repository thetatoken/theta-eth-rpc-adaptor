package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	hexutil "github.com/thetatoken/theta/common/hexutil"
	"github.com/thetatoken/theta/ledger/types"
	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

type chainIDResultWrapper struct {
	chainID string
}

// ------------------------------- eth_chainId -----------------------------------

func (e *EthRPCService) ChainId(ctx context.Context) (result string, err error) {
	logger.Infof("eth_chainId called")

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetStatus", trpc.GetStatusArgs{})
	var blockHeight uint64
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		re := chainIDResultWrapper{
			chainID: trpcResult.ChainID,
		}
		blockHeight = uint64(trpcResult.LatestFinalizedBlockHeight)
		return re, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", nil
	}
	thetaChainIDResult, ok := resultIntf.(chainIDResultWrapper)
	if !ok {
		return "", nil
	}

	thetaChainID := thetaChainIDResult.chainID
	ethChainID := types.MapChainID(thetaChainID, blockHeight).Uint64()
	result = hexutil.EncodeUint64(ethChainID)

	return result, nil
}

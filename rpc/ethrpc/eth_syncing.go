package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	"github.com/thetatoken/theta/common/hexutil"
	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

type syncingResultWrapper struct {
	*common.EthSyncingResult
	syncing bool
}

// ------------------------------- eth_syncing -----------------------------------
func (e *EthRPCService) Syncing(ctx context.Context) (result interface{}, err error) {
	logger.Infof("eth_syncing called")
	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetStatus", trpc.GetStatusArgs{})
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		re := syncingResultWrapper{&common.EthSyncingResult{}, false}
		re.syncing = trpcResult.Syncing
		if trpcResult.Syncing {
			re.StartingBlock = 1
			re.CurrentBlock = hexutil.Uint64(trpcResult.CurrentHeight)
			re.HighestBlock = hexutil.Uint64(trpcResult.LatestFinalizedBlockHeight)
			re.PulledStates = re.CurrentBlock
			re.KnownStates = re.CurrentBlock
		}
		return re, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", nil
	}
	thetaSyncingResult, ok := resultIntf.(syncingResultWrapper)
	if !ok {
		return nil, nil
	}
	if !thetaSyncingResult.syncing {
		result = false
	} else {
		result = thetaSyncingResult.EthSyncingResult
	}

	return result, nil
}

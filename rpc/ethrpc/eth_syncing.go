package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// type ethSyncingResult struct {
// 	StartingBlock tcommon.JSONUint64 `json:"startingBlock"`
// 	CurrentBlock  tcommon.JSONUint64 `json:"currentBlock"`
// 	HighestBlock  tcommon.JSONUint64 `json:"highestBlock"`
// 	PulledStates  tcommon.JSONUint64 `json:"pulledStates"` //pulledStates is the number it already downloaded
// 	KnownStates   tcommon.JSONUint64 `json:"knownStates"`  //knownStates is the number of trie nodes that the sync algo knows about
// }
type syncingResultWrapper struct {
	*common.EthSyncingResult
	syncing bool
}

// ------------------------------- eth_syncing -----------------------------------
func (e *EthRPCService) Syncing(ctx context.Context) (result interface{}, err error) {
	logger.Infof("eth_syncing called")
	result = "log4"
	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetStatus", trpc.GetStatusArgs{})
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		re := syncingResultWrapper{&common.EthSyncingResult{}, false}
		re.syncing = trpcResult.Syncing
		if trpcResult.Syncing {
			re.StartingBlock = 1
			re.CurrentBlock = trpcResult.CurrentHeight
			re.HighestBlock = trpcResult.LatestFinalizedBlockHeight
			re.PulledStates = trpcResult.CurrentHeight
			re.KnownStates = trpcResult.CurrentHeight
		}
		return re, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	thetaSyncingResult, ok := resultIntf.(syncingResultWrapper)
	if !ok {
		return nil, fmt.Errorf("failed to convert syncingResultWrapper")
	}
	if !thetaSyncingResult.syncing {
		result = false
	} else {
		result = thetaSyncingResult.EthSyncingResult
	}

	return result, nil
}

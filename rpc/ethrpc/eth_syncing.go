package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_syncing -----------------------------------

func (e *EthRPCService) Syncing(ctx context.Context) (result string, err error) {
	logger.Infof("eth_syncing called")
	result = "log4"
	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetStatus", trpc.GetStatusArgs{})
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Syncing, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	result = fmt.Sprintf("%v", resultIntf)

	return result, nil
}

package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	tcommon "github.com/thetatoken/theta/common"
	hexutil "github.com/thetatoken/theta/common/hexutil"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_blockNumber -----------------------------------

func (e *EthRPCService) BlockNumber(ctx context.Context) (result string, err error) {
	logger.Infof("eth_blockNumber called")

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetStatus", trpc.GetStatusArgs{})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.CurrentHeight, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	result = hexutil.EncodeUint64(uint64(resultIntf.(tcommon.JSONUint64)))
	return result, nil
}

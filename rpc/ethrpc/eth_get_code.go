package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getCode -----------------------------------

func (e *EthRPCService) GetCode(ctx context.Context, address string, tag string) (result string, err error) {
	logger.Infof("eth_getCode called")

	height := common.GetHeightByTag(tag)

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetCode", trpc.GetCodeArgs{Address: address, Height: height})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetCodeResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Code, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)

	if err != nil {
		return result, err
	}

	result = resultIntf.(string)
	return result, nil
}

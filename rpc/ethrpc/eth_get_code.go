package ethrpc

import (
	"context"
	"encoding/json"
	"math"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getCode -----------------------------------

func (e *EthRPCService) GetCode(ctx context.Context, address string, tag string) (result string, err error) {
	logger.Infof("eth_getCode called")

	height := common.GetHeightByTag(tag)
	if height == math.MaxUint64 {
		height = 0 // 0 is interpreted as the last height by the theta.GetAccount method
	}

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
	if result == "" {
		result = "0x"
	}

	return result, nil
}

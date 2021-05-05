package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_protocolVersion -----------------------------------

func (e *EthRPCService) ProtocolVersion(ctx context.Context) (result string, err error) {
	logger.Infof("eth_protocolVersion called")

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetVersion", trpc.GetVersionArgs{})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetVersionResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Version, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	result = resultIntf.(string)

	return result, nil
}

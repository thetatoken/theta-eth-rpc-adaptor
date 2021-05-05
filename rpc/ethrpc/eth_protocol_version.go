package ethrpc

import (
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/rpc"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_protocolVersion -----------------------------------

type ProtocolVersionArgs struct {
}

type ProtocolVersionResult struct {
	Result string `json:"result"`
}

func (t *RPCAdaptorService) ProtocolVersion(args *ProtocolVersionArgs, result *ProtocolVersionResult) (err error) {
	logger.Infof("eth_protocolVersion called with args: %v", *args)

	client := rpcc.NewRPCClient(rpc.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetVersion", trpc.GetVersionArgs{})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetVersionResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Version, nil
	}

	resultIntf, err := rpc.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return err
	}
	result.Result = resultIntf.(string)

	return nil
}

package ethrpc

import (
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	"github.com/thetatoken/theta-eth-rpc-adaptor/rpc"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- ProtocolVersion -----------------------------------

type GetVersionArgs struct {
}

func (t *rpc.RPCAdaptorServer) ProtocolVersion(args *GetVersionArgs, result *string) (err error) {
	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetVersion", trpc.GetVersionArgs{})

	parse := func(jsonBytes []byte) (*string, error) {
		result := ""
		json.Unmarshal(jsonBytes, &result)
		return &result, nil
	}

	err, result = rpc.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return err
	}

	return nil
}

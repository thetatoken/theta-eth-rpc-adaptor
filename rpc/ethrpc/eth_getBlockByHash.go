package ethrpc

import (
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	tcommon "github.com/thetatoken/theta/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

type GetBlockArgs struct {
	Hash string `json:"hash"`
}

// ------------------------------- eth_getBlockByHash -----------------------------------
func (e *EthRPCService) GetBlockByHash(args *GetBlockArgs, resp *trpc.GetBlockResult) (result string, err error) {
	logger.Infof("eth_getBlockByHash called")

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetBlock", trpc.GetBlockArgs{Hash: tcommon.HexToHash(args.Hash)})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetBlockResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	resultObj, err := json.Marshal(resultIntf)
	if err != nil {
		return "", err
	}
	result = string(resultObj)
	return result, nil
}

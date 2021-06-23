package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_call -----------------------------------

func (e *EthRPCService) Call(ctx context.Context, argObj common.EthSmartContractArgObj, tag string) (result string, err error) {
	logger.Infof("eth_call called, tx: %+v", argObj)

	sctxBytes, err := common.GetSctxBytes(argObj)
	if err != nil {
		logger.Errorf("eth_call: Failed to get smart contract bytes: %+v\n", argObj)
		return result, err
	}

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())

	rpcRes, rpcErr := client.Call("theta.CallSmartContract", trpc.CallSmartContractArgs{SctxBytes: hex.EncodeToString(sctxBytes)})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.CallSmartContractResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		logger.Infof("eth_call Theta RPC result: %+v\n", trpcResult)
		return trpcResult.VmReturn, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	result = "0x" + resultIntf.(string)

	logger.Infof("result: %v\n", result)

	return result, nil
}

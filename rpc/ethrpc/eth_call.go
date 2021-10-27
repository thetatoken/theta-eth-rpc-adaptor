package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_call -----------------------------------

// Note: "tag" could be an integer block number, or the string "latest", "earliest" or "pending". So its type needs to be interface{}
func (e *EthRPCService) Call(ctx context.Context, argObj common.EthSmartContractArgObj, tag interface{}) (result string, err error) {
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
		if len(trpcResult.VmError) > 0 {
			return trpcResult.GasUsed, fmt.Errorf(trpcResult.VmError)
		}
		return trpcResult.VmReturn, nil
	}

	//logger.Infof("eth_call rpcRes: %v, rpcErr: %v", rpcRes, rpcErr)

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		logger.Infof("eth_call error: %v", err)
		return "", err
	}
	result = "0x" + resultIntf.(string)

	logger.Infof("eth_call result: %v", result)

	return result, nil
}

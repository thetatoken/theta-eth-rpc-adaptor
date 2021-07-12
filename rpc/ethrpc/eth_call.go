package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	// "https://github.com/ethereum/go-ethereum/accounts/abi"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_call -----------------------------------

func (e *EthRPCService) Call(ctx context.Context, argObj common.EthSmartContractArgObj, tag string) (result string, err error) {
	// func (e *EthRPCService) Call(ctx context.Context, argObj common.EthSmartContractArgObj, tag string) (result interface{}, err error) {
	logger.Infof("eth_call called, tx: %+v", argObj)

	sctxBytes, err := common.GetSctxBytes(argObj)
	if err != nil {
		logger.Errorf("eth_call: Failed to get smart contract bytes: %+v\n", argObj)
		return result, err
	}
	logger.Infof("jlog3 sctxBytes", hex.EncodeToString(sctxBytes))

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())

	rpcRes, rpcErr := client.Call("theta.CallSmartContract", trpc.CallSmartContractArgs{SctxBytes: hex.EncodeToString(sctxBytes)})
	parse := func(jsonBytes []byte) (result interface{}, err error) {
		trpcResult := trpc.CallSmartContractResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		logger.Infof("eth_call Theta RPC result: %+v\n", trpcResult)

		if trpcResult.VmError != "" {
			var errStr string
			errStr, err = common.GetErrorMessageFromCallData(trpcResult.VmReturn)
			if err == nil {
				err = fmt.Errorf(errStr)
			}
		}
		return trpcResult.VmReturn, err
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	result = "0x" + resultIntf.(string)

	logger.Infof("result: %v\n", result)

	return result, nil
}

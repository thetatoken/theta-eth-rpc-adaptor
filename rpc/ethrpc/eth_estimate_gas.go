package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	tcommon "github.com/thetatoken/theta/common"
	hexutil "github.com/thetatoken/theta/common/hexutil"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_estimateGas -----------------------------------

func (e *EthRPCService) EstimateGas(ctx context.Context, argObj common.EthSmartContractArgObj) (result string, err error) {
	logger.Infof("eth_estimateGas called")

	sctxBytes, err := common.GetSctxBytes(argObj)
	if err != nil {
		logger.Errorf("eth_estimateGas: Failed to get smart contract bytes: %+v\n", argObj)
		return result, err
	}

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())

	rpcRes, rpcErr := client.Call("theta.CallSmartContract", trpc.CallSmartContractArgs{SctxBytes: hex.EncodeToString(sctxBytes)})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.CallSmartContractResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		if len(trpcResult.VmError) > 0 {
			logger.Warnf("eth_estimateGas: EVM execution failed: %v\n", trpcResult.VmError)
			return trpcResult.GasUsed, fmt.Errorf(trpcResult.VmError)
		}
		return trpcResult.GasUsed, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	result = hexutil.EncodeUint64(2 * uint64(resultIntf.(tcommon.JSONUint64)))
	return result, nil
}

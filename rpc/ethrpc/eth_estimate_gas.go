package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	"github.com/thetatoken/theta/cmd/thetacli/cmd/utils"
	tcommon "github.com/thetatoken/theta/common"
	hexutil "github.com/thetatoken/theta/common/hexutil"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_estimateGas -----------------------------------

func (e *EthRPCService) EstimateGas(ctx context.Context, argObj common.EthSmartContractArgObj, tag string) (result string, err error) {
	logger.Infof("eth_estimateGas called")

	sctxBytes, err := common.GetSctxBytes(argObj)
	if err != nil {
		utils.Error("Failed to get smart contract bytes: %+v\n", argObj)
		return result, err
	}

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())

	rpcRes, rpcErr := client.Call("theta.CallSmartContract", trpc.CallSmartContractArgs{SctxBytes: hex.EncodeToString(sctxBytes)})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.CallSmartContractResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.GasUsed, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	result = hexutil.EncodeUint64(uint64(resultIntf.(tcommon.JSONUint64)))
	return result, nil
}

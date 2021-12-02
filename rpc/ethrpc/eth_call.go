package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_call -----------------------------------

// Note: "tag" could be an integer block number, or the string "latest", "earliest" or "pending". So its type needs to be interface{}
func (e *EthRPCService) Call(ctx context.Context, argObj common.EthSmartContractArgObj, tag interface{}) (result string, err error) {
	logger.Infof("eth_call called, tx: %+v", argObj)

	blockGasLimit := viper.GetUint64(common.CfgThetaBlockGasLimit)
	gas, err := strconv.ParseUint(argObj.Gas, 16, 64)
	if err != nil || gas > blockGasLimit {
		argObj.Gas = "0x" + fmt.Sprintf("%x", blockGasLimit)
	}

	sctxBytes, err := common.GetSctxBytes(argObj)
	if err != nil {
		logger.Errorf("eth_call: Failed to get smart contract bytes: %+v\n", argObj)
		return result, err
	}

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.CallSmartContractResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		logger.Infof("eth_call Theta RPC result: %+v\n", trpcResult)
		if len(trpcResult.VmError) > 0 {
			return trpcResult.GasUsed, fmt.Errorf(trpcResult.VmError)
		}
		return trpcResult.VmReturn, nil
	}

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())

	maxRetry := 5
	for i := 0; i < maxRetry; i++ { // It might take some time for a tx to be finalized, retry a few times
		rpcRes, rpcErr := client.Call("theta.CallSmartContract", trpc.CallSmartContractArgs{
			SctxBytes: hex.EncodeToString(sctxBytes),
			Preview:   true,
		})
		//logger.Infof("eth_call rpcRes: %v, rpcErr: %v", rpcRes, rpcErr)

		resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
		if err != nil {
			if i == maxRetry-1 {
				logger.Infof("eth_call error: %v", err)
				return "", err
			}
			time.Sleep(blockInterval) // one block duration
		} else {
			result = "0x" + resultIntf.(string)
			break
		}
	}

	logger.Infof("eth_call result: %v", result)

	return result, nil
}

package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_sendTransaction -----------------------------------
func (e *EthRPCService) SendTransaction(ctx context.Context, argObj common.EthSmartContractArgObj) (result string, err error) {
	logger.Infof("eth_sendTransaction called, tx: %+v", argObj)

	blockNumber, err := e.BlockNumber(ctx)
	if err != nil {
		logger.Errorf("eth_sendTransaction, failed to get blocknumber\n")
		return "", nil
	}
	chainID, err := e.ChainId(ctx)
	if err != nil {
		logger.Errorf("eth_sendTransaction, failed to get chainID\n")
		return "", nil
	}

	signedTx, err := common.GetSignedBytes(argObj, chainID, blockNumber)
	if err != nil {
		return "", nil
	}
	logger.Infof("eth_sendTransaction broadcasting signedTX: %v\n", signedTx)
	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.BroadcastRawTransactionAsync", trpc.BroadcastRawTransactionAsyncArgs{TxBytes: signedTx})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.BroadcastRawTransactionAsyncResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.TxHash, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		//logger.Errorf("eth_sendTransaction, err: %v, result: %v", err, resultIntf.(string))
		return "", err
	}
	result = resultIntf.(string)

	logger.Infof("eth_sendTransaction, result: %v", result)

	return result, nil

}

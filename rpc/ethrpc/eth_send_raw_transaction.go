package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_sendRawTransaction -----------------------------------

func (e *EthRPCService) SendRawTransaction(ctx context.Context, txBytes string) (result string, err error) {
	logger.Infof("eth_sendRawTransaction called")

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.BroadcastRawEthTransactionAsync", trpc.BroadcastRawTransactionAsyncArgs{TxBytes: txBytes})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.BroadcastRawTransactionAsyncResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.TxHash, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		logger.Errorf("eth_sendRawTransaction, err: , result: \n", err, resultIntf.(string))
		return "", err
	}
	result = resultIntf.(string)

	logger.Infof("eth_sendRawTransaction, result: %v\n", result)

	return result, nil
}

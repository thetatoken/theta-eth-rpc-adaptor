package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/ledger/types"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getTransactionByHash -----------------------------------
func (e *EthRPCService) GetTransactionByHash(ctx context.Context, hashStr string) (result common.EthGetTransactionResult, err error) {
	logger.Infof("eth_getTransactionByHash called")

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetTransaction", trpc.GetTransactionArgs{Hash: hashStr})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetTransactionResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		var objmap map[string]json.RawMessage
		json.Unmarshal(jsonBytes, &objmap)
		if objmap["transaction"] != nil {
			//TODO: handle other types
			if types.TxType(trpcResult.Type) == types.TxSend {
				tx := types.SendTx{}
				json.Unmarshal(objmap["transaction"], &tx)
				trpcResult.Tx = &tx
			}
		}
		return trpcResult, nil
	}
	result = common.EthGetTransactionResult{}
	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return result, err
	}
	thetaGetTransactionResult, ok := resultIntf.(trpc.GetTransactionResult)
	if !ok {
		return result, fmt.Errorf("failed to convert GetBlockResult")
	}
	result.BlockHash = thetaGetTransactionResult.BlockHash
	result.BlockHeight = thetaGetTransactionResult.BlockHeight
	if thetaGetTransactionResult.Tx != nil {
		//TODO: handle other types
		if types.TxType(thetaGetTransactionResult.Type) == types.TxSend {
			tx := thetaGetTransactionResult.Tx.(*types.SendTx)
			result.From = tx.Inputs[0].Address
			result.To = tx.Outputs[0].Address
			result.Gas = (tcommon.JSONBig)(*tx.Fee.TFuelWei)
			// Theta/Tfuel value? Sum of all inputs amount?
			result.Value = (tcommon.JSONBig)(*tx.Inputs[0].Coins.TFuelWei)
			// TODO: handle r,v,s, transactionIndex maybe?
		}
	}

	return result, nil
}

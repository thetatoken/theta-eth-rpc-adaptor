package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/ledger/types"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getBlockByHash -----------------------------------
func (e *EthRPCService) GetBlockByHash(ctx context.Context, hashStr string, txDetails bool) (result common.EthGetBlockResult, err error) {
	logger.Infof("eth_getBlockByHash called")

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetBlock", trpc.GetBlockArgs{Hash: tcommon.HexToHash(hashStr)})
	result = common.EthGetBlockResult{}
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetBlockResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		logger.Infof("transactions count %d \n", len(trpcResult.Txs))
		result.Transactions = make([]interface{}, len(trpcResult.Txs))
		var objmap map[string]json.RawMessage
		json.Unmarshal(jsonBytes, &objmap)
		if objmap["transactions"] != nil {
			var txmaps []map[string]json.RawMessage
			json.Unmarshal(objmap["transactions"], &txmaps)
			for i, omap := range txmaps {
				if types.TxType(trpcResult.Txs[i].Type) == types.TxSmartContract {
					if txDetails {
						scTx := types.SmartContractTx{}
						json.Unmarshal(omap["raw"], &scTx)
						result.Transactions[i] = scTx
					}
					result.GasUsed = tcommon.JSONUint64(trpcResult.Txs[i].Receipt.GasUsed)
				} else if txDetails && types.TxType(trpcResult.Txs[i].Type) == types.TxSend {
					sTx := types.SendTx{}
					json.Unmarshal(omap["raw"], &sTx)
					result.Transactions[i] = sTx
				}
			}
		}
		return trpcResult, nil
	}
	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return result, err
	}
	theta_GetBlockResult := resultIntf.(trpc.GetBlockResult)
	result.Height = theta_GetBlockResult.Height
	result.Hash = theta_GetBlockResult.Hash
	result.Parent = theta_GetBlockResult.Parent
	result.Timestamp = theta_GetBlockResult.Timestamp
	result.Proposer = theta_GetBlockResult.Proposer
	result.TxHash = theta_GetBlockResult.TxHash
	result.StateHash = theta_GetBlockResult.StateHash
	for i, tx := range theta_GetBlockResult.Txs {
		if txDetails && (types.TxType(tx.Type) == types.TxSmartContract || types.TxType(tx.Type) == types.TxSend) {
			//already handled
		} else {
			result.Transactions[i] = tx.Hash
		}
	}
	result.GasLimit = 20000000
	return result, nil
}

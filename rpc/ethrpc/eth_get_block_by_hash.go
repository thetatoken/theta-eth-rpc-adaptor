package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/common/hexutil"
	"github.com/thetatoken/theta/ledger/types"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getBlockByHash -----------------------------------
func (e *EthRPCService) GetBlockByHash(ctx context.Context, hashStr string, txDetails bool) (result common.EthGetBlockResult, err error) {
	logger.Infof("eth_getBlockByHash called")

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetBlock", trpc.GetBlockArgs{Hash: tcommon.HexToHash(hashStr)})
	return GetBlockFromTRPCResult(rpcRes, rpcErr, txDetails)
}

func GetBlockFromTRPCResult(rpcRes *rpcc.RPCResponse, rpcErr error, txDetails bool) (result common.EthGetBlockResult, err error) {
	result = common.EthGetBlockResult{}
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetBlockResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		// logger.Infof("transactions count %d \n", len(trpcResult.Txs))
		result.Transactions = make([]interface{}, len(trpcResult.Txs))
		if txDetails {
			var objmap map[string]json.RawMessage
			json.Unmarshal(jsonBytes, &objmap)
			if objmap["transactions"] != nil {
				var txmaps []map[string]json.RawMessage
				json.Unmarshal(objmap["transactions"], &txmaps)
				for i, omap := range txmaps {
					tx := common.EthGetTransactionResult{}
					if types.TxType(trpcResult.Txs[i].Type) == types.TxSmartContract {
						scTx := types.SmartContractTx{}
						json.Unmarshal(omap["raw"], &scTx)
						result.Transactions[i] = scTx
						result.GasUsed = hexutil.Uint64(trpcResult.Txs[i].Receipt.GasUsed)
					} else if types.TxType(trpcResult.Txs[i].Type) == types.TxSend {
						sTx := types.SendTx{}
						json.Unmarshal(omap["raw"], &sTx)
						result.Transactions[i] = sTx
					} else if types.TxType(trpcResult.Txs[i].Type) == types.TxCoinbase {
						cTx := types.CoinbaseTx{}
						json.Unmarshal(omap["raw"], &cTx)
						tx.From = cTx.Proposer.Address
						tx.Gas = hexutil.Uint64(0)
						tx.Value = hexutil.Uint64(cTx.Proposer.Coins.TFuelWei.Uint64())
						tx.Input = "0x"
						data := cTx.Proposer.Signature.ToBytes()
						GetRSVfromSignature(data, &tx)
						result.Transactions[i] = tx
					}
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
	result.Height = hexutil.Uint64(theta_GetBlockResult.Height)
	result.Hash = theta_GetBlockResult.Hash
	result.Parent = theta_GetBlockResult.Parent
	result.Timestamp = hexutil.Uint64(theta_GetBlockResult.Timestamp.ToInt().Uint64())
	result.Proposer = theta_GetBlockResult.Proposer
	result.TxHash = theta_GetBlockResult.TxHash
	result.StateHash = theta_GetBlockResult.StateHash
	for i, tx := range theta_GetBlockResult.Txs {
		if txDetails && (types.TxType(tx.Type) == types.TxSmartContract || types.TxType(tx.Type) == types.TxSend || types.TxType(tx.Type) == types.TxCoinbase) {
			//already handled
		} else {
			result.Transactions[i] = tx.Hash
		}
	}
	result.GasLimit = 20000000
	result.LogsBloom = "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	result.ExtraData = "0x"
	result.Nonce = "0x0000000000000000"
	result.Uncles = []tcommon.Hash{}

	return result, nil
}

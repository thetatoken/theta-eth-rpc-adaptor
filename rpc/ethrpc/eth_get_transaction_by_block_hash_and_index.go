package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/common/hexutil"
	"github.com/thetatoken/theta/ledger/types"
	trpc "github.com/thetatoken/theta/rpc"

	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getTransactionByBlockHashAndIndex -----------------------------------
func (e *EthRPCService) GetTransactionByBlockHashAndIndex(ctx context.Context, hashStr string, txIndexStr string) (result common.EthGetTransactionResult, err error) {
	logger.Infof("GetTransactionByBlockHashAndIndex called")
	txIndex := common.GetHeightByTag(txIndexStr)
	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetBlock", trpc.GetBlockArgs{Hash: tcommon.HexToHash(hashStr)})
	return GetIndexedTransactionFromBlock(rpcRes, rpcErr, txIndex)
}

func GetIndexedTransactionFromBlock(rpcRes *rpcc.RPCResponse, rpcErr error, txIndex tcommon.JSONUint64) (result common.EthGetTransactionResult, err error) {
	result = common.EthGetTransactionResult{}
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetBlockResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		if txIndex >= tcommon.JSONUint64(len(trpcResult.Txs)) {
			return result, fmt.Errorf("transaction index out of range")
		}
		result.TransactionIndex = hexutil.Uint64(txIndex)
		var objmap map[string]json.RawMessage
		json.Unmarshal(jsonBytes, &objmap)
		result.BlockHash = trpcResult.Hash
		result.BlockHeight = hexutil.Uint64(trpcResult.Height)
		result.Nonce = hexutil.Uint64(0)

		if objmap["transactions"] != nil {
			var txmaps []map[string]json.RawMessage
			json.Unmarshal(objmap["transactions"], &txmaps)
			indexedTx := trpcResult.Txs[txIndex]
			omap := txmaps[txIndex]
			result.TxHash = indexedTx.Hash
			if types.TxType(indexedTx.Type) == types.TxSmartContract {
				tx := types.SmartContractTx{}
				json.Unmarshal(omap["raw"], &tx)
				result.From = tx.From.Address
				result.To = tx.To.Address
				result.GasPrice = hexutil.Uint64(tx.GasPrice.Uint64())
				result.Gas = hexutil.Uint64(tx.GasLimit)
				result.Value = hexutil.Uint64(tx.From.Coins.TFuelWei.Uint64())
				result.Input = tx.Data.String()
				data := tx.From.Signature.ToBytes()
				GetRSVfromSignature(data, &result)
			} else if types.TxType(indexedTx.Type) == types.TxSend {
				tx := types.SendTx{}
				json.Unmarshal(omap["raw"], &tx)
				result.From = tx.Inputs[0].Address
				result.To = tx.Outputs[0].Address
				result.Gas = hexutil.Uint64(tx.Fee.TFuelWei.Uint64())
				result.Value = hexutil.Uint64(tx.Inputs[0].Coins.TFuelWei.Uint64())
				data := tx.Inputs[0].Signature.ToBytes()
				GetRSVfromSignature(data, &result)
			}
		}
		return trpcResult, nil
	}
	_, err = common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return result, err
	}
	return result, nil
}

package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/common/hexutil"
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
			if types.TxType(trpcResult.Type) == types.TxSend {
				tx := types.SendTx{}
				json.Unmarshal(objmap["transaction"], &tx)
				trpcResult.Tx = &tx
			}
			if types.TxType(trpcResult.Type) == types.TxSmartContract {
				tx := types.SmartContractTx{}
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
	thetaGetTransactionResult := resultIntf.(trpc.GetTransactionResult)
	result.BlockHash = thetaGetTransactionResult.BlockHash
	result.BlockHeight = thetaGetTransactionResult.BlockHeight
	result.TxHash = thetaGetTransactionResult.TxHash
	if thetaGetTransactionResult.Tx != nil {
		if types.TxType(thetaGetTransactionResult.Type) == types.TxSend {
			tx := thetaGetTransactionResult.Tx.(*types.SendTx)
			result.From = tx.Inputs[0].Address
			result.To = tx.Outputs[0].Address
			result.Gas = (tcommon.JSONBig)(*tx.Fee.TFuelWei)
			result.Value = (tcommon.JSONBig)(*tx.Inputs[0].Coins.TFuelWei)
			data := tx.Inputs[0].Signature.ToBytes()
			GetRSVfromSignature(data, &result)
		}
		if types.TxType(thetaGetTransactionResult.Type) == types.TxSmartContract {
			tx := thetaGetTransactionResult.Tx.(*types.SmartContractTx)
			result.From = tx.From.Address
			result.To = tx.To.Address
			result.GasPrice = tcommon.JSONBig(*tx.GasPrice)
			result.Gas = tcommon.JSONBig(*new(big.Int).SetUint64(tx.GasLimit))
			result.Value = (tcommon.JSONBig)(*tx.From.Coins.TFuelWei)
			result.Input = tx.Data
			data := tx.From.Signature.ToBytes()
			GetRSVfromSignature(data, &result)
		}
	}
	result.TransactionIndex, err = GetTransactionIndex(result.BlockHash, result.TxHash, client)
	if err != nil {
		return result, err
	}
	return result, nil
}

type Tx struct {
	types.Tx `json:"raw"`
	Type     byte         `json:"type"`
	Hash     tcommon.Hash `json:"hash"`
}

func GetTransactionIndex(blockHash tcommon.Hash, transactionHash tcommon.Hash, client *rpcc.RPCClient) (hexutil.Uint64, error) {
	rpcRes, rpcErr := client.Call("theta.GetBlock", trpc.GetBlockArgs{Hash: blockHash})
	if rpcErr != nil {
		return 0, rpcErr
	}
	jsonBytes, err := json.MarshalIndent(rpcRes.Result, "", "    ")
	if err != nil {
		return 0, err
	}
	var objmap map[string]json.RawMessage
	json.Unmarshal(jsonBytes, &objmap)
	var txs []Tx
	if objmap["transactions"] != nil {
		json.Unmarshal(objmap["transactions"], &txs)
	}

	for i, tx := range txs {
		if tx.Hash == transactionHash {
			return hexutil.Uint64(i), nil
		}
	}
	return 0, fmt.Errorf("could not find hash for tx")
}

func GetRSVfromSignature(data []byte, txResult *common.EthGetTransactionResult) error {
	copy(txResult.R[:], data[0:32])
	copy(txResult.S[:], data[32:64])
	txResult.V = tcommon.JSONUint64(data[64])
	return nil
}

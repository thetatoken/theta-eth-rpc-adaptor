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

// ------------------------------- eth_getTransactionReceipt -----------------------------------
func (e *EthRPCService) GetTransactionReceipt(ctx context.Context, hashStr string) (interface{}, error) {
	logger.Infof("eth_getTransactionReceipt called")

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetTransaction", trpc.GetTransactionArgs{Hash: hashStr})
	result := common.EthGetReceiptResult{}

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetTransactionResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		var objmap map[string]json.RawMessage
		json.Unmarshal(jsonBytes, &objmap)
		if objmap["transaction"] != nil {
			if types.TxType(trpcResult.Type) == types.TxSend {
				tx := types.SendTx{}
				json.Unmarshal(objmap["transaction"], &tx)
				// trpcResult.Tx = &tx
				result.From = tx.Inputs[0].Address
				result.To = tx.Outputs[0].Address
			}
			if types.TxType(trpcResult.Type) == types.TxSmartContract {
				tx := types.SmartContractTx{}
				json.Unmarshal(objmap["transaction"], &tx)
				// trpcResult.Tx = &tx
				result.From = tx.From.Address
				result.To = tx.To.Address
			}
		}
		return trpcResult, nil
	}
	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return result, err
	}
	thetaGetTransactionResult := resultIntf.(trpc.GetTransactionResult)
	result.BlockHash = thetaGetTransactionResult.BlockHash
	result.BlockHeight = hexutil.Uint64(thetaGetTransactionResult.BlockHeight)
	result.TxHash = thetaGetTransactionResult.TxHash
	result.GasUsed = hexutil.Uint64(thetaGetTransactionResult.Receipt.GasUsed)
	result.Logs = make([]common.EthLogObj, len(thetaGetTransactionResult.Receipt.Logs))
	for i, log := range thetaGetTransactionResult.Receipt.Logs {
		result.Logs[i] = ThetaLogToEthLog(log)
		result.Logs[i].BlockHash = result.BlockHash
		result.Logs[i].BlockHeight = result.BlockHeight
		result.Logs[i].TxHash = result.TxHash
		result.Logs[i].LogIndex = hexutil.Uint64(i)
	}
	//TODO: handle logIndex & TransactionIndex of logs
	result.TransactionIndex, result.CumulativeGasUsed, err = GetTransactionIndexAndCumulativeGasUsed(result.BlockHash, result.TxHash, result.Logs, client)
	if err != nil {
		return nil, err
	}
	result.Status = 1
	return result, nil
}

func GetTransactionIndexAndCumulativeGasUsed(blockHash tcommon.Hash, transactionHash tcommon.Hash, logs []common.EthLogObj, client *rpcc.RPCClient) (hexutil.Uint64, hexutil.Uint64, error) {
	rpcRes, rpcErr := client.Call("theta.GetBlock", trpc.GetBlockArgs{Hash: blockHash})
	if rpcErr != nil {
		return 0, 0, rpcErr
	}
	jsonBytes, err := json.MarshalIndent(rpcRes.Result, "", "    ")
	if err != nil {
		return 0, 0, err
	}
	var objmap map[string]json.RawMessage
	json.Unmarshal(jsonBytes, &objmap)
	var txs []common.Tx
	if objmap["transactions"] != nil {
		json.Unmarshal(objmap["transactions"], &txs)
	}
	var cumulativeGas hexutil.Uint64
	var logIndex int
	for i, tx := range txs {
		if types.TxType(tx.Type) == types.TxSmartContract {
			cumulativeGas += hexutil.Uint64(tx.Receipt.GasUsed)
			if tx.Hash != transactionHash {
				logIndex += len(tx.Receipt.Logs)
			}
		}
		if tx.Hash == transactionHash {
			for j, _ := range logs {
				log := &logs[j]
				log.LogIndex = hexutil.Uint64(logIndex)
				log.TransactionIndex = hexutil.Uint64(i)
				logger.Infof("jlog2 i is %d, log.TransactionIndex is %d, logs[i].TransactionIndex is %d \n", i, log.TransactionIndex, logs[j].TransactionIndex)
			}
			return hexutil.Uint64(i), cumulativeGas, nil
		}
	}
	return 0, 0, fmt.Errorf("could not find hash for tx")
}

func ThetaLogToEthLog(log *types.Log) common.EthLogObj {
	result := common.EthLogObj{}
	result.Address = log.Address
	result.Data = log.Data
	result.Topics = log.Topics
	return result
}

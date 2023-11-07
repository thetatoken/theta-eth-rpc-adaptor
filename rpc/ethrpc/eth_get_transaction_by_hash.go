package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/common/hexutil"
	"github.com/thetatoken/theta/ledger/types"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getTransactionByHash -----------------------------------
func (e *EthRPCService) GetTransactionByHash(ctx context.Context, hashStr string) (result common.EthGetTransactionResult, err error) {
	logger.Infof("eth_getTransactionByHash called, txHash: %v", hashStr)

	result = common.EthGetTransactionResult{}
	var resultIntf interface{}
	var thetaGetTransactionResult trpc.GetTransactionResult

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	maxRetry := 5
	for i := 0; i < maxRetry; i++ { // It might take some time for a block to be finalized, retry a few times
		rpcRes, rpcErr := client.Call("theta.GetTransaction", trpc.GetTransactionArgs{Hash: hashStr})

		if rpcErr != nil {
			logger.Warnf("eth_getTransactionByHash failed, err: %v", rpcErr)
		}

		parse := func(jsonBytes []byte) (interface{}, error err) {
			err = fmt.Errorf("Failed to parse the transaction")
			defer func() {
				recover() // recover from parsing exception.
			}()

			trpcResult := trpc.GetTransactionResult{}
			json.Unmarshal(jsonBytes, &trpcResult)
			if (trpcResult.BlockHash == tcommon.Hash{}) {
				return trpcResult, nil // Thet tx is not finalized yet. Just let the outer loop retry. Otherwise, tx unmashal might crash if we continue the parsing
			}
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
			logger.Infof("eth_getTransactionByHash EvmRet: %+v\n", trpcResult.Receipt.EvmRet)
			return trpcResult, nil
		}
		resultIntf, err = common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
		if err != nil {
			logger.Warnf("eth_getTransactionByHash failed, err: %v", err)
			return result, err
		}

		thetaGetTransactionResult = resultIntf.(trpc.GetTransactionResult)
		if (thetaGetTransactionResult.BlockHash != tcommon.Hash{}) {
			break
		}

		time.Sleep(blockInterval) // one block duration
	}

	result.BlockHash = thetaGetTransactionResult.BlockHash
	result.BlockHeight = hexutil.Uint64(thetaGetTransactionResult.BlockHeight)

	trimmedHashStr := hashStr
	if strings.HasPrefix(hashStr, "0x") {
		trimmedHashStr = hashStr[2:]
	}
	txHash, _ := hex.DecodeString(trimmedHashStr)
	result.TxHash = tcommon.BytesToHash(txHash) // For ethers.js compatibility, need to return the ETH tx hash (i.e. the query parameter)

	logger.Infof("eth_getTransactionByHash, hashStr: %v, result.TxHash: %v", hashStr, result.TxHash.Hex())

	nativeTxHash := thetaGetTransactionResult.TxHash // need use native tx hash to find the tx index, instead of the ETH tx hash
	if thetaGetTransactionResult.Tx != nil {
		if types.TxType(thetaGetTransactionResult.Type) == types.TxSend {
			tx := thetaGetTransactionResult.Tx.(*types.SendTx)
			result.From = tx.Inputs[0].Address
			if (tx.Outputs[0].Address == tcommon.Address{}) {
				result.To = nil // conform to ETH standard
			} else {
				result.To = &tx.Outputs[0].Address
			}
			result.Gas = hexutil.Uint64(tx.Fee.TFuelWei.Uint64())
			result.Value = "0x" + tx.Inputs[0].Coins.TFuelWei.Text(16)
			data := tx.Inputs[0].Signature.ToBytes()
			result.Nonce = hexutil.Uint64(tx.Inputs[0].Sequence) - 1 // off-by-one: Ethereum's account nonce starts from 0, while Theta's account sequnce starts from 1
			GetRSVfromSignature(data, &result)
		}
		if types.TxType(thetaGetTransactionResult.Type) == types.TxSmartContract {
			tx := thetaGetTransactionResult.Tx.(*types.SmartContractTx)
			result.From = tx.From.Address
			if (tx.To.Address == tcommon.Address{}) {
				result.To = nil // conform to ETH standard
			} else {
				result.To = &tx.To.Address
			}
			result.GasPrice = "0x" + tx.GasPrice.Text(16)
			result.Gas = hexutil.Uint64(tx.GasLimit)
			result.Value = "0x" + tx.From.Coins.TFuelWei.Text(16)
			//result.Input = tx.Data.String()
			result.Input = "0x" + hex.EncodeToString(tx.Data)
			data := tx.From.Signature.ToBytes()
			result.Nonce = hexutil.Uint64(tx.From.Sequence) - 1 // off-by-one: Ethereum's account nonce starts from 0, while Theta's account sequnce starts from 1
			GetRSVfromSignature(data, &result)
		}
	}
	result.TransactionIndex, err = GetTransactionIndex(result.BlockHash, nativeTxHash, client)
	if err != nil {
		return result, err
	}

	//resultJsonBytes, _ := json.MarshalIndent(result, "", "    ")
	//logger.Infof("eth_getTransactionByHash, result: %v", string(resultJsonBytes))

	return result, nil
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
	var txs []common.Tx
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
	txResult.V = hexutil.Uint64(data[64])
	return nil
}

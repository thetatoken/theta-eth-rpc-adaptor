package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	rpcc "github.com/ybbus/jsonrpc"

	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/ledger/types"

	// "github.com/thetatoken/theta/ledger/types"
	trpc "github.com/thetatoken/theta/rpc"
)

type TxTmp struct {
	Tx   json.RawMessage `json:"raw"`
	Type byte            `json:"type"`
	Hash tcommon.Hash    `json:"hash"`
}

// ------------------------------- eth_gasPrice -----------------------------------

func (e *EthRPCService) GasPrice(ctx context.Context) (result string, err error) {
	logger.Infof("eth_gasPrice called")

	currentHeight, err := common.GetCurrentHeight()

	if err != nil {
		return "", err
	}

	// fmt.Printf("currentHeight: %v\n", currentHeight)
	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetBlockByHeight", trpc.GetBlockByHeightArgs{Height: currentHeight})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetBlockResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		var objmap map[string]json.RawMessage
		json.Unmarshal(jsonBytes, &objmap)
		if objmap["transactions"] != nil {
			//TODO: handle other types
			txs := []trpc.Tx{}
			tmpTxs := []TxTmp{}
			json.Unmarshal(objmap["transactions"], &tmpTxs)
			for _, tx := range tmpTxs {
				newTx := trpc.Tx{}
				newTx.Type = tx.Type
				newTx.Hash = tx.Hash
				if types.TxType(tx.Type) == types.TxSmartContract {
					transaction := types.SmartContractTx{}
					json.Unmarshal(tx.Tx, &transaction)
					// fmt.Printf("transaction: %+v\n", transaction)
					newTx.Tx = &transaction
				}
				txs = append(txs, newTx)
			}
			trpcResult.Txs = txs
		}
		return trpcResult, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	thetaGetBlockResult, ok := resultIntf.(trpc.GetBlockResult)
	if !ok {
		return "", fmt.Errorf("failed to convert GetBlockResult")
	}
	totalGasPrice := big.NewInt(0)
	count := 0
	for _, tx := range thetaGetBlockResult.Txs {
		if types.TxType(tx.Type) != types.TxSmartContract {
			continue
		}
		if tx.Tx != nil {
			transaction := tx.Tx.(*types.SmartContractTx)
			count++
			totalGasPrice = new(big.Int).Add(transaction.GasPrice, totalGasPrice)
		}
	}
	gasPrice := big.NewInt(1000000000000)
	if count != 0 {
		gasPrice = new(big.Int).Div(totalGasPrice, big.NewInt(int64(count)))
	}
	fmt.Printf("gasPrice: %v\n", gasPrice)
	result = "0x" + gasPrice.Text(16)
	return result, nil
}

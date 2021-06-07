package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	tcommon "github.com/thetatoken/theta/common"
	hexutil "github.com/thetatoken/theta/common/hexutil"
	"github.com/thetatoken/theta/ledger/types"
	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

type EthGetLogsArgs struct {
	FromBlock string          `json:"fromBlock"`
	ToBlock   string          `json:"toBlock"`
	Address   tcommon.Address `json:"address"`
	Topics    []tcommon.Hash  `json:"topics"`
	Blockhash tcommon.Hash    `json:"blockhash"`
}

type EthGetLogsResult struct {
	Removed          bool            `json: "removed"`
	LogIndex         string          `json: "logIndex"`
	TransactionIndex string          `json: "transactionIndex"`
	TransactionHash  tcommon.Hash    `json: "transactionHash"`
	BlockHash        tcommon.Hash    `json: "blockHash"`
	BlockNumber      string          `json: "blockNumber"`
	Address          tcommon.Address `json: "address"`
	Data             []byte          `json:"data"`
	Topics           []tcommon.Hash  `json:"topics"`
}

// ------------------------------- eth_getLogs -----------------------------------

func (e *EthRPCService) GetLogs(ctx context.Context, args EthGetLogsArgs) (result []EthGetLogsResult, err error) {
	logger.Infof("eth_getLogs called")

	blocks := []*trpc.GetBlockResultInner{}
	fmt.Printf("block hash: %v\n", args.Blockhash)
	fmt.Printf("block hash.hex: %v\n", args.Blockhash.Hex())
	if args.Blockhash.Hex() != "0x0000000000000000000000000000000000000000000000000000000000000000" {
		fmt.Printf("block hash: %v\n", args.Blockhash)
		client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
		rpcRes, rpcErr := client.Call("theta.GetBlock", trpc.GetBlockArgs{Hash: args.Blockhash})

		parse := func(jsonBytes []byte) (interface{}, error) {
			trpcResult := trpc.GetBlockResult{GetBlockResultInner: &trpc.GetBlockResultInner{}}
			json.Unmarshal(jsonBytes, &trpcResult)
			var objmap map[string]json.RawMessage
			json.Unmarshal(jsonBytes, &objmap)
			if objmap["transactions"] != nil {
				txs := []trpc.Tx{}
				json.Unmarshal(objmap["transactions"], &txs)
				trpcResult.Txs = txs
			}
			return trpcResult, nil
		}

		resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
		if err != nil {
			return result, err
		}
		block := resultIntf.(trpc.GetBlockResult)
		blocks = append(blocks, block.GetBlockResultInner)
	} else {
		blockStart := tcommon.JSONUint64(0)
		if args.FromBlock != "" {
			blockStart = common.GetHeightByTag(args.FromBlock)
			fmt.Printf("blockStart: %v\n", blockStart)
		}

		blockEnd := tcommon.JSONUint64(0)
		if args.ToBlock != "" {
			blockEnd = common.GetHeightByTag(args.ToBlock)
			if blockEnd == 0 {
				currentHeight, err := common.GetCurrentHeight()
				if err != nil {
					return result, err
				}
				blockEnd = currentHeight
			}
		} else {
			currentHeight, err := common.GetCurrentHeight()
			if err != nil {
				return result, err
			}
			blockEnd = currentHeight
		}
		fmt.Printf("blockEnd: %v\n", blockEnd)

		client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
		rpcRes, rpcErr := client.Call("theta.GetBlocksByRange", trpc.GetBlocksByRangeArgs{Start: blockStart, End: blockEnd})

		parse := func(jsonBytes []byte) (interface{}, error) {
			trpcResult := trpc.GetBlocksResult{}
			json.Unmarshal(jsonBytes, &trpcResult)
			return trpcResult, nil
		}

		resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
		if err != nil {
			return result, err
		}
		blocks = resultIntf.(trpc.GetBlocksResult)
	}
	fmt.Printf("blocks: %v\n", blocks)

	for _, block := range blocks {
		fmt.Printf("txs: %+v\n", block.Txs)
		for txIndex, tx := range block.Txs {
			if types.TxType(tx.Type) != types.TxSmartContract {
				continue
			}
			// if tx.Receipt != nil {
			receipt := *tx.Receipt
			fmt.Printf("receipt: %v\n", receipt)
			fmt.Printf("receipt.Logs: %v\n", receipt.Logs)
			fmt.Printf("args.topics: %v\n", args.Topics)
			for logIndex, log := range receipt.Logs {
				if len(args.Topics) > 0 {
					for _, topic := range log.Topics {
						for _, t := range args.Topics {
							fmt.Printf("topic: %v\n", topic)
							fmt.Printf("t: %v\n", t)
							if topic == t {
								res := EthGetLogsResult{}
								res.Removed = false
								res.LogIndex = common.Int2hex2str(logIndex)
								res.TransactionIndex = common.Int2hex2str(txIndex)
								res.TransactionHash = tx.Hash
								res.BlockHash = block.Hash
								res.BlockNumber = hexutil.EncodeUint64(uint64(block.Height))
								res.Address = receipt.ContractAddress
								res.Data = log.Data
								res.Topics = log.Topics
								result = append(result, res)
							}
						}
					}
				} else {
					res := EthGetLogsResult{}
					res.Removed = false
					res.LogIndex = common.Int2hex2str(logIndex)
					res.TransactionIndex = common.Int2hex2str(txIndex)
					res.TransactionHash = tx.Hash
					res.BlockHash = block.Hash
					res.BlockNumber = hexutil.EncodeUint64(uint64(block.Height))
					res.Address = receipt.ContractAddress
					res.Data = log.Data
					res.Topics = log.Topics
					result = append(result, res)
				}

				// }
			}
		}
	}
	return result, nil
}

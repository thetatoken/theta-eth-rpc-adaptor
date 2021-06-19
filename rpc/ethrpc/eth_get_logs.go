package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	logger.Infof("eth_getLogs called, fromBlock: %v, toBlock: %v, address: %v, blockHash: %v, topics: %v\n",
		args.FromBlock, args.ToBlock, args.Address, args.Blockhash.Hex(), args.Topics)

	maxRetry := 5
	blocks := []*trpc.GetBlockResultInner{}
	if args.Blockhash.Hex() != "0x0000000000000000000000000000000000000000000000000000000000000000" {
		client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())

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

		var block trpc.GetBlockResult
		for i := 0; i < maxRetry; i++ { // It might take some time for a tx to be finalized, retry a few times
			if i == maxRetry {
				return []EthGetLogsResult{}, fmt.Errorf("failed to retrieve block %v", args.Blockhash.Hex())
			}

			rpcRes, rpcErr := client.Call("theta.GetBlock", trpc.GetBlockArgs{Hash: args.Blockhash})
			resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
			if err == nil {
				block = resultIntf.(trpc.GetBlockResult)
				break
			}

			logger.Warnf("eth_getLogs, theta.GetBlock returned error: %v", err)
			time.Sleep(blockInterval) // one block duration
		}

		blocks = append(blocks, block.GetBlockResultInner)
	} else {
		blockStart := tcommon.JSONUint64(0)
		if args.FromBlock != "" {
			blockStart = common.GetHeightByTag(args.FromBlock)
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

		if blockStart > blockEnd {
			tmp := blockStart
			blockStart = blockEnd
			blockEnd = tmp
		}
		blockStart -= 2 // Theta requires two consecutive committed blocks for finalization

		logger.Infof("blockStart: %v, blockEnd: %v", blockStart, blockEnd)

		client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())

		parse := func(jsonBytes []byte) (interface{}, error) {
			trpcResult := trpc.GetBlocksResult{}
			json.Unmarshal(jsonBytes, &trpcResult)
			return trpcResult, nil
		}

		for i := 0; i < maxRetry; i++ { // It might take some time for a tx to be finalized, retry a few times
			if i == maxRetry {
				return []EthGetLogsResult{}, fmt.Errorf("failed to retrieve blocks from %v to %v", blockStart, blockEnd)
			}

			rpcRes, rpcErr := client.Call("theta.GetBlocksByRange", trpc.GetBlocksByRangeArgs{Start: blockStart, End: blockEnd})
			resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
			blocks = resultIntf.(trpc.GetBlocksResult)
			if err == nil && len(blocks) > 0 {
				break
			}

			logger.Warnf("eth_getLogs, theta.GetBlocksByRange returned error: %v", err)
			time.Sleep(blockInterval) // one block duration
		}

	}
	logger.Infof("blocks: %v\n", blocks)

	for _, block := range blocks {
		logger.Infof("txs: %+v\n", block.Txs)
		for txIndex, tx := range block.Txs {
			if types.TxType(tx.Type) != types.TxSmartContract {
				continue
			}

			if tx.Receipt == nil {
				logger.Errorf("No receipt for tx: %v", tx.Hash.Hex())
				continue
			}

			// if tx.Receipt != nil {
			receipt := *tx.Receipt
			logger.Infof("receipt: %v\n", receipt)
			logger.Infof("receipt.Logs: %v\n", receipt.Logs)
			logger.Infof("args.topics: %v\n", args.Topics)
			for logIndex, log := range receipt.Logs {
				if len(args.Topics) > 0 {
					for _, topic := range log.Topics {
						for _, t := range args.Topics {
							logger.Infof("topic: %v\n", topic)
							logger.Infof("t: %v\n", t)
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

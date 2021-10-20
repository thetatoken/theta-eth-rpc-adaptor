package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	tcommon "github.com/thetatoken/theta/common"
	hexutil "github.com/thetatoken/theta/common/hexutil"
	"github.com/thetatoken/theta/ledger/types"
	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// type EthGetLogsArgs struct {
// 	FromBlock string          `json:"fromBlock"`
// 	ToBlock   string          `json:"toBlock"`
// 	Address   tcommon.Address `json:"address"`
// 	Topics    []tcommon.Hash  `json:"topics"`
// 	Blockhash tcommon.Hash    `json:"blockhash"`
// }

type EthGetLogsArgs struct {
	FromBlock string        `json:"fromBlock"`
	ToBlock   string        `json:"toBlock"`
	Address   interface{}   `json:"address"`
	Topics    []interface{} `json:"topics"`
	Blockhash tcommon.Hash  `json:"blockhash"`
}

type EthGetLogsResult struct {
	//Removed          bool            `json:"removed"`
	LogIndex         string          `json:"logIndex"`
	TransactionIndex string          `json:"transactionIndex"`
	TransactionHash  tcommon.Hash    `json:"transactionHash"`
	BlockHash        tcommon.Hash    `json:"blockHash"`
	BlockNumber      string          `json:"blockNumber"`
	Address          tcommon.Address `json:"address"`
	Data             string          `json:"data"`
	Topics           []tcommon.Hash  `json:"topics"`
	Type             string          `json:"type"`
}

// ------------------------------- eth_getLogs -----------------------------------

//
// Reference: https://docs.alchemy.com/alchemy/guides/eth_getlogs
//
func (e *EthRPCService) GetLogs(ctx context.Context, args EthGetLogsArgs) (result []EthGetLogsResult, err error) {
	logger.Infof("eth_getLogs called, fromBlock: %v, toBlock: %v, address: %v, blockHash: %v, topics: %v\n",
		args.FromBlock, args.ToBlock, args.Address, args.Blockhash.Hex(), args.Topics)

	result = []EthGetLogsResult{}

	addresses, err := parseAddresses(args.Address)
	if err != nil {
		return result, err
	}

	topicsFilter, err := parseTopicsFilter(args.Topics)
	if err != nil {
		return result, err
	}

	parse := func(jsonBytes []byte) (interface{}, error) {
		//logger.Infof("eth_getLogs.parse, jsonBytes: %v", string(jsonBytes))

		trpcResult := common.ThetaGetBlockResult{ThetaGetBlockResultInner: &common.ThetaGetBlockResultInner{}}
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

	maxRetry := 5
	blocks := []*common.ThetaGetBlockResultInner{}
	if args.Blockhash.Hex() != "0x0000000000000000000000000000000000000000000000000000000000000000" {
		client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())

		var block common.ThetaGetBlockResult
		for i := 0; i < maxRetry; i++ { // It might take some time for a tx to be finalized, retry a few times
			if i == maxRetry {
				return []EthGetLogsResult{}, fmt.Errorf("failed to retrieve block %v", args.Blockhash.Hex())
			}

			rpcRes, rpcErr := client.Call("theta.GetBlock", trpc.GetBlockArgs{Hash: args.Blockhash})
			resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
			if err == nil {
				block = resultIntf.(common.ThetaGetBlockResult)
				break
			}

			logger.Warnf("eth_getLogs, theta.GetBlock returned error: %v", err)
			time.Sleep(blockInterval) // one block duration
		}

		if block.ThetaGetBlockResultInner != nil {
			blocks = append(blocks, block.ThetaGetBlockResultInner)
		}
	} else {
		currentHeight, err := common.GetCurrentHeight()
		if err != nil {
			return result, err
		}

		blockStart := currentHeight
		if args.FromBlock != "" {
			blockStart = common.GetHeightByTag(args.FromBlock)
			if blockStart == math.MaxUint64 {
				blockStart = currentHeight
			}
		}

		blockEnd := currentHeight
		if args.ToBlock != "" {
			blockEnd = common.GetHeightByTag(args.ToBlock)
			if blockEnd == math.MaxUint64 {
				blockEnd = currentHeight
			}
		}

		if blockStart > blockEnd {
			tmp := blockStart
			blockStart = blockEnd
			blockEnd = tmp
		}
		// if blockStart >= 2 {
		// 	blockStart -= 2 // Theta requires two consecutive committed blocks for finalization
		// }

		logger.Infof("blockStart: %v, blockEnd: %v", blockStart, blockEnd)

		client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
		for i := 0; i < maxRetry; i++ { // It might take some time for a tx to be finalized, retry a few times
			if i == maxRetry {
				return []EthGetLogsResult{}, fmt.Errorf("failed to retrieve blocks from %v to %v", blockStart, blockEnd)
			}

			success := true
			for h := uint64(blockStart); h <= uint64(blockEnd); h++ {
				rpcRes, rpcErr := client.Call("theta.GetBlockByHeight", trpc.GetBlockByHeightArgs{Height: tcommon.JSONUint64(h)})
				resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
				if err != nil {
					success = false
					blocks = []*common.ThetaGetBlockResultInner{}
					logger.Warnf("eth_getLogs, theta.GetBlocksByHeight returned error: %v", err)
					break
				}

				block := resultIntf.(common.ThetaGetBlockResult)
				if block.ThetaGetBlockResultInner != nil {
					blocks = append(blocks, block.ThetaGetBlockResultInner)
				}
			}

			if success {
				break
			}

			time.Sleep(blockInterval) // one block duration
		}
	}
	//logger.Infof("blocks: %v\n", blocks)

	filterByAddress := !(len(addresses) == 0 || (len(addresses) == 1) && (addresses[0] == tcommon.Address{}))

	logger.Debugf("filterByAddress: %v, addresses: %v", filterByAddress, addresses)

	for _, block := range blocks {
		logger.Debugf("txs: %+v\n", block.Txs)
		for txIndex, tx := range block.Txs {
			if types.TxType(tx.Type) != types.TxSmartContract {
				continue
			}

			if tx.Receipt == nil {
				logger.Errorf("No receipt for tx: %v", tx.Hash.Hex())
				continue
			}

			receipt := *tx.Receipt
			logger.Debugf("receipt: %v\n", receipt)
			logger.Debugf("receipt.Logs: %v\n", receipt.Logs)
			logger.Debugf("topicsFilter: %v\n", topicsFilter)

			for logIndex, log := range receipt.Logs {

				logger.Debugf("filterByAddress: %v, addresses: %v, log.Address: %v, addrMatch: %v", filterByAddress, addresses, log.Address, addressMatch(addresses, log.Address))
				if filterByAddress && !addressMatch(addresses, log.Address) {
					continue
				}

				if topicsMatch(topicsFilter, log) {
					res := EthGetLogsResult{}
					res.Type = "mined"
					res.LogIndex = common.Int2hex2str(logIndex)
					res.TransactionIndex = common.Int2hex2str(txIndex)
					res.TransactionHash = tx.Hash
					res.BlockHash = block.Hash
					res.BlockNumber = hexutil.EncodeUint64(uint64(block.Height))
					res.Address = receipt.ContractAddress
					res.Data = "0x" + hex.EncodeToString(log.Data)
					res.Topics = log.Topics
					result = append(result, res)
				}
			}
		}
	}

	resultJson, _ := json.Marshal(result)
	logger.Infof("eth_getLogs, result: %v", string(resultJson))

	return result, nil
}

func parseAddresses(argsAddress interface{}) ([]tcommon.Address, error) {
	// Some clients may call eth_getLogs with a single address or a list of addresses
	var addresses []tcommon.Address
	addrVal := argsAddress
	switch addrVal.(type) {
	case string:
		address := tcommon.HexToAddress(addrVal.(string))
		addresses = append(addresses, address)
	case []interface{}:
		for _, vi := range addrVal.([]interface{}) {
			val := vi.(string)
			address := tcommon.HexToAddress(val)
			addresses = append(addresses, address)
		}
	default:
		//return []tcommon.Address{}, fmt.Errorf("invalid args.Address type: %v", argsAddress)
		return []tcommon.Address{}, nil
	}

	return addresses, nil
}

func parseTopicsFilter(argsTopics []interface{}) ([][]tcommon.Hash, error) {
	// some clients, e.g. the Graph calls the eth_getLogs methods with topics formatted as a list of list, i.e. topics : [[0x..., 0x....]]
	// needs special handling to convert it into a list of hashs
	var topicsFilter [][]tcommon.Hash
	for _, val := range argsTopics {
		switch val.(type) {
		case string:
			topic := tcommon.HexToHash(val.(string))
			topicsFilter = append(topicsFilter, []tcommon.Hash{topic})
		case []interface{}:
			for _, item := range val.([]interface{}) {
				topic := tcommon.HexToHash(item.(string))
				topicsFilter = append(topicsFilter, []tcommon.Hash{topic})
			}
		case [][]interface{}:
			for _, itemList := range val.([][]interface{}) {
				topicList := []tcommon.Hash{}
				for _, item := range itemList {
					topic := tcommon.HexToHash(item.(string))
					topicList = append(topicList, topic)
				}
				topicsFilter = append(topicsFilter, topicList)
			}
		default:
			return [][]tcommon.Hash{}, fmt.Errorf("invalid args.Topics type: %v", argsTopics)
		}
	}

	return topicsFilter, nil
}

func addressMatch(addresses []tcommon.Address, contractAddress tcommon.Address) bool {
	for _, address := range addresses {
		if address == contractAddress {
			return true
		}
	}
	return false
}

// Reference: https://docs.alchemy.com/alchemy/guides/eth_getlogs#a-note-on-specifying-topic-filters
func topicsMatch(topicsFilter [][]tcommon.Hash, log *types.Log) bool {
	numFilters := len(topicsFilter)
	numLogTopics := len(log.Topics)
	for i := 0; i < numFilters; i++ {
		if i >= numLogTopics {
			break
		}

		logTopic := log.Topics[i]
		if !topicIncludedIn(logTopic, topicsFilter[i]) {
			return false
		}
	}

	return true
}

func topicIncludedIn(logTopic tcommon.Hash, topicList []tcommon.Hash) bool {
	if len(topicList) == 0 {
		return true
	}

	for _, topic := range topicList {
		if len(topic) == 0 {
			return true
		}

		if (topic == tcommon.Hash{}) {
			return true
		}

		if topic == logTopic {
			return true
		}
	}

	return false
}

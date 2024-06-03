package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

// Reference: https://docs.alchemy.com/alchemy/guides/eth_getlogs
func (e *EthRPCService) GetLogs(ctx context.Context, args EthGetLogsArgs) (result []EthGetLogsResult, err error) {
	logger.Infof("eth_getLogs called, fromBlock: %v, toBlock: %v, address: %v, blockHash: %v, topics: %v\n",
		args.FromBlock, args.ToBlock, args.Address, args.Blockhash.Hex(), args.Topics)

	start := time.Now()

	result = []EthGetLogsResult{}

	addresses, err := parseAddresses(args.Address)
	if err != nil {
		return result, err
	}

	topicsFilter, err := parseTopicsFilter(args.Topics)
	if err != nil {
		return result, err
	}

	maxRetry := 5
	blocks := []*common.ThetaGetBlockResultInner{}
	if args.Blockhash.Hex() != "0x0000000000000000000000000000000000000000000000000000000000000000" {
		err = retrieveBlockByHash(args.Blockhash, &blocks, maxRetry)
	} else {
		err = retrieveBlocksByRange(args.FromBlock, args.ToBlock, &blocks, maxRetry)
	}
	if err != nil {
		return result, err
	}

	queryBlocksTime := time.Since(start)
	start = time.Now()

	filterByAddress := !(len(addresses) == 0 || (len(addresses) == 1) && (addresses[0] == tcommon.Address{}))
	logger.Debugf("filterByAddress: %v, addresses: %v", filterByAddress, addresses)

	extractLogs(addresses, topicsFilter, filterByAddress, blocks, &result)

	resultJson, _ := json.Marshal(result)
	logger.Infof("eth_getLogs, queryBlocksTime: %v, result: %v", queryBlocksTime, string(resultJson))

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
	var topicsFilter [][]tcommon.Hash
	for _, val := range argsTopics {
		switch val.(type) {
		case string:
			topic := tcommon.HexToHash(val.(string))
			topicsFilter = append(topicsFilter, []tcommon.Hash{topic})
		case []interface{}:
			topicList := []tcommon.Hash{}
			for _, item := range val.([]interface{}) {
				topic := tcommon.HexToHash(item.(string))
				topicList = append(topicList, topic)
			}
			topicsFilter = append(topicsFilter, topicList)
		case nil:
			break
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
		logger.Debugf("topic: %v, logTopic: %v, topic == logTopic: %v\n", topic.Hex(), logTopic.Hex(), topic == logTopic)

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

func retrieveBlockByHash(blockhash tcommon.Hash, blocks *[](*common.ThetaGetBlockResultInner), maxRetry int) (err error) {
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

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())

	var block common.ThetaGetBlockResult
	for i := 0; i < maxRetry; i++ { // It might take some time for a tx to be finalized, retry a few times
		if i == maxRetry {
			return fmt.Errorf("failed to retrieve block %v", blockhash.Hex())
		}

		rpcRes, rpcErr := client.Call("theta.GetBlock", trpc.GetBlockArgs{Hash: blockhash})
		resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
		if err == nil {
			block = resultIntf.(common.ThetaGetBlockResult)
			break
		}

		logger.Warnf("eth_getLogs, theta.GetBlock returned error: %v", err)
		time.Sleep(blockInterval) // one block duration
	}

	if block.ThetaGetBlockResultInner != nil {
		(*blocks) = append((*blocks), block.ThetaGetBlockResultInner)
	}

	return nil
}

func retrieveBlocksByRange(fromBlock string, toBlock string, blocks *[](*common.ThetaGetBlockResultInner), maxRetry int) (err error) {
	parse := func(jsonBytes []byte) (interface{}, error) {
		//logger.Infof("eth_getLogs.parse, jsonBytes: %v", string(jsonBytes))

		trpcResult := common.ThetaGetBlocksResult{}

		var objList []json.RawMessage
		json.Unmarshal(jsonBytes, &objList)
		for _, blockJsonBytes := range objList {
			getBlockResult := common.ThetaGetBlockResult{ThetaGetBlockResultInner: &common.ThetaGetBlockResultInner{}}
			json.Unmarshal(blockJsonBytes, &getBlockResult)
			var objmap map[string]json.RawMessage
			json.Unmarshal(blockJsonBytes, &objmap)
			if objmap["transactions"] != nil {
				txs := []trpc.Tx{}
				json.Unmarshal(objmap["transactions"], &txs)
				getBlockResult.Txs = txs
			}
			trpcResult = append(trpcResult, getBlockResult.ThetaGetBlockResultInner)
		}

		return trpcResult, nil
	}

	currentHeight, err := common.GetCurrentHeight()
	if err != nil {
		return err
	}

	blockStart := currentHeight
	if fromBlock != "" {
		blockStart = common.GetHeightByTag(fromBlock)
		if blockStart == math.MaxUint64 {
			blockStart = currentHeight
		}
	}

	blockEnd := currentHeight
	if toBlock != "" {
		blockEnd = common.GetHeightByTag(toBlock)
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

	blockRangeLimit := viper.GetUint64(common.CfgQueryGetLogsBlockRange)
	queryBlockRange := blockEnd - blockStart + 1

	logger.Infof("eth_getLogs, theta.GetBlocksByRange blockStart: %v, blockEnd: %v, blockRange: %v", blockStart, blockEnd, queryBlockRange)
	if queryBlockRange > tcommon.JSONUint64(blockRangeLimit) {
		logger.Infof("queried block range too large")
		return fmt.Errorf("eth_getLogs, theta.GetBlocksByRange block range too large, we currently allow querying for at most %v blocks at a time (start: %v, end: %v)", blockRangeLimit, blockStart, blockEnd)
	}

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	for i := 0; i < maxRetry; i++ { // It might take some time for a tx to be finalized, retry a few times
		if i == maxRetry {
			return fmt.Errorf("eth_getLogs, theta.GetBlocksByRange failed to retrieve blocks from %v to %v", blockStart, blockEnd)
		}

		rpcRes, rpcErr := client.Call("theta.GetBlocksByRange", trpc.GetBlocksByRangeArgs{Start: tcommon.JSONUint64(blockStart), End: tcommon.JSONUint64(blockEnd)})
		rpcResJson, err := json.Marshal(rpcRes)
		if err != nil {
			logger.Warnf("eth_getLogs, theta.GetBlocksByRange returned error: %v", err)
		}
		if len(string(rpcResJson)) > viper.GetInt(common.CfgLogRpcResponseSizeThreshold) {
			logger.WithFields(log.Fields{
				"rpc":            "eth_getLogs",
				"func":           "theta.GetBlocksByRange",
				"responseLength": len(string(rpcResJson)),
				"blockRange":     queryBlockRange,
				"blockStart":     blockStart,
				"blockEnd":       blockEnd,
			}).Infof("resp size")
		}
		resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
		if err != nil {
			blocks = &[]*common.ThetaGetBlockResultInner{}
			logger.Warnf("eth_getLogs, theta.GetBlocksByRange returned error: %v", err)
			time.Sleep(blockInterval) // one block duration
			continue
		}

		getBlocksRes := resultIntf.(common.ThetaGetBlocksResult)
		for _, block := range getBlocksRes {
			if block != nil {
				*blocks = append((*blocks), block)
			}
		}

		break
	}

	return nil
}

func extractLogs(addresses []tcommon.Address, topicsFilter [][]tcommon.Hash, filterByAddress bool, blocks [](*common.ThetaGetBlockResultInner), result *([]EthGetLogsResult)) {
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
					res.Address = log.Address
					res.Data = "0x" + hex.EncodeToString(log.Data)
					res.Topics = log.Topics
					*result = append(*result, res)
				}
			}
		}
	}
}

package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	tcommon "github.com/thetatoken/theta/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

type Bytes8 [8]byte
type Bytes256 [256]byte

type eth_GetBlockResult struct {
	Height    tcommon.JSONUint64 `json:"number"`
	Hash      tcommon.Hash       `json:"hash"`
	Parent    tcommon.Hash       `json:"parentHash"`
	Timestamp *tcommon.JSONBig   `json:"timestamp"`
	Proposer  tcommon.Address    `json:"miner"`
	TxHash    tcommon.Hash       `json:"transactionsRoot"`
	StateHash tcommon.Hash       `json:"stateRoot"`

	Nonce           Bytes8             `json:"nonce"`
	Sha3Uncles      tcommon.Hash       `json:"sha3Uncles"`
	LogsBloom       Bytes256           `json:"logsBloom"`
	Difficulty      tcommon.JSONUint64 `json:"difficulty"`
	TotalDifficulty tcommon.JSONUint64 `json:"totalDifficulty"`
	Size            tcommon.JSONUint64 `json:"size"`
	GasLimit        tcommon.JSONUint64 `json:"gasLimit"`
	GasUsed         tcommon.JSONUint64 `json:"gasUsed"`
	ExtraData       []byte             `json:"extraData"`
	Uncles          []tcommon.Hash     `json:"uncles"`
	//TODO : add transactions
	// Txs  []Tx         `json:"transactions"`

	// ChainID   string             `json:"chain_id"`
	// Epoch     tcommon.JSONUint64 `json:"epoch"`
	// HCC                core.CommitCertificate   `json:"hcc"`
	// GuardianVotes      *core.AggregatedVotes    `json:"guardian_votes"`
	// EliteEdgeNodeVotes *core.AggregatedEENVotes `json:"elite_edge_node_votes"`
	// Children []tcommon.Hash   `json:"children"`
	// Status   core.BlockStatus `json:"status"`
}

// ------------------------------- eth_getBlockByHash -----------------------------------
func (e *EthRPCService) GetBlockByHash(ctx context.Context, hashStr string, tx bool) (result eth_GetBlockResult, err error) {
	logger.Infof("eth_getBlockByHash called")

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetBlock", trpc.GetBlockArgs{Hash: tcommon.HexToHash(hashStr)})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetBlockResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult, nil
	}
	result = eth_GetBlockResult{}
	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return result, err
	}
	theta_GetBlockResult, ok := resultIntf.(trpc.GetBlockResult)
	if !ok {
		return result, fmt.Errorf("failed to convert GetBlockResult")
	}
	result.Height = theta_GetBlockResult.Height
	result.Hash = theta_GetBlockResult.Hash
	result.Parent = theta_GetBlockResult.Parent
	result.Timestamp = theta_GetBlockResult.Timestamp
	result.Proposer = theta_GetBlockResult.Proposer
	result.TxHash = theta_GetBlockResult.TxHash
	result.StateHash = theta_GetBlockResult.StateHash
	//TODO : handle tx
	//TODO : handle other fields
	return result, nil
}

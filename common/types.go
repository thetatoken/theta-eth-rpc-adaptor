package common

import (
	tcommon "github.com/thetatoken/theta/common"
)

type Bytes8 [8]byte
type Bytes256 [256]byte

type ethGetTransactionResult struct {
	BlockHash        tcommon.Hash       `json:"blockHash"`
	BlockHeight      tcommon.JSONUint64 `json:"blockNumber"`
	From             tcommon.Address    `json:"from"`
	To               tcommon.Address    `json:"to"`
	Gas              tcommon.JSONUint64 `json:"gas"`
	GasPrice         tcommon.JSONUint64 `json:"gasPrice"`
	TxHash           tcommon.Hash       `json:"hash"`
	Nonce            Bytes8             `json:"nonce"`
	Input            []byte             `json:"input"`
	TransactionIndex tcommon.JSONUint64 `json:"transactionIndex"`
	Value            tcommon.JSONUint64 `json:"value"`
	V                tcommon.JSONUint64 `json:"v"` //ECDSA recovery id
	R                tcommon.Hash       `json:"r"` //ECDSA signature r
	S                tcommon.Hash       `json:"s"` //ECDSA signature s
}

type ethGetBlockResult struct {
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
}

type ethSyncingResult struct {
	StartingBlock tcommon.JSONUint64 `json:"startingBlock"`
	CurrentBlock  tcommon.JSONUint64 `json:"currentBlock"`
	HighestBlock  tcommon.JSONUint64 `json:"highestBlock"`
	PulledStates  tcommon.JSONUint64 `json:"pulledStates"` //pulledStates is the number it already downloaded
	KnownStates   tcommon.JSONUint64 `json:"knownStates"`  //knownStates is the number of trie nodes that the sync algo knows about
}

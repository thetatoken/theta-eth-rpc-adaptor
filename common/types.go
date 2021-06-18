package common

import (
	"github.com/thetatoken/theta/blockchain"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/common/hexutil"
	"github.com/thetatoken/theta/ledger/types"
)

type Bytes8 [8]byte

//TODO: change more tcommon.JSONUint64 to hexutil.Uint64
type EthGetTransactionResult struct {
	BlockHash        tcommon.Hash    `json:"blockHash"`
	BlockHeight      hexutil.Uint64  `json:"blockNumber"`
	From             tcommon.Address `json:"from"`
	To               tcommon.Address `json:"to"`
	Gas              hexutil.Uint64  `json:"gas"`
	GasPrice         hexutil.Uint64  `json:"gasPrice"`
	TxHash           tcommon.Hash    `json:"hash"`
	Nonce            hexutil.Uint64  `json:"nonce"`
	Input            []byte          `json:"input"`
	TransactionIndex hexutil.Uint64  `json:"transactionIndex"`
	Value            hexutil.Uint64  `json:"value"`
	V                hexutil.Uint64  `json:"v"` //ECDSA recovery id
	R                tcommon.Hash    `json:"r"` //ECDSA signature r
	S                tcommon.Hash    `json:"s"` //ECDSA signature s
}

type EthGetBlockResult struct {
	Height    hexutil.Uint64  `json:"number"`
	Hash      tcommon.Hash    `json:"hash"`
	Parent    tcommon.Hash    `json:"parentHash"`
	Timestamp hexutil.Uint64  `json:"timestamp"`
	Proposer  tcommon.Address `json:"miner"`
	TxHash    tcommon.Hash    `json:"transactionsRoot"`
	StateHash tcommon.Hash    `json:"stateRoot"`

	ReiceptHash     tcommon.Hash   `json:"receiptsRoot"`
	Nonce           Bytes8         `json:"nonce"`
	Sha3Uncles      tcommon.Hash   `json:"sha3Uncles"`
	LogsBloom       tcommon.Bytes  `json:"logsBloom"`
	Difficulty      hexutil.Uint64 `json:"difficulty"`
	TotalDifficulty hexutil.Uint64 `json:"totalDifficulty"`
	Size            hexutil.Uint64 `json:"size"`
	GasLimit        hexutil.Uint64 `json:"gasLimit"`
	GasUsed         hexutil.Uint64 `json:"gasUsed"`
	ExtraData       []byte         `json:"extraData"`
	Uncles          []tcommon.Hash `json:"uncles"`
	Transactions    []interface{}  `json:"transactions"`
}

type EthSyncingResult struct {
	StartingBlock hexutil.Uint64 `json:"startingBlock"`
	CurrentBlock  hexutil.Uint64 `json:"currentBlock"`
	HighestBlock  hexutil.Uint64 `json:"highestBlock"`
	PulledStates  hexutil.Uint64 `json:"pulledStates"` //pulledStates is the number it already downloaded
	KnownStates   hexutil.Uint64 `json:"knownStates"`  //knownStates is the number of trie nodes that the sync algo knows about
}

type EthGetReceiptResult struct {
	BlockHash         tcommon.Hash    `json:"blockHash"`
	BlockHeight       hexutil.Uint64  `json:"blockNumber"`
	TxHash            tcommon.Hash    `json:"transactionHash"`
	TransactionIndex  hexutil.Uint64  `json:"transactionIndex"`
	ContractAddress   tcommon.Address `json:"contractAddress"`
	From              tcommon.Address `json:"from"`
	To                tcommon.Address `json:"to"`
	GasUsed           hexutil.Uint64  `json:"gasUsed"`
	CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
	Logs              []EthLogObj     `json:"logs"`
	LogsBloom         tcommon.Bytes   `json:"logsBloom"`
	Status            hexutil.Uint64  `json:"status"`
}

type Tx struct {
	types.Tx `json:"raw"`
	Type     byte                       `json:"type"`
	Hash     tcommon.Hash               `json:"hash"`
	Receipt  *blockchain.TxReceiptEntry `json:"receipt"`
}

type EthLogObj struct {
	Address          tcommon.Address `json:"address"`
	BlockHash        tcommon.Hash    `json:"blockHash"`
	BlockHeight      hexutil.Uint64  `json:"blockNumber"`
	LogIndex         hexutil.Uint64  `json:"logIndex"`
	Removed          bool            `json:"removed"`
	Topics           []tcommon.Hash  `json:"topics"`
	TxHash           tcommon.Hash    `json:"transactionHash"`
	TransactionIndex hexutil.Uint64  `json:"transactionIndex"`
	Data             tcommon.Bytes   `json:"data"`
}

type EthSmartContractArgObj struct {
	From     tcommon.Address `json:"from"`
	To       tcommon.Address `json:"to"`
	Gas      string          `json:"gas"`
	GasPrice string          `json:"gasPrice"`
	Value    string          `json:"value"`
	Data     string          `json:"data"`
}

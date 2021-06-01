package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	"github.com/thetatoken/theta/cmd/thetacli/cmd/utils"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/ledger/types"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

type ethSmartContractArgObj struct {
	From     tcommon.Address `json:"from"`
	To       tcommon.Address `json:"to"`
	Gas      string          `json:"gas"`
	GasPrice string          `json:"gasPrice"`
	Value    string          `json:"value"`
	Data     string          `json:"data"`
}

// ------------------------------- eth_call -----------------------------------

func (e *EthRPCService) Call(ctx context.Context, argObj ethSmartContractArgObj, tag string) (result string, err error) {
	logger.Infof("eth_call called")

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())

	fromAddress := argObj.From.String()

	rpcRes, rpcErr := client.Call("theta.GetAccount", trpc.GetAccountArgs{Address: fromAddress})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetAccountResult{Account: &types.Account{}}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Account.Sequence, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)

	if err != nil {
		return "", err
	}
	sequence := resultIntf.(uint64) + 1

	from := types.TxInput{
		Address: tcommon.HexToAddress(fromAddress),
		Coins: types.Coins{
			ThetaWei: new(big.Int).SetUint64(0),
			TFuelWei: new(big.Int).SetUint64(common.Str2hex2unit(argObj.Value)),
		},
		Sequence: sequence,
	}

	to := types.TxOutput{
		Address: tcommon.HexToAddress(argObj.To.String()),
	}

	gasPrice, ok := types.ParseCoinAmount(argObj.GasPrice)
	if !ok {
		utils.Error("Failed to parse gas price")
	}

	sctx := &types.SmartContractTx{
		From:     from,
		To:       to,
		GasLimit: 500000,
		GasPrice: gasPrice,
		Data:     []byte(argObj.Data),
	}

	sctxBytes, err := types.TxToBytes(sctx)

	if err != nil {
		utils.Error("Failed to encode smart contract transaction: %v\n", sctx)
	}

	client = rpcc.NewRPCClient(common.GetThetaRPCEndpoint())

	rpcRes, rpcErr = client.Call("theta.CallSmartContract", trpc.CallSmartContractArgs{SctxBytes: hex.EncodeToString(sctxBytes)})

	parse = func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.CallSmartContractResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.VmReturn, nil
	}

	resultIntf, err = common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	result = resultIntf.(string)
	return result, nil
}

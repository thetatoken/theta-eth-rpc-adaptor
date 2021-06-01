package common

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/thetatoken/theta/cmd/thetacli/cmd/utils"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/ledger/types"
	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

func GetThetaRPCEndpoint() string {
	thetaRPCEndpoint := viper.GetString(CfgThetaRPCEndpoint)
	return thetaRPCEndpoint
}

func HandleThetaRPCResponse(rpcRes *rpcc.RPCResponse, rpcErr error, parse func(jsonBytes []byte) (interface{}, error)) (result interface{}, err error) {
	if rpcErr != nil {
		return nil, fmt.Errorf("failed to get theta RPC response: %v", rpcErr)
	}
	if rpcRes.Error != nil {
		return nil, fmt.Errorf("theta RPC returns an error: %v", rpcRes.Error)
	}

	var jsonBytes []byte
	jsonBytes, err = json.MarshalIndent(rpcRes.Result, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to parse theta RPC response: %v, %s", err, string(jsonBytes))
	}

	result, err = parse(jsonBytes)
	return
}

func GetHeightByTag(tag string) (height tcommon.JSONUint64) {
	switch tag {
	case "latest":
		height = tcommon.JSONUint64(0)
	case "earliest":
		height = tcommon.JSONUint64(1)
	case "pending":
		height = tcommon.JSONUint64(0)
	default:
		height = tcommon.JSONUint64(Str2hex2unit(tag))
	}
	return height
}

func Str2hex2unit(str string) uint64 {
	// remove 0x suffix if found in the input string
	if strings.HasPrefix(str, "0x") {
		str = strings.TrimPrefix(str, "0x")
	}

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(str, 16, 64)
	return uint64(result)
}

func GetSctxBytes(arg EthSmartContractArgObj) (sctxBytes []byte, err error) {
	sequence, seqErr := GetSeqByAddress(arg.From)
	if seqErr != nil {
		utils.Error("Failed to get sequence by address: %v\n", arg.From)
		return sctxBytes, seqErr
	}
	from := types.TxInput{
		Address: tcommon.HexToAddress(arg.From.String()),
		Coins: types.Coins{
			ThetaWei: new(big.Int).SetUint64(0),
			TFuelWei: new(big.Int).SetUint64(Str2hex2unit(arg.Value)),
		},
		Sequence: sequence,
	}

	to := types.TxOutput{
		Address: tcommon.HexToAddress(arg.To.String()),
	}

	gasPrice, ok := types.ParseCoinAmount(arg.GasPrice)
	if !ok {
		utils.Error("Failed to parse gas price")
	}

	sctx := &types.SmartContractTx{
		From:     from,
		To:       to,
		GasLimit: 500000,
		GasPrice: gasPrice,
		Data:     []byte(arg.Data),
	}

	sctxBytes, err = types.TxToBytes(sctx)

	if err != nil {
		utils.Error("Failed to encode smart contract transaction: %v\n", sctx)
		return sctxBytes, err
	}
	return sctxBytes, nil
}

func GetSeqByAddress(address tcommon.Address) (sequence uint64, err error) {
	client := rpcc.NewRPCClient(GetThetaRPCEndpoint())

	rpcRes, rpcErr := client.Call("theta.GetAccount", trpc.GetAccountArgs{Address: address.String()})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetAccountResult{Account: &types.Account{}}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Account.Sequence, nil
	}

	resultIntf, err := HandleThetaRPCResponse(rpcRes, rpcErr, parse)

	if err != nil {
		return sequence, err
	}
	sequence = resultIntf.(uint64) + 1

	return sequence, nil
}

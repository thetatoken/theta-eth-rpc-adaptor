package common

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/ledger/types"
	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

var logger *log.Entry = log.WithFields(log.Fields{"prefix": "common"})

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

func Int2hex2str(num int) string {
	return "0x" + strconv.FormatInt(int64(num), 16)
}

func HexToBytes(hexStr string) ([]byte, error) {
	trimmedHexStr := strings.TrimPrefix(hexStr, "0x")
	data, err := hex.DecodeString(trimmedHexStr)
	return data, err
}

func GetSctxBytes(arg EthSmartContractArgObj) (sctxBytes []byte, err error) {
	sequence, seqErr := GetSeqByAddress(arg.From)
	if seqErr != nil {
		logger.Errorf("Failed to get sequence by address: %v\n", arg.From)
		// return sctxBytes, seqErr
		sequence = 1
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

	gasPriceStr := "0wei"
	if arg.GasPrice != "" {
		gasPriceStr = arg.GasPrice + "wei"
	}

	gasPrice, ok := types.ParseCoinAmount(gasPriceStr)
	if !ok {
		err = errors.New("failed to parse gas price")
		logger.Errorf(fmt.Sprintf("%v", err))
		return sctxBytes, err
	}

	data, err := HexToBytes(arg.Data)
	if err != nil {
		logger.Errorf("Failed to decode data: %v, err: %v\n", arg.Data, err)
		return sctxBytes, err
	}

	gas := uint64(1000000)
	if arg.Gas != "" {
		gas = Str2hex2unit(arg.Gas)
	}
	fmt.Printf("gas: %v\n", gas)

	sctx := &types.SmartContractTx{
		From:     from,
		To:       to,
		GasLimit: gas,
		GasPrice: gasPrice,
		Data:     data,
	}

	sctxBytes, err = types.TxToBytes(sctx)

	if err != nil {
		logger.Errorf("Failed to encode smart contract transaction: %v\n", sctx)
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

func GetCurrentHeight() (height tcommon.JSONUint64, err error) {
	client := rpcc.NewRPCClient(GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetStatus", trpc.GetStatusArgs{})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.CurrentHeight, nil
	}

	resultIntf, err := HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return height, err
	}
	height = resultIntf.(tcommon.JSONUint64)
	return height, nil
}

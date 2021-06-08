package ethrpc

import (
	"context"
	"strconv"
	"strings"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	tcommon "github.com/thetatoken/theta/common"
	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getBlockByNumber -----------------------------------
func (e *EthRPCService) GetBlockByNumber(ctx context.Context, numberStr string, txDetails bool) (result common.EthGetBlockResult, err error) {
	logger.Infof("eth_getBlockByNumber called")
	// height := common.GetHeightByTag(numberStr)
	height := GetHeightByTag(numberStr)
	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetBlockByHeight", trpc.GetBlockByHeightArgs{
		Height: height})
	return GetBlockFromTRPCResult(rpcRes, rpcErr, txDetails)
}

//TODO : use common.GetHeightByTag
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

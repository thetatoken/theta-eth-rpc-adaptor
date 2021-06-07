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

// ------------------------------- eth_getBlockByHash -----------------------------------
func (e *EthRPCService) GetBlockByNumber(ctx context.Context, numberStr string, txDetails bool) (result common.EthGetBlockResult, err error) {
	logger.Infof("eth_getBlockByNumber called")
	var height tcommon.JSONUint64
	//TODO: handle strings
	if numberStr == "earliest" {

	} else if numberStr == "latest" {

	} else if numberStr == "pending" {

	} else {
		// 	hashNum := tcommon.HexToHash(numberStr)
		cleaned := strings.Replace(numberStr, "0x", "", -1)
		tmpNum, err := strconv.ParseUint(cleaned, 16, 64)
		if err != nil {
			return common.EthGetBlockResult{}, err
		}
		height = tcommon.JSONUint64(tmpNum)
	}
	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetBlockByHeight", trpc.GetBlockByHeightArgs{
		Height: height})
	return GetBlockFromTRPCResult(rpcRes, rpcErr, txDetails)
}

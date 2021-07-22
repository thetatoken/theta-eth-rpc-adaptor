package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"

	trpc "github.com/thetatoken/theta/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getStorageAt -----------------------------------

func (e *EthRPCService) GetStorageAt(ctx context.Context, address string, storagePosition string, tag string) (result string, err error) {
	logger.Infof("eth_getStorageAt called")

	height := common.GetHeightByTag(tag)

	client := rpcc.NewRPCClient(common.GetThetaRPCEndpoint())
	rpcRes, rpcErr := client.Call("theta.GetStorageAt", trpc.GetStorageAtArgs{
		Address:         address,
		StoragePosition: storagePosition,
		Height:          height})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStorageAtResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Value, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}

	result = resultIntf.(string)
	if result == "0000000000000000000000000000000000000000000000000000000000000000" {
		result = "0x0"
	}
	return result, nil
}

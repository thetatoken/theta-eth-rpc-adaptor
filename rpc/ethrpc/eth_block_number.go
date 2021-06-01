package ethrpc

import (
	"context"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	hexutil "github.com/thetatoken/theta/common/hexutil"
)

// ------------------------------- eth_blockNumber -----------------------------------

func (e *EthRPCService) BlockNumber(ctx context.Context) (result string, err error) {
	logger.Infof("eth_blockNumber called")

	blockNumber, err := common.GetCurrentHeight()

	if err != nil {
		return "", err
	}

	result = hexutil.EncodeUint64(uint64(blockNumber))
	return result, nil
}

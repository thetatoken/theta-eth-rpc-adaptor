package ethrpc

import (
	"context"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
)

// ------------------------------- eth_accounts -----------------------------------

func (e *EthRPCService) Accounts(ctx context.Context) (result []string, err error) {
	logger.Infof("eth_accounts called")
	return common.TestWalletArr, nil
}

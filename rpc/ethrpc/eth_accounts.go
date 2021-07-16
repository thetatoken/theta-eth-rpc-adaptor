package ethrpc

import (
	"context"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
)

// ------------------------------- eth_accounts -----------------------------------

func (e *EthRPCService) Accounts(ctx context.Context) (result []string, err error) {
	logger.Infof("eth_accounts called")
	for key, _ := range common.TestWallets {
		result = append(result, key)
	}
	return result, nil
}

package ethrpc

import (
	"context"
)

// ------------------------------- eth_accounts -----------------------------------

func (e *EthRPCService) Accounts(ctx context.Context, address string, tag string) (result string, err error) {
	logger.Infof("eth_accounts called")

	return result, nil
}

package web3rpc

import (
	"context"
)

// ------------------------------- web3_clientVersion -----------------------------------

func (e *Web3RPCService) ClientVersion(ctx context.Context) (result string, err error) {
	logger.Infof("web3_clientVersion called")

	return result, nil
}

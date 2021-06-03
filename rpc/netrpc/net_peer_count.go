package netrpc

import (
	"context"
)

// ------------------------------- net_peerCount -----------------------------------

func (e *NetRPCService) PeerCount(ctx context.Context) (result string, err error) {
	logger.Infof("net_peerCount called")

	return result, nil
}

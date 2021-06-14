package netrpc

import (
	"context"
)

// ------------------------------- net_listening -----------------------------------

func (e *NetRPCService) Listening(ctx context.Context) (result bool, err error) {
	logger.Infof("net_listening called")

	return result, nil
}

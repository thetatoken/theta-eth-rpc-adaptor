package netrpc

import (
	"context"
)

// ------------------------------- net_version -----------------------------------

func (e *NetRPCService) Version(ctx context.Context) (result string, err error) {
	logger.Infof("net_version called")

	return result, nil
}

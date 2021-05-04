package common

import (
	"github.com/spf13/viper"
)

const (
	// CfgConfigPath defines custom config path
	CfgConfigPath = "config.path"

	// CfgThetaRPCEndpoint configures the Theta RPC endpoint
	CfgThetaRPCEndpoint = "theta.rpcEndpoint"

	// CfgRPCEnabled sets whether to run RPC service.
	CfgRPCEnabled = "rpc.enabled"
	// CfgRPCAddress sets the binding address of RPC service.
	CfgRPCAddress = "rpc.address"
	// CfgRPCPort sets the port of RPC service.
	CfgRPCPort = "rpc.port"
	// CfgRPCMaxConnections limits concurrent connections accepted by RPC server.
	CfgRPCMaxConnections = "rpc.maxConnections"
	// CfgRPCTimeoutSecs set a timeout for RPC.
	CfgRPCTimeoutSecs = "rpc.timeoutSecs"

	// CfgLogLevels sets the log level.
	CfgLogLevels = "log.levels"
	// CfgLogPrintSelfID determines whether to print node's ID in log (Useful in simulation when
	// there are more than one node running).
	CfgLogPrintSelfID = "log.printSelfID"

	// CfgForceGCEnabled to enable force GC
	CfgForceGCEnabled = "gc.enabled"
)

func init() {
	viper.SetDefault(CfgThetaRPCEndpoint, "http://127.0.0.1:16888/rpc")

	viper.SetDefault(CfgRPCAddress, "0.0.0.0")
	viper.SetDefault(CfgRPCPort, "18888")
	viper.SetDefault(CfgRPCMaxConnections, 2048)
	viper.SetDefault(CfgRPCTimeoutSecs, 600)

	viper.SetDefault(CfgLogLevels, "*:debug")
	viper.SetDefault(CfgLogPrintSelfID, false)
}

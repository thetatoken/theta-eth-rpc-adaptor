package cmds

import (
	"context"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/thetatoken/theta-eth-rpc-adaptor/node"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Theta ETH RPC Adaptor",
	Run:   runStart,
}

func init() {
	RootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) {
	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(context.Background())

	n := node.NewNode()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	done := make(chan struct{})
	go func() {
		<-c
		signal.Stop(c)
		cancel()
		// Wait at most 5 seconds before forcefully shutting down.
		<-time.After(time.Duration(5) * time.Second)
		close(done)
	}()

	n.Start(ctx)

	go func() {
		n.Wait()
		close(done)
	}()

	<-done
	log.Infof("")
	log.Infof("Graceful exit.")
}

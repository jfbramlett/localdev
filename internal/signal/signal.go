package signal

import (
	"context"
	"fmt"
	splicelogger "github.com/splice/platform/infra/libs/golang/logger"
	"github.com/splice/platform/localdev/v2/localdev/internal/command"
	"os"
	"os/signal"
	"syscall"
)

// WaitForSignal hooks into the user session to capture commands like ctrl-c so we can cleanly shut down running
// operations
func WaitForSignal(ctx context.Context, commands ...command.Command) {
	logger, _ := splicelogger.FromContext(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		for _, c := range commands {
			if err := c.Stop(ctx); err != nil {
				logger.WithField("exe_command", c).WithError(err).Error("failed stopping command")
			}
		}
		fmt.Println(sig)
		done <- true
	}()

	<-done
}

package main

import (
	"context"
	splicelogger "github.com/splice/platform/infra/libs/golang/logger"
	"github.com/splice/platform/localdev/v2/localdev/internal/runagainststaging"
	"os"

	"github.com/spf13/cobra"
)

const (
	flagRebuild = "rebuild"
)

// runAgainstStagingCommand represents a command that, based on the current working directory, will attempt to run the given
// service against our staging environment. It relies on conventions to keep things simple - it will create an ssh tunnel
// to our staging DB (primary and replica) and will run docker-compose against a `docker-compose-staging.yml` that is expected
// to be in the current directory. It relies on a file `private.env` for feeding environment variables into the docker-compose
// and this file is expected to be in the root of $SPLICE_PLATFORM.
// To stop the environment press CTRL-C in the terminal.
var runAgainstStagingCommand = &cobra.Command{
	Use:   "run-against-staging",
	Short: "runs a local service against staging",
	Run:   runAgainstStaging,
}

func init() {
	flags := runAgainstStagingCommand.Flags()
	flags.Bool(flagRebuild, false, "flag indicating if we should rebuild before running")

	rootCmd.AddCommand(runAgainstStagingCommand)
}

func runAgainstStaging(cmd *cobra.Command, _ []string) {
	logger := splicelogger.New()
	logger = logger.WithField("command", "run-against-staging")
	ctx := splicelogger.ContextWithLogger(context.Background(), logger)

	rebuild, err := cmd.Flags().GetBool(flagRebuild)
	if err != nil {
		logger.WithError(err).Error("failed to resolve our rebuild flag")
		os.Exit(1)
	}
	runner := runagainststaging.NewRunner(rebuild)

	if err := runner.Run(ctx); err != nil {
		logger.WithError(err).Error("failed executing our runner")
		os.Exit(1)
	}
}

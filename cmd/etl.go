package main

import (
	"context"
	"github.com/splice/platform/localdev/v2/localdev/internal/etl"

	splicelogger "github.com/splice/platform/infra/libs/golang/logger"

	"github.com/spf13/cobra"
)

// The serveCmd will execute the generate command
var extractAndLoadCmd = &cobra.Command{
	Use:   "etl",
	Short: "starts the extract and load process",
	Run:   extractAndLoad,
}

func init() {
	// Here we create the command line flags for our app, and bind them to our package-local
	// config variable.
	flags := extractAndLoadCmd.Flags()
	flags.String("source_db_dsn", "<user>:<password>@tcp(localhost:57527)/splice_staging", "the source database")
	flags.String("target_db_dsn", "root:@tcp(localhost:3307)/splice_local", "the database to load")

	// Add the "serve" sub-command to the root command.
	rootCmd.AddCommand(extractAndLoadCmd)
}

func extractAndLoad(cmd *cobra.Command, args []string) {
	logger := splicelogger.New()
	ctx := splicelogger.ContextWithLogger(context.Background(), logger)

	sourceDSN, err := cmd.Flags().GetString("source_db_dsn")
	if err != nil {
		logger.WithError(err).Error("failed to resolve source database dsn")
		return
	}

	targetDSN, err := cmd.Flags().GetString("target_db_dsn")
	if err != nil {
		logger.WithError(err).Error("failed to resolve target database dsn")
		return
	}

	driver, err := etl.NewManualDriver(sourceDSN, targetDSN)
	if err != nil {
		logger.WithError(err).Error("failed to initialize our driver")
		return
	}

	err = driver.Run(ctx)
	if err != nil {
		logger.WithError(err).Error("failed executing our driver")
		return
	}
}

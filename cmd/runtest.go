package main

import (
	"context"
	splicelogger "github.com/splice/platform/infra/libs/golang/logger"
	"github.com/splice/platform/localdev/v2/localdev/internal/runtest"
	"os"

	"github.com/spf13/cobra"
)

const (
	flagTestType = "type"
)

// runTestCommand represents a command that, based on the current working directory, will attempt to run the given
// suite of tests as specified by the flag `--type` where `type` can be `unit`, `integration` or `all`.
// The type of tests are distinguished using a build tag (at the top of a test file using: `// +build unit` or `// +build integration`).
// The tests are run using gotestsum.
var runTestCommand = &cobra.Command{
	Use:   "run-test",
	Short: "runs a set of tests",
	Run:   runTests,
}

func init() {
	flags := runTestCommand.Flags()
	flags.String(flagTestType, "unit", "the class of test to run (can be one of 'unit', 'integration', or 'all'")

	rootCmd.AddCommand(runTestCommand)
}

func runTests(cmd *cobra.Command, _ []string) {
	logger := splicelogger.New()
	logger = logger.WithField("command", "run-test")

	testType, err := cmd.Flags().GetString(flagTestType)
	if err != nil {
		logger.WithError(err).Error("failed to resolve test type arg")
		os.Exit(1)
	}
	logger = logger.WithField("type", testType)
	ctx := splicelogger.ContextWithLogger(context.Background(), logger)

	runner := runtest.NewRunTest(testType)
	if err := runner.Run(ctx); err != nil {
		logger.WithError(err).Error("failed running tests")
		os.Exit(1)
	}
}

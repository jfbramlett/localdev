package main

import (
	"context"
	splicelogger "github.com/splice/platform/infra/libs/golang/logger"
	"github.com/splice/platform/localdev/v2/localdev/internal/awslogin"
	"os"

	"github.com/spf13/cobra"
)

// runUnitTestCommand represents a command that, based on the current working directory, will attempt to run the given
// suite of unit tests. Since these are unit tests there is no bootstrap as these tests are expected to run without
// any external processes configured. It will simply run gotestsum against the current working directory providing a
// tag of `unit`.
var awsLoginCommand = &cobra.Command{
	Use:   "aws-login",
	Short: "logs into aws",
	Run:   awsLogin,
}

func init() {
	rootCmd.AddCommand(awsLoginCommand)
}

func awsLogin(_ *cobra.Command, _ []string) {
	logger := splicelogger.New()
	logger = logger.WithField("command", "aws-login")
	ctx := splicelogger.ContextWithLogger(context.Background(), logger)

	login := awslogin.NewAWSLogin()
	if err := login.ECRLogin(ctx); err != nil {
		logger.WithError(err).Error("failed logging into AWS")
		os.Exit(1)
	}
}

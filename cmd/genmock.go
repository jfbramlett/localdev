package main

import (
	"github.com/splice/platform/localdev/v2/localdev/internal/mockgen"
	"os"

	splicelogger "github.com/splice/platform/infra/libs/golang/logger"

	"github.com/spf13/cobra"
)

const flagType = "type"

// genMockCommand is a command used to generate a sample mock for a proto
var genMockCommand = &cobra.Command{
	Use:   "genmock",
	Short: "generates a reference mock representation of a proto",
	Run:   genMock,
}

func init() {
	flags := genMockCommand.Flags()
	flags.String(flagType, "", "the proto type to generate for")

	rootCmd.AddCommand(genMockCommand)
}

func genMock(cmd *cobra.Command, _ []string) {
	logger := splicelogger.New()
	logger = logger.WithField("command", "mockgen")

	mockGen := mockgen.NewMockGenerator()

	typ, err := cmd.Flags().GetString(flagType)
	if err != nil {
		logger.WithError(err).Error("Flag required")
		os.Exit(1)
	}

	if _, err := mockGen.GenerateMock(typ); err != nil {
		logger.WithError(err).Error("failed logging into AWS")
		os.Exit(1)
	}
}

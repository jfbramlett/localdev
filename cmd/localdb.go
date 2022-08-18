package main

import (
	"github.com/spf13/cobra"
)

// The localdbCmd will execute the generate command
var localdbCmd = &cobra.Command{
	Use:   "localdb",
	Short: "starts a localdb and runs the migrations against it",
	Run:   localdb,
}

func init() {
	// Add the "localdb" sub-command to the root command.
	rootCmd.AddCommand(localdbCmd)
}

func localdb(_ *cobra.Command, _ []string) {
}

package main

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the root command for the app.
var rootCmd = &cobra.Command{
	Use:   "localdev",
	Short: "localdev",
}

// Execute rootCmd.
func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"context"
	splicerouter2 "github.com/splice/platform/localdev/v2/localdev/internal/splicerouter"
	"os"
	"path"

	splicelogger "github.com/splice/platform/infra/libs/golang/logger"

	"github.com/spf13/cobra"
)

const (
	flagConfig  = "config"
	flagWebRoot = "webRoot"
)

// routerCommand runs a proxy service that can be used to route requests internally or externally based on a JSON
// configuration. This services as a mini "gateway" allowing you to configure some components to run locally while others
// might live in a shared environment. In addition, it can serve static content for some routes allowing for mock data
// (for say an integration test). The mock data relies on convention where for a given path you define a directory containing
// the response data with the files named for the RPC method being invoked. So if we were proxying
// `twirp/splice.api.entitlement.v1.EntitlementService/GetEntitlements` to a mock response we would define a route for
// `twirp/splice.api.entitlement.v1.EntitlementService` (or you can include the method if you want to do this for a single
// method vs all for the service) and then in our mock data directory have a file named `GetEntitlements.json`.
// There are 2 separate configurations, one for internal service-to-service routing and another for external requests. This
// allows for internal calls between services which have slightly different requirements than those that come from external
// that would go through our traditional API gateway.
var routerCommand = &cobra.Command{
	Use:   "router",
	Short: "runs the splice router",
	Run:   router,
}

func init() {
	webRoot := path.Join(os.Getenv("SPLICE_PLATFORM"), "infra", "cmd", "localdev", "web")
	flags := routerCommand.Flags()
	flags.String(flagConfig, splicerouter2.RouterConfig, "the configuration file describing how splice router will route messages")
	flags.String(flagWebRoot, webRoot, "the root directory containing our web files")

	rootCmd.AddCommand(routerCommand)
}

func router(cmd *cobra.Command, _ []string) {
	logger := splicelogger.New()
	logger = logger.WithField("command", "router")
	ctx := splicelogger.ContextWithLogger(context.Background(), logger)

	cfgFile, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		logger.WithError(err).Error("error getting flag holding splice router config")
		os.Exit(1)
	}

	webRoot, err := cmd.Flags().GetString(flagWebRoot)
	if err != nil {
		logger.WithError(err).Error("error getting flag holding web root")
		os.Exit(1)
	}

	router := splicerouter2.NewSpliceRouter(cfgFile, webRoot)
	if err := router.Run(ctx); err != nil {
		logger.WithError(err).Error("failed running our splice router")
		os.Exit(1)
	}
}

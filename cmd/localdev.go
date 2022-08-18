package main

import (
	"context"
	splicerouter2 "github.com/splice/platform/localdev/v2/localdev/internal/splicerouter"
	"os"
	"path"

	splicelogger "github.com/splice/platform/infra/libs/golang/logger"

	"github.com/spf13/cobra"
	//"github.com/webview/webview"
)

var localdevCommand = &cobra.Command{
	Use:   "localdev",
	Short: "runs the splice router",
	Run:   localdev,
}

func init() {
	rootCmd.AddCommand(localdevCommand)
}

func localdev(cmd *cobra.Command, _ []string) {
	logger := splicelogger.New()
	logger = logger.WithField("command", "localdev")
	ctx := splicelogger.ContextWithLogger(context.Background(), logger)

	webRoot := path.Join(os.Getenv("SPLICE_PLATFORM"), "infra", "cmd", "localdev", "web")
	router := splicerouter2.NewSpliceRouter(splicerouter2.RouterConfig, webRoot)
	go func() {
		_ = router.Run(ctx)
	}()
	/* USING WEBVIEW HAS ISSUES WITH CI BUILDS AS WE'LL NEED TO ADD ADDITIONAL LIBS ONTO THE BUILD IMAGE
		SO COMMENTING OUT FOR NOW, BUT IF WE WANT A LOCAL UI THIS IS A GOOD OPTION
	   	debug := true
	   	w := webview.New(debug)

	   	copyPasteShortcut := `
	   window.addEventListener("keypress", (event) => {
	     if (event.metaKey && event.key === 'c') {
	       document.execCommand("copy")
	       event.preventDefault();
	     }
	     if (event.metaKey && event.key === 'v') {
	       document.execCommand("paste")
	       event.preventDefault();
	     }
	   })
	   `

	   	w.Eval(copyPasteShortcut)
	   	defer w.Destroy()
	   	w.SetTitle("Splice Local Development Toolkit")
	   	w.SetSize(800, 600, webview.HintNone)
	   	w.Navigate("http://localhost:8001")
	   	w.Run()
	*/
}

package cmd

import (
	"github.com/opub/scoreplus/web"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts application web server.",
	Long:  "Launches application web server on local port.",
	Run: func(cmd *cobra.Command, args []string) {
		web.Start()
	},
}

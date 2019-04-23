package cmd

import (
	"fmt"
	"os"

	"github.com/opub/scoreplus/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scoreplus",
	Short: "Score+ is a score tracking and sharing application.",
	Long:  "See https://github.com/opub/scoreplus for more details.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nBuild Date: %s\nGit Commit: %s\nGit Branch: %s\nGit State: %s\n", util.Version, util.BuildDate, util.GitCommit, util.GitBranch, util.GitState)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(buildCmd)
}

//Execute is main application entry point for all of the supported commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

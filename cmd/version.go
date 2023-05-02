package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of SnipForge",
	Long: `Display the version of SnipForge

  Find more information at: https://github.com/peetya/snipforge-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("SnipForge %s-%s (%s)\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

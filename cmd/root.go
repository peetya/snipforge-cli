package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var logLevel string

var rootCmd = &cobra.Command{
	Use:   "snipforge",
	Short: "SnipForge - AI Code Snippet Generator",
	Long: `SnipForge - AI Code Snippet Generator

  Find more information at: https://github.com/peetya/snipforge-cli`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", logrus.InfoLevel.String(), "The log level")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return setupLogger()
	}
}

func setupLogger() error {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	logrus.SetFormatter(&logrus.TextFormatter{})

	return nil
}

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var logLevel string

var rootCmd = &cobra.Command{
	Use:   "snipforge",
	Short: "A CLI to generate code snippets",
	Long: `SnipForge - A CLI tool for generating code snippets

Easily generate code snippets in various programming languages using natural language descriptions with SnipForge. 
Powered by ChatGPT, SnipForge helps you create custom code snippets quickly and efficiently.`,
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
		if err := initConfig(); err != nil {
			return err
		}
		if err := setupLogger(); err != nil {
			return err
		}
		return nil
	}
}

func initConfig() error {
	viper.SetConfigName("snipforge")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	viper.SetEnvPrefix("SNIPFORGE")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	return nil
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

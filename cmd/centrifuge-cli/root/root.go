/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package root

import (
	"fmt"

	"github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/centrifuge-cli/internal"
	"github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/centrifuge-cli/root/publish"
	"github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/centrifuge-cli/root/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	prettyLogFlagName    = "pretty-log"
	prettyLogEnvVariable = "PRETTY_LOG"

	logLevelFlagName    = "log-level"
	logLevelEnvVariable = "LOG_LEVEL"
)

var (
	prettyLog bool
	logLevel  string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(cmd *cobra.Command) {
	cobra.CheckErr(cmd.Execute())
}

// ExecuteE adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ExecuteE(cmd *cobra.Command) error {
	return cmd.Execute()
}

// InitRootCmd initializes the root command
func InitRootCmd() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "A centrifugo client CLI tool",
		Long:  ``,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kafkaCLI.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	var flagName string

	flagName = logLevelFlagName

	flagName = prettyLogFlagName
	rootCmd.PersistentFlags().BoolVar(&prettyLog, flagName, false, fmt.Sprintf("i.e. --%s=true", flagName))
	viper.BindPFlag(flagName, rootCmd.PersistentFlags().Lookup(flagName))

	flagName = "port"
	rootCmd.PersistentFlags().IntVar(&internal.Port, flagName, 8000, fmt.Sprintf("[required] i.e. --%s=8000", flagName))
	viper.BindPFlag(flagName, rootCmd.PersistentFlags().Lookup(flagName))

	flagName = "channel"
	rootCmd.PersistentFlags().StringVar(&internal.Channel, flagName, "chat:index", fmt.Sprintf("[required] i.e. --%s=chat:index", flagName))
	viper.BindPFlag(flagName, rootCmd.PersistentFlags().Lookup(flagName))

	version.Init(rootCmd)
	publish.Init(rootCmd)
	return rootCmd
}

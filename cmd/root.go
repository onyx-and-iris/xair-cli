/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/onyx-and-iris/xair-cli/internal/xair"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xair-cli",
	Short: "A command-line utility to interact with Behringer X Air mixers via OSC",
	Long: `xair-cli is a command-line tool that allows users to send OSC messages
to Behringer X Air mixers for remote control and configuration. It supports
various commands to manage mixer settings directly from the terminal.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		level, err := log.ParseLevel(viper.GetString("loglevel"))
		if err != nil {
			return err
		}
		log.SetLevel(level)

		client, err := xair.NewClient(viper.GetString("host"), viper.GetInt("port"))
		if err != nil {
			return err
		}
		cmd.SetContext(WithContext(cmd.Context(), client))

		client.StartListening()
		err, resp := client.RequestInfo()
		if err != nil {
			return err
		}

		log.Infof("Received mixer info: %+v", resp)

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client != nil {
			client.Stop()
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("host", "H", "mixer.local", "host address of the X Air mixer")
	rootCmd.PersistentFlags().IntP("port", "p", 10024, "Port number of the X Air mixer")
	rootCmd.PersistentFlags().
		StringP("loglevel", "l", "warn", "Log level (debug, info, warn, error, fatal, panic)")

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("XAIR_CLI")
	viper.AutomaticEnv()
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
}

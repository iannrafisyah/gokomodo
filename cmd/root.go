package cmd

import (
	"github.com/iannrafisyah/gokomodo/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "service",
	Short: "Service App",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	config.SetConfig()
}

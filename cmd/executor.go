package cmd

import (
	"os"
	"webhook/pkg"

	"github.com/spf13/cobra"
)

const (
	MainServiceName = "webhook"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: MainServiceName,
	}
	cmd.AddCommand(serveCmd)
	cmd.AddCommand(migrateCmd)

	return cmd
}

func Execute() {
	pkg.InitLogger()

	rootCmd := NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

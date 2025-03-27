package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the API server",
	RunE:  runAPICmd,
}

func runAPICmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

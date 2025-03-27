package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func NewAPICmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the API server",
		RunE:  runAPICmd,
	}
	return serveCmd
}

func runAPICmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

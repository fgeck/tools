package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var removeName string

func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove a tool bookmark",
		Long:    "Remove a CLI tool bookmark by name",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.DeleteTool(context.Background(), removeName); err != nil {
				return fmt.Errorf("failed to remove tool: %w", err)
			}

			fmt.Printf("Successfully removed tool: %s\n", removeName)
			return nil
		},
	}

	cmd.Flags().StringVarP(&removeName, "name", "n", "", "Tool name to remove (required)")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

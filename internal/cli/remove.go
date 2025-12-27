package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var removeCommand string

func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove an example bookmark",
		Long:    "Remove an example by its command (primary key)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.DeleteExample(context.Background(), removeCommand); err != nil {
				return fmt.Errorf("failed to remove example: %w", err)
			}

			fmt.Printf("Successfully removed example: %s\n", removeCommand)
			return nil
		},
	}

	cmd.Flags().StringVarP(&removeCommand, "command", "c", "", "Command to remove (required)")
	_ = cmd.MarkFlagRequired("command")

	return cmd
}

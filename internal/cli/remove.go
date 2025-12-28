package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	removeCommand  string
	removeToolName string
)

func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove example bookmark(s)",
		Long: `Remove examples by command or tool name.

Use -c to remove a specific example by its command (primary key).
Use -n to remove all examples for a tool name.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// Must specify either command or tool name, but not both
			if removeCommand == "" && removeToolName == "" {
				return fmt.Errorf("must specify either --command (-c) or --name (-n)")
			}
			if removeCommand != "" && removeToolName != "" {
				return fmt.Errorf("cannot specify both --command and --name, choose one")
			}

			// Remove by command (single example)
			if removeCommand != "" {
				if err := svc.DeleteBookmark(ctx, removeCommand); err != nil {
					return fmt.Errorf("failed to remove example: %w", err)
				}
				fmt.Printf("Successfully removed example: %s\n", removeCommand)
				return nil
			}

			// Remove by tool name (all examples for that tool)
			if removeToolName != "" {
				if err := svc.DeleteToolBookmarks(ctx, removeToolName); err != nil {
					return fmt.Errorf("failed to remove examples for tool '%s': %w", removeToolName, err)
				}
				fmt.Printf("Successfully removed all examples for tool: %s\n", removeToolName)
				return nil
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&removeCommand, "command", "c", "", "Remove specific example by command")
	cmd.Flags().StringVarP(&removeToolName, "name", "n", "", "Remove all examples for tool name")

	return cmd
}

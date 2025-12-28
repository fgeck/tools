package cli

import (
	"context"
	"fmt"

	"github.com/fgeck/tools/internal/dto"
	"github.com/spf13/cobra"
)

var (
	addToolName   string
	addDesc       string
	addExampleCmd string
)

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new example bookmark",
		Long: `Add a new example to the bookmark manager.

Each example requires:
- Tool name: For grouping (e.g., "lsof")
- Description: What it does (e.g., "list all ports at port 54321")
- Command: The actual command (e.g., "lsof -i :54321")`,
		RunE: func(cmd *cobra.Command, args []string) error {
			req := dto.CreateBookmarkRequest{
				Command:     addExampleCmd,
				ToolName:    addToolName,
				Description: addDesc,
			}

			resp, err := svc.CreateBookmark(context.Background(), req)
			if err != nil {
				return fmt.Errorf("failed to add example: %w", err)
			}

			fmt.Printf("Successfully added command: %s for tool: %s\n", resp.Command, resp.ToolName)
			return nil
		},
	}

	cmd.Flags().StringVarP(&addToolName, "name", "n", "", "Tool name for grouping (required)")
	cmd.Flags().StringVarP(&addDesc, "description", "d", "", "Description - what it does (required)")
	cmd.Flags().StringVarP(&addExampleCmd, "command", "c", "", "The actual command to execute (required)")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("description")
	_ = cmd.MarkFlagRequired("command")

	return cmd
}

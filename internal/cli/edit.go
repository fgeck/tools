package cli

import (
	"context"
	"fmt"

	"github.com/fgeck/tools/internal/dto"
	"github.com/spf13/cobra"
)

var (
	editCommand     string
	editNewToolName string
	editNewDesc     string
	editNewCommand  string
)

func newEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit",
		Aliases: []string{"e", "update"},
		Short:   "Edit an existing example bookmark",
		Long: `Edit an existing example by specifying its current command.
You can update the tool name, description, and/or command.
Only the fields you provide will be updated.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// At least one field must be provided for update
			if editNewToolName == "" && editNewDesc == "" && editNewCommand == "" {
				return fmt.Errorf("at least one field must be provided for update (--new-tool, --new-description, or --new-command)")
			}

			req := dto.UpdateExampleRequest{
				Command:        editCommand,
				NewToolName:    editNewToolName,
				NewDescription: editNewDesc,
				NewCommand:     editNewCommand,
			}

			resp, err := svc.UpdateExample(context.Background(), req)
			if err != nil {
				return fmt.Errorf("failed to edit example: %w", err)
			}

			fmt.Printf("Successfully updated example: %s\n", resp.Command)
			return nil
		},
	}

	cmd.Flags().StringVarP(&editCommand, "command", "c", "", "Current command to edit (required)")
	cmd.Flags().StringVarP(&editNewToolName, "new-tool", "t", "", "New tool name")
	cmd.Flags().StringVarP(&editNewDesc, "new-description", "d", "", "New description")
	cmd.Flags().StringVarP(&editNewCommand, "new-command", "n", "", "New command")

	_ = cmd.MarkFlagRequired("command")

	return cmd
}

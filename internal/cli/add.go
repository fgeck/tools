package cli

import (
	"context"
	"fmt"

	"github.com/fgeck/tools/internal/dto"
	"github.com/spf13/cobra"
)

var (
	addName        string
	addCommand     string
	addDescription string
	addExamples    []string
)

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new tool bookmark",
		Long:    "Add a new CLI tool bookmark with name, command, description, and examples",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := dto.CreateToolRequest{
				Name:        addName,
				Command:     addCommand,
				Description: addDescription,
				Examples:    addExamples,
			}

			resp, err := svc.CreateTool(context.Background(), req)
			if err != nil {
				return fmt.Errorf("failed to add tool: %w", err)
			}

			fmt.Printf("Successfully added tool: %s\n", resp.Name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&addName, "name", "n", "", "Tool name (required)")
	cmd.Flags().StringVarP(&addCommand, "command", "c", "", "Command/path to executable (required)")
	cmd.Flags().StringVarP(&addDescription, "description", "d", "", "Tool description")
	cmd.Flags().StringArrayVarP(&addExamples, "example", "e", []string{}, "Usage example (can be specified multiple times)")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("command")

	return cmd
}

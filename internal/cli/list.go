package cli

import (
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List all tool bookmarks",
		Long:    "Display all CLI tool bookmarks in a formatted table",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listExamples()
		},
	}

	return cmd
}

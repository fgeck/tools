package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fgeck/tools/internal/service"
	"github.com/spf13/cobra"
)

var (
	svc     service.ToolService
	rootCmd *cobra.Command
)

// Initialize sets up the CLI with the provided service
func Initialize(toolService service.ToolService) {
	svc = toolService

	rootCmd = &cobra.Command{
		Use:   "tools",
		Short: "A bookmark manager for your terminal",
		Long: `The single CLI tool to view, add or remove CLI tools.
Consider it as a bookmark manager for your terminal.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default behavior: list all tools
			return listTools()
		},
	}

	// Add subcommands
	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newRemoveCmd())
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// listTools is a shared function for displaying tools in table format
func listTools() error {
	resp, err := svc.ListTools(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list tools: %w", err)
	}

	if resp.Count == 0 {
		fmt.Println("No tools found. Use 'tools add' to add your first tool.")
		return nil
	}

	// Create tabwriter for aligned output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print header
	fmt.Fprintln(w, "NAME\tCOMMAND\tDESCRIPTION\tEXAMPLES")
	fmt.Fprintln(w, "----\t-------\t-----------\t--------")

	// Print rows
	for _, tool := range resp.Tools {
		examples := ""
		if len(tool.Examples) > 0 {
			examples = tool.Examples[0]
			if len(tool.Examples) > 1 {
				examples += fmt.Sprintf(" (+%d more)", len(tool.Examples)-1)
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			tool.Name,
			tool.Command,
			tool.Description,
			examples,
		)
	}

	w.Flush()
	fmt.Printf("\nTotal: %d tools\n", resp.Count)

	return nil
}

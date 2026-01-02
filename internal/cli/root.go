package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fgeck/tools/internal/service"
	"github.com/fgeck/tools/internal/tui"
	"github.com/fgeck/tools/internal/utils"
	"github.com/spf13/cobra"
)

var (
	svc     service.BookmarkService
	rootCmd *cobra.Command
	useCLI  bool
)

// Initialize sets up the CLI with the provided service
func Initialize(exampleService service.BookmarkService) {
	svc = exampleService

	rootCmd = &cobra.Command{
		Use:   "tools",
		Short: "A bookmark manager for your terminal",
		Long: `The single CLI tool to view, add or remove CLI tools.
Consider it as a bookmark manager for your terminal.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default behavior: launch TUI unless --cli flag is set
			if useCLI {
				return listExamples()
			}
			return tui.Run(svc)
		},
	}

	// Add global flag
	rootCmd.PersistentFlags().BoolVar(&useCLI, "cli", false, "Use classic CLI mode instead of TUI")

	// Add subcommands
	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newEditCmd())
	rootCmd.AddCommand(newRemoveCmd())
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// listExamples is a shared function for displaying examples in table format
func listExamples() error {
	resp, err := svc.ListBookmarks(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list examples: %w", err)
	}

	if resp.Count == 0 {
		fmt.Println("No examples found. Use 'tools add' to add your first example.")
		return nil
	}

	// Create tabwriter for aligned output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print header
	_, _ = fmt.Fprintln(w, "TOOL\tDESCRIPTION\tCOMMAND")
	_, _ = fmt.Fprintln(w, "----\t-----------\t-------")

	// Define column widths for wrapping
	const (
		descriptionWidth = 40
		commandWidth     = 50
	)

	// Print rows with wrapping support
	for _, example := range resp.Examples {
		rows := utils.SplitWrappedRows(
			example.ToolName,
			example.Description,
			example.Command,
			descriptionWidth,
			commandWidth,
		)

		for _, row := range rows {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", row[0], row[1], row[2])
		}
	}

	_ = w.Flush()
	fmt.Printf("\nTotal: %d examples\n", resp.Count)

	return nil
}

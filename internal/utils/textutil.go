package utils

import (
	"strings"

	"github.com/muesli/reflow/wordwrap"
)

// WrapText wraps text to the specified width using word boundaries.
// If width is <= 0, returns the original text without wrapping.
func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	return wordwrap.String(text, width)
}

// WrapToLines wraps text and returns it as a slice of lines.
// If width is <= 0, returns the original text as a single-element slice.
func WrapToLines(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}
	wrapped := wordwrap.String(text, width)
	return strings.Split(wrapped, "\n")
}

// SplitWrappedRows takes column data and widths, returns multi-row representation.
// The first column (tool) is only shown on the first row; continuation rows have empty tool column.
// This maintains alignment in tabular output formats like tabwriter.
func SplitWrappedRows(tool, description, command string, descWidth, cmdWidth int) [][]string {
	descLines := WrapToLines(description, descWidth)
	cmdLines := WrapToLines(command, cmdWidth)

	// Determine how many rows we need
	maxLines := len(descLines)
	if len(cmdLines) > maxLines {
		maxLines = len(cmdLines)
	}

	rows := make([][]string, maxLines)
	for i := 0; i < maxLines; i++ {
		row := make([]string, 3)

		// Tool name only on first line
		if i == 0 {
			row[0] = tool
		} else {
			row[0] = "" // Empty for continuation lines
		}

		// Description
		if i < len(descLines) {
			row[1] = descLines[i]
		} else {
			row[1] = ""
		}

		// Command
		if i < len(cmdLines) {
			row[2] = cmdLines[i]
		} else {
			row[2] = ""
		}

		rows[i] = row
	}

	return rows
}

package utils

import (
	"strings"
	"testing"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		width          int
		wantMaxLineLen int // Maximum line length in result
	}{
		{
			name:           "short text no wrap",
			input:          "short text",
			width:          20,
			wantMaxLineLen: 10,
		},
		{
			name:           "long text wraps",
			input:          "this is a very long text that should wrap at word boundaries",
			width:          20,
			wantMaxLineLen: 20,
		},
		{
			name:           "zero width returns original",
			input:          "test text",
			width:          0,
			wantMaxLineLen: 9,
		},
		{
			name:           "negative width returns original",
			input:          "test text",
			width:          -5,
			wantMaxLineLen: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapText(tt.input, tt.width)
			lines := strings.Split(result, "\n")
			for _, line := range lines {
				if len(line) > tt.wantMaxLineLen {
					t.Errorf("Line exceeds max length: got %d, want <= %d, line: %q", len(line), tt.wantMaxLineLen, line)
				}
			}
		})
	}
}

func TestWrapToLines(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		width         int
		wantMinLines  int
		wantMaxLinLen int
	}{
		{
			name:          "short text single line",
			input:         "short",
			width:         20,
			wantMinLines:  1,
			wantMaxLinLen: 20,
		},
		{
			name:          "long text multiple lines",
			input:         "this is a very long text that should wrap at word boundaries",
			width:         20,
			wantMinLines:  2,
			wantMaxLinLen: 20,
		},
		{
			name:          "zero width single line",
			input:         "test",
			width:         0,
			wantMinLines:  1,
			wantMaxLinLen: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := WrapToLines(tt.input, tt.width)

			if len(lines) < tt.wantMinLines {
				t.Errorf("Expected at least %d lines, got %d", tt.wantMinLines, len(lines))
			}

			for _, line := range lines {
				if len(line) > tt.wantMaxLinLen {
					t.Errorf("Line too long: %d chars (max %d), line: %q", len(line), tt.wantMaxLinLen, line)
				}
			}
		})
	}
}

func TestSplitWrappedRows(t *testing.T) {
	tests := []struct {
		name        string
		tool        string
		description string
		command     string
		descWidth   int
		cmdWidth    int
		wantRows    int
	}{
		{
			name:        "single line no wrapping",
			tool:        "kubectl",
			description: "short desc",
			command:     "kubectl get pods",
			descWidth:   20,
			cmdWidth:    20,
			wantRows:    1,
		},
		{
			name:        "description wraps",
			tool:        "kubectl",
			description: "this is a very long description that will wrap to multiple lines",
			command:     "kubectl get pods",
			descWidth:   15,
			cmdWidth:    20,
			wantRows:    6, // Description wraps to multiple lines
		},
		{
			name:        "command wraps",
			tool:        "lsof",
			description: "short",
			command:     "lsof -i :8080 | grep LISTEN | awk '{print $2}'",
			descWidth:   20,
			cmdWidth:    15,
			wantRows:    4, // Command wraps to multiple lines
		},
		{
			name:        "both wrap",
			tool:        "docker",
			description: "list all running containers with verbose output",
			command:     "docker ps --all --format '{{.ID}} {{.Names}} {{.Status}}'",
			descWidth:   15,
			cmdWidth:    20,
			wantRows:    4, // Both wrap
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := SplitWrappedRows(tt.tool, tt.description, tt.command, tt.descWidth, tt.cmdWidth)

			if len(rows) != tt.wantRows {
				t.Errorf("Expected %d rows, got %d", tt.wantRows, len(rows))
			}

			// First row should have tool name
			if len(rows) > 0 && rows[0][0] != tt.tool {
				t.Errorf("First row should have tool name %s, got %s", tt.tool, rows[0][0])
			}

			// Subsequent rows should have empty tool column
			for i := 1; i < len(rows); i++ {
				if rows[i][0] != "" {
					t.Errorf("Row %d should have empty tool column, got %s", i, rows[i][0])
				}
			}

			// Verify each row has 3 columns
			for i, row := range rows {
				if len(row) != 3 {
					t.Errorf("Row %d should have 3 columns, got %d", i, len(row))
				}
			}
		})
	}
}

func TestSplitWrappedRowsEdgeCases(t *testing.T) {
	t.Run("empty strings", func(t *testing.T) {
		rows := SplitWrappedRows("", "", "", 20, 20)
		if len(rows) != 1 {
			t.Errorf("Expected 1 row for empty strings, got %d", len(rows))
		}
		if rows[0][0] != "" || rows[0][1] != "" || rows[0][2] != "" {
			t.Errorf("Expected all empty strings in row")
		}
	})

	t.Run("unicode characters", func(t *testing.T) {
		rows := SplitWrappedRows("echo", "print unicode: 你好 世界", "echo 'Hello 世界'", 15, 15)
		if len(rows) < 1 {
			t.Errorf("Expected at least 1 row with unicode, got %d", len(rows))
		}
	})

	t.Run("very long single word", func(t *testing.T) {
		longWord := strings.Repeat("a", 100)
		rows := SplitWrappedRows("test", longWord, "cmd", 20, 20)
		// Should handle long words by breaking them
		if len(rows) < 1 {
			t.Errorf("Expected at least 1 row with long word, got %d", len(rows))
		}
	})

	t.Run("zero widths", func(t *testing.T) {
		rows := SplitWrappedRows("tool", "description", "command", 0, 0)
		if len(rows) != 1 {
			t.Errorf("Expected 1 row with zero widths, got %d", len(rows))
		}
	})
}

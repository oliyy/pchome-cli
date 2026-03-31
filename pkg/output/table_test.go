package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestTruncateEnd_ASCII(t *testing.T) {
	if got := TruncateEnd("abcdef", 6); got != "abcdef" {
		t.Fatalf("expected no truncation, got %q", got)
	}
	if got := TruncateEnd("abcdef", 5); got != "ab..." {
		t.Fatalf("unexpected truncation: %q", got)
	}
	if got := TruncateEnd("abcdef", 3); got != "..." {
		t.Fatalf("unexpected truncation: %q", got)
	}
}

func TestPrintTable_RedditMarkdownStyle(t *testing.T) {
	var buf bytes.Buffer
	PrintTable(&buf, [][]string{
		{"Col1", "Numeric Column"},
		{"Value 1", "10.0"},
		{"Separate", "-2,027.1"},
	})

	got := strings.TrimRight(buf.String(), "\n")
	lines := strings.Split(got, "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d:\n%s", len(lines), got)
	}

	// Header + separator + 2 data rows; ignore trailing whitespace differences.
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " ")
	}

	if !strings.Contains(lines[0], "|") || !strings.Contains(lines[0], "Col1") || !strings.Contains(lines[0], "Numeric Column") {
		t.Fatalf("unexpected header row: %q", lines[0])
	}
	if !strings.Contains(lines[1], "|") || !strings.Contains(lines[1], "-") {
		t.Fatalf("unexpected separator row: %q", lines[1])
	}

	// Numeric column should be right-aligned within its padded width.
	if !strings.HasSuffix(lines[2], "10.0") {
		t.Fatalf("expected first data row to end with numeric value, got %q", lines[2])
	}
	if !strings.HasSuffix(lines[3], "-2,027.1") {
		t.Fatalf("expected second data row to end with numeric value, got %q", lines[3])
	}
}

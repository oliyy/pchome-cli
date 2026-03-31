package output

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

var numberColRe = regexp.MustCompile(`^(\s*-?(\d|,| |[.])*\s*)$`)

// PrintTable prints a Reddit Markdown style ASCII table (as in ozh/ascii-tables).
// It keeps the CLI dependency-light while producing copy/paste friendly output.
func PrintTable(w io.Writer, rows [][]string) {
	if len(rows) == 0 {
		return
	}

	colCount := 0
	for _, r := range rows {
		if len(r) > colCount {
			colCount = len(r)
		}
	}
	if colCount == 0 {
		return
	}

	widths := make([]int, colCount)
	for _, r := range rows {
		for i := 0; i < colCount; i++ {
			var cell string
			if i < len(r) {
				cell = r[i]
			}
			cw := runewidth.StringWidth(cell)
			if cw > widths[i] {
				widths[i] = cw
			}
		}
	}

	isNumberCol := make([]bool, colCount)
	for i := range isNumberCol {
		isNumberCol[i] = true
	}
	for rIdx := 1; rIdx < len(rows); rIdx++ { // skip header row
		r := rows[rIdx]
		for i := 0; i < colCount; i++ {
			var cell string
			if i < len(r) {
				cell = r[i]
			}
			if cell == "" {
				continue
			}
			if !numberColRe.MatchString(cell) {
				isNumberCol[i] = false
			}
		}
	}

	// Header row.
	fmt.Fprintln(w, strings.TrimRight(buildRedditRow(rows[0], widths, isNumberCol, true), " "))
	if len(rows) == 1 {
		return
	}

	// Header separator row.
	fmt.Fprintln(w, strings.TrimRight(buildRedditSeparator(widths), " "))

	// Data rows.
	for _, r := range rows[1:] {
		fmt.Fprintln(w, strings.TrimRight(buildRedditRow(r, widths, isNumberCol, false), " "))
	}
}

func buildRedditSeparator(widths []int) string {
	var b strings.Builder
	b.Grow(64)

	for i, w := range widths {
		if i == 0 {
			b.WriteString(" ")
		} else {
			b.WriteString("|")
		}
		if w < 0 {
			w = 0
		}
		b.WriteString(strings.Repeat("-", w+2))
	}

	// cMR is a single trailing space in ozh/ascii-tables' reddit style.
	b.WriteString(" ")
	return b.String()
}

func buildRedditRow(row []string, widths []int, isNumberCol []bool, isHeader bool) string {
	var b strings.Builder
	b.Grow(64)

	for i := 0; i < len(widths); i++ {
		var cell string
		if i < len(row) {
			cell = row[i]
		}

		align := "l"
		if isHeader {
			align = "c"
		} else if i < len(isNumberCol) && isNumberCol[i] {
			align = "r"
		}

		cell = padDisplayWidth(cell, widths[i], align)

		if i == 0 {
			// With no left side, ascii-tables emits two spaces instead of "| ".
			b.WriteString("  ")
		} else {
			b.WriteString("| ")
		}
		b.WriteString(cell)
		b.WriteString(" ")
	}
	return b.String()
}

func padDisplayWidth(s string, width int, align string) string {
	additional := width - runewidth.StringWidth(s)
	if additional <= 0 {
		return s
	}
	switch align {
	case "r":
		return strings.Repeat(" ", additional) + s
	case "c":
		left := additional / 2
		right := additional - left
		return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
	default: // "l"
		return s + strings.Repeat(" ", additional)
	}
}

// TruncateEnd truncates a string to at most maxWidth (display width) and appends "..." when needed.
func TruncateEnd(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if runewidth.StringWidth(s) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return strings.Repeat(".", maxWidth)
	}

	limit := maxWidth - 3
	var b strings.Builder
	b.Grow(len(s))

	w := 0
	for _, r := range s {
		rw := runewidth.RuneWidth(r)
		if w+rw > limit {
			break
		}
		b.WriteRune(r)
		w += rw
	}
	b.WriteString("...")
	return b.String()
}

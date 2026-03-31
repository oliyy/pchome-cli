package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func withTimeout(cmd *cobra.Command, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return context.WithCancel(cmd.Context())
	}
	return context.WithTimeout(cmd.Context(), timeout)
}

func printJSON(cmd *cobra.Command, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}

func printNDJSON(cmd *cobra.Command, values any) error {
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetEscapeHTML(false)

	switch typed := values.(type) {
	case []any:
		for _, value := range typed {
			if err := enc.Encode(value); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("ndjson output is only supported for list results")
	}

	return nil
}

func formatMoney(n int) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return sign + s
	}
	var b strings.Builder
	b.Grow(len(s) + len(s)/3)
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	return sign + b.String()
}

func formatRating(v *float64) string {
	if v == nil {
		return ""
	}
	if math.Abs(*v-math.Round(*v)) < 1e-9 {
		return strconv.FormatInt(int64(math.Round(*v)), 10)
	}
	return fmt.Sprintf("%.1f", *v)
}

func formatRatingWithReviews(score *float64, reviews *int) string {
	scoreStr := formatRating(score)
	reviewStr := formatInt(reviews)

	switch {
	case scoreStr != "" && reviewStr != "":
		return fmt.Sprintf("%s (%s)", scoreStr, reviewStr)
	case scoreStr != "":
		return scoreStr
	case reviewStr != "":
		return fmt.Sprintf("(%s)", reviewStr)
	default:
		return ""
	}
}

func formatInt(v *int) string {
	if v == nil {
		return ""
	}
	return strconv.Itoa(*v)
}

func formatBoolMarker(v *bool) string {
	if v == nil {
		return "?"
	}
	if *v {
		return "Y"
	}
	return "N"
}

func parseCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		out = append(out, part)
	}
	return out
}

func effectiveNameWidth(base int, compact, wide bool) int {
	switch {
	case compact:
		return 40
	case wide:
		if base < 80 {
			return 80
		}
		return base
	default:
		return base
	}
}

func joinArgs(args []string) string {
	return strings.TrimSpace(strings.Join(args, " "))
}

func columnsFlagDefault(columns []string) string {
	if len(columns) == 0 {
		return ""
	}
	return strings.Join(columns, ",")
}

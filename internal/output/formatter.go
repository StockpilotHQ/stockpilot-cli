package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// Paginated unwraps {"count":N,"results":[...]} or falls back to a raw array.
func Paginated(data []byte) ([]map[string]any, error) {
	var wrapper struct {
		Results []map[string]any `json:"results"`
	}
	if err := json.Unmarshal(data, &wrapper); err == nil && wrapper.Results != nil {
		return wrapper.Results, nil
	}
	var list []map[string]any
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// PaginatedCount returns both the results and the total count from the wrapper.
func PaginatedCount(data []byte) ([]map[string]any, int, error) {
	var wrapper struct {
		Count   int              `json:"count"`
		Results []map[string]any `json:"results"`
	}
	if err := json.Unmarshal(data, &wrapper); err == nil && wrapper.Results != nil {
		return wrapper.Results, wrapper.Count, nil
	}
	var list []map[string]any
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, 0, err
	}
	return list, len(list), nil
}

// JSON prints data as formatted JSON to stdout.
func JSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// Table renders a table to stdout.
func Table(headers []string, rows [][]string) {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader(headers)
	t.SetBorder(false)
	t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetColumnSeparator("  ")
	t.SetHeaderLine(false)
	t.AppendBulk(rows)
	t.Render()
}

// Error prints an error message and exits 1.
func Fatal(msg string) {
	fmt.Fprintln(os.Stderr, "Error: "+msg)
	os.Exit(1)
}

// Drift result formatter
package internal

import (
	"encoding/json"
	"fmt"
)

func PrintDriftReport(drifts []Drift, format string) {
	switch format {
	case "json":
		printJSON(drifts)
	case "markdown":
		printMarkdown(drifts)
	default:
		printTable(drifts)
	}
}

func printJSON(drifts []Drift) {
	bytes, err := json.MarshalIndent(drifts, "", "  ")
	if err != nil {
		fmt.Println("Failed to encode JSON output")
		return
	}
	fmt.Println(string(bytes))
}

func printMarkdown(drifts []Drift) {
	fmt.Println("| Resource Type | Resource ID | Issue | Severity |")
	fmt.Println("|---------------|-------------|--------|----------|")
	for _, d := range drifts {
		fmt.Printf("| %s | %s | %s | %s |\n", d.ResourceType, d.ID, d.Issue, d.Severity)
	}
}

func printTable(drifts []Drift) {
	for _, d := range drifts {
		fmt.Printf("- [%s] %s (%s): %s\n", d.Severity, d.ID, d.ResourceType, d.Issue)
	}
}

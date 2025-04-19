// Main CLI entry
package main

import (
	"flag"
	"fmt"
	"log"

	"iac-drift-detector/internal"
)

func main() {
	// CLI flags
	stateFile := flag.String("state-file", "", "Path to Terraform .tfstate file")
	region := flag.String("region", "us-east-1", "AWS region")
	output := flag.String("output", "table", "Output format: table | json | markdown")
	flag.Parse()

	if *stateFile == "" {
		log.Fatal("--state-file is required")
	}

	// Parse Terraform state
	state, err := internal.ParseTFState(*stateFile)
	if err != nil {
		log.Fatalf("Failed to parse Terraform state: %v", err)
	}

	// Fetch AWS live resources
	awsResources, err := internal.FetchAWSResources(*region)
	if err != nil {
		log.Fatalf("Failed to fetch AWS resources: %v", err)
	}

	// Compare state and AWS
	drift := internal.CompareStateWithAWS(state, awsResources)

	// Output results
	if len(drift) == 0 {
		fmt.Println("✅ No drift detected.")
		return
	}

	fmt.Println("⚠️ Drift detected:")
	for _, d := range drift {
		fmt.Printf("- [%s] %s (%s): %s\n", d.Severity, d.ID, d.ResourceType, d.Issue)
	}
}

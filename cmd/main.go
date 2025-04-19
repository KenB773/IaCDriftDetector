// Main CLI entry
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/KenB773/IaCDriftDetector/internal"
)

func main() {
	// Custom help output
	flag.Usage = func() {
		fmt.Println("Usage: drift-detector [OPTIONS]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
	}

	// CLI flags
	configPath := flag.String("config", "", "Path to config.yaml file")
	region := flag.String("region", "", "AWS region")
	stateFile := flag.String("state-file", "", "Path to Terraform .tfstate file")
	output := flag.String("output", "", "Output format: table | json | markdown")
	outputFile := flag.String("output-file", "", "Write output to file")
	include := flag.String("include", "", "Comma-separated list of resource types to include")
	dryRun := flag.Bool("dry-run", false, "Run without fetching from AWS")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Handle --version
	if *showVersion {
		fmt.Println("IaC Drift Detector", internal.Version)
		return
	}

	// Load config from file if specified
	var cfg *internal.Config
	if *configPath != "" {
		loadedCfg, err := internal.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		cfg = loadedCfg
	} else {
		cfg = &internal.Config{}
	}

	// Apply CLI flags
	cliFlags := &internal.Config{
		Region:     *region,
		StateFile:  *stateFile,
		Output:     *output,
		OutputFile: *outputFile,
		DryRun:     *dryRun,
	}
	if *include != "" {
		cliFlags.Include = strings.Split(*include, ",")
	}
	cfg = internal.MergeConfigWithFlags(cfg, cliFlags)

	if cfg.StateFile == "" {
		log.Fatal("--state-file is required (or define it in the config file)")
	}

	// Parse Terraform state
	state, err := internal.ParseTFState(cfg.StateFile)
	if err != nil {
		log.Fatalf("Failed to parse Terraform state: %v", err)
	}

	var awsResources []internal.FetchedResource
	if cfg.DryRun {
		color.Yellow("[DRY RUN] Skipping AWS resource fetch.")
	} else {
		awsResources, err = internal.FetchAWSResources(cfg.Region)
		if err != nil {
			log.Fatalf("Failed to fetch AWS resources: %v", err)
		}
	}

	// Compare state and AWS
	drift := internal.CompareStateWithAWS(state, awsResources, cfg.Include)

	// Output results
	if len(drift) == 0 {
		color.Green("✅ No drift detected.")
	} else {
		color.Red("⚠️  Drift detected:")
		internal.PrintDriftReport(drift, cfg.Output)
	}

	// Save output to file if specified
	if cfg.OutputFile != "" {
		err := internal.SaveToFile(cfg.OutputFile, drift)
		if err != nil {
			log.Fatalf("Failed to write drift report to file: %v", err)
		}
		color.Cyan("📄 Drift report saved to %s\n", cfg.OutputFile)
	}
}

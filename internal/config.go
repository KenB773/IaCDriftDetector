// Config file parser
package internal

import (
	"gopkg.in/yaml.v3"
	"os"
)

const Version = "v0.1.0"

type Config struct {
	Region      string   `yaml:"region"`
	Include     []string `yaml:"include"`
	Output      string   `yaml:"output"`
	OutputFile  string   `yaml:"output_file"`
	StateFile   string   `yaml:"state_file"`
	DryRun      bool     `yaml:"dry_run"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func MergeConfigWithFlags(cfg *Config, flags *Config) *Config {
	if flags.Region != "" {
		cfg.Region = flags.Region
	}
	if len(flags.Include) > 0 {
		cfg.Include = flags.Include
	}
	if flags.Output != "" {
		cfg.Output = flags.Output
	}
	if flags.OutputFile != "" {
		cfg.OutputFile = flags.OutputFile
	}
	if flags.StateFile != "" {
		cfg.StateFile = flags.StateFile
	}
	cfg.DryRun = cfg.DryRun || flags.DryRun

	return cfg
}

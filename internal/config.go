// Config file parser
package internal

import (
	"gopkg.in/yaml.v3"
	"os"
)

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

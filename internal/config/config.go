package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type JobArgs []string

type JobsList map[string][]JobArgs

type TriggerConfig struct {
	Type       string `yaml:"type"` // trigger type (git,test,etc)
	Job        string `yaml:"job"`
	GitRepoURL string `yaml:"git-repo-url"`
	GitBranch  string `yaml:"git-branch"`
}

type Configuration struct {
	Jobs     map[string]JobsList `yaml:"jobs"`
	Triggers []TriggerConfig     `yaml:"triggers"`
}

func LoadConfig(file *os.File) (*Configuration, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg Configuration
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

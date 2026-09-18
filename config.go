package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Env struct {
	HostedOn     string `yaml:"hosted_on"`
	AWSProfile   string `yaml:"aws_profile"`
	GCPProjectID string `yaml:"gcp_project_id"`
	InfraDir     string `yaml:"infra_dir"`
}

type BackendBucket struct {
	HostedOn   string `yaml:"hosted_on"`
	BucketName string `yaml:"bucket_name"`
	Region     string `yaml:"region"`
}

type Config struct {
	Envs           map[string]Env                    `yaml:"envs"`
	BackendBuckets map[string]BackendBucket          `yaml:"backend_buckets"`
	Tfvars         map[string]map[string]interface{} `yaml:"tfvars"`
}

// loadConfig reads, parses, and validates the YAML config file at the given path
func loadConfig(path string) (Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	if err := config.validate(); err != nil {
		return config, err
	}

	return config, nil
}

// validate checks the config is complete and internally consistent
func (c Config) validate() error {
	if len(c.Envs) == 0 {
		return fmt.Errorf("config has no environments defined under 'envs'")
	}

	for name, env := range c.Envs {
		if env.InfraDir == "" {
			return fmt.Errorf("env %q is missing infra_dir", name)
		}
		if env.HostedOn == "" {
			return fmt.Errorf("env %q is missing hosted_on", name)
		}

		if _, ok := c.BackendBuckets[name]; !ok {
			return fmt.Errorf("env %q has no matching entry in backend_buckets", name)
		}
		if _, ok := c.Tfvars[name]; !ok {
			return fmt.Errorf("env %q has no matching entry in tfvars", name)
		}
	}

	return nil
}

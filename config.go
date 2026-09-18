package main

import (
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
	Envs           map[string]Env                   `yaml:"envs"`
	BackendBuckets map[string]BackendBucket          `yaml:"backend_buckets"`
	Tfvars         map[string]map[string]interface{} `yaml:"tfvars"`
}

// loadConfig reads and parses the YAML config file at the given path
func loadConfig(path string) (Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}
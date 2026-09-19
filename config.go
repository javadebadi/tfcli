package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// HostedOn restricts which cloud providers are valid
type HostedOn string

const (
	AWS HostedOn = "AWS"
	GCP HostedOn = "GCP"
)

// isValidHostedOn checks whether this HostedOn value is one of the known providers
func (h HostedOn) isValidHostedOn() bool {
	switch h {
	case AWS, GCP:
		return true
	}
	return false
}

type Env struct {
	HostedOn     HostedOn `yaml:"hosted_on"`
	AWSProfile   string   `yaml:"aws_profile"`
	GCPProjectID string   `yaml:"gcp_project_id"`
	InfraDir     string   `yaml:"infra_dir"`
}

type BackendBucket struct {
	HostedOn   HostedOn `yaml:"hosted_on"`
	BucketName string   `yaml:"bucket_name"`
	Region     string   `yaml:"region"`
}

type Config struct {
	ProjectRoot    string                            `yaml:"project_root"`
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

func (c Config) validate() error {
	if err := c.validateEnvs(); err != nil {
		return err
	}
	if err := c.validateBackendBuckets(); err != nil {
		return err
	}
	return nil
}

func (c Config) validateEnvs() error {
	if len(c.Envs) == 0 {
		return fmt.Errorf("config has no environments defined under 'envs'")
	}

	for name, env := range c.Envs {
		if env.InfraDir == "" {
			return fmt.Errorf("env %q is missing infra_dir", name)
		}
		if !env.HostedOn.isValidHostedOn() {
			return fmt.Errorf("env %q has invalid hosted_on %q (must be AWS or GCP)", name, env.HostedOn)
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

func (c Config) validateBackendBuckets() error {

	for name, bucketInfo := range c.BackendBuckets {
		if bucketInfo.BucketName == "" {
			return fmt.Errorf("backend bucket for env %q has no bucket_name ", name)
		}
		if !bucketInfo.HostedOn.isValidHostedOn() {
			return fmt.Errorf("backend bucket for env %q has invalid hosted_on %q", name, bucketInfo.HostedOn)
		}
	}
	return nil
}

func (c Config) getBackendBucket(env string) (string, error) {
	bucketInfo, ok := c.BackendBuckets[env]
	if !ok {
		return "", fmt.Errorf("no bucket info for env %v", env)
	}
	return bucketInfo.BucketName, nil

}

func (c Config) getInfraDir(env string) (string, error) {
	infraDir := c.Envs[env].InfraDir
	if c.ProjectRoot != "" {
		infraDir = filepath.Join(c.ProjectRoot, c.Envs[env].InfraDir)
	}

	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		return "", fmt.Errorf("infra_dir %q does not exist", infraDir)
	}

	return infraDir, nil
}

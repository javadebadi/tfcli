package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Envs: map[string]Env{
					"prod": {HostedOn: "GCP", InfraDir: "terraform"},
				},
				BackendBuckets: map[string]BackendBucket{
					"prod": {HostedOn: "GCP", BucketName: "terraform-state-prod", Region: "us-central1"},
				},
				Tfvars: map[string]map[string]interface{}{
					"prod": {},
				},
			},
			wantErr: false,
		},
		{
			name: "config has no envs",
			config: Config{
				Envs: map[string]Env{
					// intentionally empty — "prod" key missing
				},
				BackendBuckets: map[string]BackendBucket{
					"prod": {HostedOn: "GCP", BucketName: "terraform-state-prod", Region: ""},
				},
				Tfvars: map[string]map[string]interface{}{
					"prod": {},
				},
			},
			wantErr: true,
		},
		{
			name: "missing infra_dir",
			config: Config{
				Envs: map[string]Env{
					"prod": {HostedOn: "GCP", InfraDir: ""},
				},
				BackendBuckets: map[string]BackendBucket{
					"prod": {},
				},
				Tfvars: map[string]map[string]interface{}{
					"prod": {},
				},
			},
			wantErr: true,
		},
		{
			name: "missing hosted_on",
			config: Config{
				Envs: map[string]Env{
					"prod": {HostedOn: "", InfraDir: "terraform"},
				},
				BackendBuckets: map[string]BackendBucket{
					"prod": {},
				},
				Tfvars: map[string]map[string]interface{}{
					"prod": {},
				},
			},
			wantErr: true,
		},
		{
			name: "existing but invalid hosted_on",
			config: Config{
				Envs: map[string]Env{
					"prod": {HostedOn: "ALIBABA", InfraDir: "terraform"},
				},
				BackendBuckets: map[string]BackendBucket{
					"prod": {},
				},
				Tfvars: map[string]map[string]interface{}{
					"prod": {},
				},
			},
			wantErr: true,
		},
		{
			name: "missing backend_buckets entry",
			config: Config{
				Envs: map[string]Env{
					"prod": {HostedOn: "GCP", InfraDir: "terraform"},
				},
				BackendBuckets: map[string]BackendBucket{
					// intentionally empty — "prod" key missing
				},
				Tfvars: map[string]map[string]interface{}{
					"prod": {},
				},
			},
			wantErr: true,
		},
		{
			name: "missing backend_buckets bucket_name",
			config: Config{
				Envs: map[string]Env{
					"prod": {HostedOn: "GCP", InfraDir: "terraform"},
				},
				BackendBuckets: map[string]BackendBucket{
					"prod": {HostedOn: "GCP", BucketName: "", Region: ""},
				},
				Tfvars: map[string]map[string]interface{}{
					"prod": {},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid backend_buckets hosted_on",
			config: Config{
				Envs: map[string]Env{
					"prod": {HostedOn: "GCP", InfraDir: "terraform"},
				},
				BackendBuckets: map[string]BackendBucket{
					"prod": {HostedOn: "BAD_HOSTED_ON", BucketName: "terraform-state-prod", Region: ""},
				},
				Tfvars: map[string]map[string]interface{}{
					"prod": {},
				},
			},
			wantErr: true,
		},
		{
			name: "missing tfvars entry",
			config: Config{
				Envs: map[string]Env{
					"prod": {HostedOn: "GCP", InfraDir: "terraform"},
				},
				BackendBuckets: map[string]BackendBucket{
					"prod": {HostedOn: "GCP", BucketName: "terraform-state-prod", Region: ""},
				},
				Tfvars: map[string]map[string]interface{}{
					// intentionally empty — "prod" key missing
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given a config from the table (tt.config)

			// When validate is called
			err := tt.config.validate()

			// Then it should match the expected error state
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	// Given a valid YAML file written to a temp directory
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	yamlContent := `
envs:
  prod:
    hosted_on: GCP
    infra_dir: terraform
backend_buckets:
  prod:
    hosted_on: GCP
    bucket_name: my-bucket
tfvars:
  prod:
    db_name: knaph
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// When loadConfig is called
	config, err := loadConfig(configPath)

	// Then it should succeed and parse correctly
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if config.Envs["prod"].GCPProjectID != "" {
		t.Errorf("unexpected GCPProjectID: %v", config.Envs["prod"].GCPProjectID)
	}
	if config.Tfvars["prod"]["db_name"] != "knaph" {
		t.Errorf("expected db_name=knaph, got: %v", config.Tfvars["prod"]["db_name"])
	}
}

func TestLoadConfig_ValidYaml_InvalidContent(t *testing.T) {
	// Given a valid YAML file written to a temp directory
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	yamlContent := `
envs:
  prod:
    hosted_on: INVALID
    infra_dir: terraform
backend_buckets:
  prod:
    hosted_on: INVALID
    bucket_name: my-bucket
tfvars:
  prod:
    db_name: knaph
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// When loadConfig is called
	_, err := loadConfig(configPath)

	// Then it should succeed and parse correctly
	if err == nil {
		t.Error("expected and error for valid YAML but invalid content of the yaml file")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	// Given a path that doesn't exist
	// When loadConfig is called
	_, err := loadConfig("/nonexistent/path/config.yaml")

	// Then it should return an error
	if err == nil {
		t.Error("expected an error for missing file, got nil")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Given a file with broken YAML syntax
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("not: valid: yaml: ["), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// When loadConfig is called
	_, err := loadConfig(configPath)

	// Then it should return an error
	if err == nil {
		t.Error("expected an error for invalid YAML, got nil")
	}
}

func TestGetBackendBucket(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		wantErr  bool
		expected string
	}{
		{
			name: "get backend bucket for a valid config",
			config: Config{
				Envs: map[string]Env{
					"prod": {HostedOn: "GCP", InfraDir: "terraform"},
				},
				BackendBuckets: map[string]BackendBucket{
					"prod": {HostedOn: "GCP", BucketName: "terraform-state-prod", Region: "us-central1"},
				},
				Tfvars: map[string]map[string]interface{}{
					"prod": {},
				},
			},
			wantErr:  false,
			expected: "terraform-state-prod",
		},
		{
			name: "Invalid Config",
			config: Config{
				Envs: map[string]Env{
					"prod": {HostedOn: "GCP", InfraDir: "terraform"},
				},
				BackendBuckets: map[string]BackendBucket{
					// intentionally empty — "prod" key missing
				},
				Tfvars: map[string]map[string]interface{}{
					"prod": {},
				},
			},
			wantErr:  true,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given a config from the table (tt.config)

			// When getBackendBucket is called to get backend bucket for "prod" environment
			val, err := tt.config.getBackendBucket(("prod"))

			// Then it should match the expected error state
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("getBackendBucket(\"prod\") error = %v, wantErr %v", err, tt.wantErr)
			}
			if val != tt.expected {
				t.Errorf("returned %v != expected %v", val, tt.expected)
			}
		})
	}

}

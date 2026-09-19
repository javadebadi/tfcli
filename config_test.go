package main

import "testing"

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

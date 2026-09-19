// tools to manage terraform backend bucket used to store state of terraform
package main

import "fmt"

func createGcpBucket(bucketName string, gcpProjectId string) error {
	// Runs this gsutil command
	// gsutil mb -p gcpProjectId gs://bucketName
	return runCommand("gsutil", "mb", "-p", gcpProjectId, fmt.Sprintf("gs://%s", bucketName))
}

func doesGcpBucketExist(bucketName string) (bool, error) {

	err := runCommand("gsutil", "ls", "-b", fmt.Sprintf("gs://%s", bucketName))

	if err != nil {
		return false, err
	}

	return true, nil
}

func createGcpBucketIfNotExists(bucketName string, gcpProjectId string) error {
	exists, err := doesGcpBucketExist(bucketName)

	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("bucket:%v already exists!", bucketName)
		return nil
	}

	return createGcpBucket(bucketName, gcpProjectId)
}

func createBackendBucket(c Config, env string) error {
	envConfig := c.Envs[env]

	switch envConfig.HostedOn {
	case GCP:
		return createGcpBucket(c.BackendBuckets[env].BucketName, envConfig.GCPProjectID)
	default:
		return fmt.Errorf("hosted_on %q not implemented", envConfig.HostedOn)
	}
}

func createBackendBucketIfNotExists(c Config, env string) error {
	envConfig := c.Envs[env]

	switch envConfig.HostedOn {
	case GCP:
		return createGcpBucketIfNotExists(c.BackendBuckets[env].BucketName, envConfig.GCPProjectID)
	default:
		return fmt.Errorf("hosted_on %q not implemented", envConfig.HostedOn)
	}
}

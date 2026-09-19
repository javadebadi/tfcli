package main

import (
	"fmt"
)

func tfInit(c Config, env string) error {
	infraDir, err := c.getInfraDir(env)
	if err != nil {
		return err
	}

	bucketName, err := c.getBackendBucket(env)
	if err != nil {
		return err
	}

	args := []string{
		"init",
		"--upgrade",
		"-reconfigure",
		fmt.Sprintf("-backend-config=bucket=%s", bucketName),
	}

	if c.Envs[env].HostedOn == GCP {
		args = append(args, "-backend-config=prefix=terraform/state")
	}

	return runCommandInDir(infraDir, "terraform", args...)
}

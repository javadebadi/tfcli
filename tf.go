package main

import (
	"fmt"
)

func tfInit(c Config, env string) error {
	bucketName, err := c.getBackendBucket(env)
	if err != nil {
		return err
	}

	args := []string{
		"--upgrade",
		"-reconfigure",
		fmt.Sprintf("-backend-config=bucket=%s", bucketName),
	}

	if c.Envs[env].HostedOn == GCP {
		args = append(args, "-backend-config=prefix=terraform/state")
	}

	return runTerraform(c, env, "init", args...)
}

func tfPlan(c Config, env string) error {
	args := []string{
		"-input=false",
		fmt.Sprintf("-var-file=%s", getEnvVarsFileName(env)),
	}
	return runTerraform(c, env, "plan", args...)
}

func tfValidate(c Config, env string) error {
	return runTerraform(c, env, "validate")
}

func tfAutoformat(c Config, env string) error {
	args := []string{
		"-recursive",
	}
	return runTerraform(c, env, "fmt", args...)
}

func tfCheckAutoformat(c Config, env string) error {
	// to check autoformat is run (can be used in CI pipelines)
	args := []string{
		"-recursive",
		"-check",
	}
	return runTerraform(c, env, "fmt", args...)
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// runCommand runs an external command in the current working directory
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runCommandInDir runs an external command inside the given directory
func runCommandInDir(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runTerraform runs a terraform subcommand (init, plan, apply...) with the
// given extra args, inside the environment's infra dir
func runTerraform(c Config, env string, subcommand string, extraArgs ...string) error {
	infraDir, err := c.getInfraDir(env)
	if err != nil {
		return err
	}

	args := append([]string{subcommand}, extraArgs...)
	return runCommandInDir(infraDir, "terraform", args...)
}

func envToLower(env string) string {
	return strings.ToLower(env)
}

func getEnvVarsFileName(env string) string {
	// returns filename of envvars
	return fmt.Sprintf("%s.tfvars", envToLower(env))
}

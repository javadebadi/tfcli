package main

import (
	"os"
	"os/exec"
)

// runCommand runs any external command, streaming its output to our terminal.
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

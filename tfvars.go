package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func writeTfvars(env string, infraDir string, vars map[string]interface{}) (err error) {
	if err := os.MkdirAll(infraDir, 0755); err != nil {
		return err
	}

	path := filepath.Join(infraDir, getEnvVarsFileName(env))

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	for key, value := range vars {
		if _, err := fmt.Fprintf(f, "%s = %q\n", key, value); err != nil {
			return err
		}
	}

	fmt.Println("wrote", path)
	return nil
}

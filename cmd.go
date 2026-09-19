package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	envName    string
	configPath string
	config     Config
)

var rootCmd = &cobra.Command{
	Use:     "tfctl",
	Short:   "A CLI to make working with Terraform environments easier",
	Version: version,

	// PersistentPreRun runs before this command's Run, AND before any
	// child subcommand's Run (plan, apply, etc.) — so tfvars are always
	// fresh before terraform ever sees them.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		for c := cmd; c != nil; c = c.Parent() {
			if c.Name() == "completion" || c.Name() == "help" || strings.HasPrefix(c.Name(), "__complete") {
				return
			}
		}

		if envName == "" {
			fmt.Println("--env is required")
			os.Exit(1)
		}

		var err error
		config, err = loadConfig(configPath)
		if err != nil {
			fmt.Println("error loading config:", err)
			os.Exit(1)
		}

		env, ok := config.Envs[envName]
		if !ok {
			fmt.Println("unknown environment:", envName)
			os.Exit(1)
		}

		if err := writeTfvars(envName, env.InfraDir, config.Tfvars[envName]); err != nil {
			fmt.Println("error writing tfvars:", err)
			os.Exit(1)
		}
	},
}

var tfvarsCmd = &cobra.Command{
	Use:   "tfvars",
	Short: "Generate a .tfvars file for the given environment",
	Run: func(cmd *cobra.Command, args []string) {
		// nothing to do here — PersistentPreRun already wrote the file
		fmt.Println("tfvars generated for", envName)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&envName, "env", "", "environment name (e.g. prod, staging)")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "config.yaml", "path to the config YAML file")

	rootCmd.AddCommand(tfvarsCmd)
}

// Execute runs the root command — called from main()
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

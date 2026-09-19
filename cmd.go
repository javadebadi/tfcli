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

		// check required flags exist
		if envName == "" {
			fmt.Println("--env is required")
			os.Exit(1)
		}

		if configPath == "" {
			fmt.Println("--config is required")
			os.Exit(1)
		}

		var err error
		config, err = loadConfig(configPath)
		if err != nil {
			fmt.Println("error loading config:", err)
			os.Exit(1)
		}

		_, ok := config.Envs[envName]
		if !ok {
			fmt.Println("unknown environment:", envName)
			os.Exit(1)
		}

		infraDir, err := config.getInfraDir(envName)
		if err != nil {
			fmt.Println("error in finding infraDir:", infraDir)
			os.Exit(1)
		}

		if err := writeTfvars(envName, infraDir, config.Tfvars[envName]); err != nil {
			fmt.Println("error writing tfvars:", err)
			os.Exit(1)
		}
	},
}

// a help to run terraform init before the taking another action
func runWithInit(c Config, env string, action func(Config, string) error) error {
	if err := tfInit(c, env); err != nil {
		return fmt.Errorf("terraform init failed: %w", err)
	}
	return action(c, env)
}

var tfvarsCmd = &cobra.Command{
	Use:   "tfvars",
	Short: "Generate a .tfvars file for the given environment",
	Run: func(cmd *cobra.Command, args []string) {
		// nothing to do here — PersistentPreRun already wrote the file
		fmt.Println("tfvars generated for", envName)
	},
}

var createBackendBucketCmd = &cobra.Command{
	Use:   "create-backend-bucket",
	Short: "Generate the terraform backend bucket for the given environment to store terraform state",
	Run: func(cmd *cobra.Command, args []string) {
		err := createBackendBucket(config, envName)
		if err == nil {
			fmt.Println("Backend bucket exists for env=", envName)
		} else {
			fmt.Println("Error in creating backend bucket: ", err)
		}
	},
}
var createBackendBucketIfNotExistsCmd = &cobra.Command{
	Use:   "create-backend-bucket-if-not-exists",
	Short: "If the backend bucket does not exist, it generates the terraform backend bucket for the given environment to store terraform state.",
	Run: func(cmd *cobra.Command, args []string) {
		err := createBackendBucketIfNotExists(config, envName)
		if err != nil {
			fmt.Println("Error in creating or getting backend bucket: ", err)
		}
	},
}
var tfInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Runs terraform init for a given environment",
	Run: func(cmd *cobra.Command, args []string) {
		err := tfInit(config, envName)
		if err != nil {
			fmt.Println("Error in terraform init: ", err)
		}
	},
}
var tfPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Runs terraform plan for a given environment",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runWithInit(config, envName, tfPlan); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}
var tfValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Runs terraform validate",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runWithInit(config, envName, tfValidate); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}
var tfAutoformatCmd = &cobra.Command{
	Use:   "autoformat",
	Short: "Runs terraform autoformat/fmt",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runWithInit(config, envName, tfAutoformat); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}
var tfCheckAutoformatCmd = &cobra.Command{
	Use:   "check-autoformat",
	Short: "Runs terraform autoformat/fmt to check format is correct",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runWithInit(config, envName, tfCheckAutoformat); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&envName, "env", "", "environment name (e.g. prod, staging)")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "path to the config YAML file")

	rootCmd.AddCommand(tfvarsCmd)
	rootCmd.AddCommand(createBackendBucketCmd)
	rootCmd.AddCommand(createBackendBucketIfNotExistsCmd)
	rootCmd.AddCommand(tfInitCmd)
	rootCmd.AddCommand(tfPlanCmd)
	rootCmd.AddCommand(tfValidateCmd)
	rootCmd.AddCommand(tfAutoformatCmd)
	rootCmd.AddCommand(tfCheckAutoformatCmd)
}

// Execute runs the root command — called from main()
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

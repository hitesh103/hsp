package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var envFlag string
var maskedFlag bool

var varCmd = &cobra.Command{
	Use:   "var",
	Short: "Manage variables",
	Long:  "List, set, delete, and manage environment variables",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all variables for an environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := LoadConfig()
		if err != nil {
			return err
		}

		envName := cfg.ActiveEnv
		if envFlag != "" {
			envName = envFlag
		}

		env, exists := cfg.Environments[envName]
		if !exists {
			return fmt.Errorf("environment %q not found", envName)
		}

		fmt.Printf("Environment: %s\n", envName)
		if len(env) == 0 {
			fmt.Println("  (no variables)")
			return nil
		}

		for key, value := range env {
			if maskedFlag {
				masked := MaskValue(key, value)
				if masked == "***" {
					fmt.Printf("  %-15s: %s (masked)\n", key, masked)
				} else {
					fmt.Printf("  %-15s: %s\n", key, value)
				}
			} else {
				fmt.Printf("  %-15s: %s\n", key, value)
			}
		}

		return nil
	},
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a variable in an environment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg, err := LoadConfig()
		if err != nil {
			return err
		}

		envName := cfg.ActiveEnv
		if envFlag != "" {
			envName = envFlag
		}

		if cfg.Environments[envName] == nil {
			cfg.Environments[envName] = make(map[string]string)
		}

		cfg.Environments[envName][key] = value

		if err := SaveConfig(); err != nil {
			return err
		}

		fmt.Printf("Set %s = %s in environment '%s'\n", key, value, envName)
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a variable from an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		cfg, err := LoadConfig()
		if err != nil {
			return err
		}

		envName := cfg.ActiveEnv
		if envFlag != "" {
			envName = envFlag
		}

		env, exists := cfg.Environments[envName]
		if !exists {
			return fmt.Errorf("environment %q not found", envName)
		}

		if _, exists := env[key]; !exists {
			return fmt.Errorf("variable %q not found in environment %q", key, envName)
		}

		delete(env, key)

		if err := SaveConfig(); err != nil {
			return err
		}

		fmt.Printf("Deleted %s from environment '%s'\n", key, envName)
		return nil
	},
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Print current config to stdout (YAML format)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := LoadConfig()
		if err != nil {
			return err
		}

		data, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}

		fmt.Print(string(data))
		return nil
	},
}

func init() {
	varCmd.AddCommand(listCmd)
	varCmd.AddCommand(setCmd)
	varCmd.AddCommand(deleteCmd)
	varCmd.AddCommand(exportCmd)

	varCmd.PersistentFlags().StringVar(&envFlag, "env", "", "target environment")
	listCmd.Flags().BoolVar(&maskedFlag, "masked", false, "mask secret values in list output")
}
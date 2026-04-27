package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var listEnv bool

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environments",
	Long:  "Switch, list, create, and delete environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := LoadConfig()
		if err != nil {
			return err
		}

		if listEnv {
			var envs []string
			for name := range cfg.Environments {
				envs = append(envs, name)
			}
			for _, name := range envs {
				marker := "  "
				current := ""
				if name == cfg.ActiveEnv {
					marker = "* "
					current = " (current)"
				}
				fmt.Printf("%s%s%s\n", marker, name, current)
			}
			return nil
		}

		if len(args) == 0 {
			var envs []string
			for name := range cfg.Environments {
				envs = append(envs, name)
			}
			fmt.Printf("Active environment: %s\n", cfg.ActiveEnv)
			fmt.Printf("Environments available: %s\n", strings.Join(envs, ", "))
			return nil
		}

		name := args[0]
		if _, exists := cfg.Environments[name]; !exists {
			fmt.Printf("Environment '%s' not found. Create it? (y/n): ", name)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if strings.ToLower(input) != "y" {
				fmt.Println("Aborted.")
				return nil
			}
			cfg.Environments[name] = make(map[string]string)
			if err := SaveConfig(); err != nil {
				return err
			}
			fmt.Printf("Created environment '%s'\n", name)
		}

		cfg.ActiveEnv = name
		if err := SaveConfig(); err != nil {
			return err
		}
		fmt.Printf("Switched to '%s'\n", name)
		return nil
	},
}

var createEnvCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := LoadConfig()
		if err != nil {
			return err
		}

		if _, exists := cfg.Environments[name]; exists {
			return fmt.Errorf("environment '%s' already exists", name)
		}

		cfg.Environments[name] = make(map[string]string)
		if err := SaveConfig(); err != nil {
			return err
		}

		fmt.Printf("Created environment '%s'\n", name)
		return nil
	},
}

var deleteEnvCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if name == "default" {
			return fmt.Errorf("cannot delete 'default' environment")
		}

		cfg, err := LoadConfig()
		if err != nil {
			return err
		}

		if _, exists := cfg.Environments[name]; !exists {
			return fmt.Errorf("environment '%s' not found", name)
		}

		delete(cfg.Environments, name)

		if cfg.ActiveEnv == name {
			cfg.ActiveEnv = "default"
		}

		if err := SaveConfig(); err != nil {
			return err
		}

		fmt.Printf("Deleted environment '%s'\n", name)
		return nil
	},
}

func init() {
	envCmd.Flags().BoolVar(&listEnv, "list", false, "list all environments")
	envCmd.AddCommand(createEnvCmd)
	envCmd.AddCommand(deleteEnvCmd)
}
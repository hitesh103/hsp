/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hsp",
	Short: "HTTP Superpowers - Easiest HTTP client in the terminal",
	Long: `HSP is an interactive HTTP client that makes API testing as easy as Postman, but in your terminal.

No need to remember curl syntax - just run 'hsp request' and answer simple prompts!

Features:
  • Interactive request builder - step-by-step guided flow
  • Auto-format JSON bodies and set Content-Type headers
  • Easy header and query parameter management
  • Request preview before sending
  • Automatic request history
  • Pretty-printed JSON responses

Examples:
  hsp request          - Start interactive request builder
  hsp get <url>        - Quick GET request
  hsp post <url>       - Quick POST request`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hsp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

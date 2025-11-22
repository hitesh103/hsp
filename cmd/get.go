package cmd

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"
)

var headers []string
var prettyOutput bool

var getCmd = &cobra.Command{
	Use:   "get [url]",
	Short: "Send a GET request",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: URL required")
			return
		}

		url := args[0]

		// Build request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Request Build Error:", err)
			return
		}

		// Apply custom headers
		for _, h := range headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}

		client := &http.Client{}

		// Measure request time
		start := time.Now()
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Request Error:", err)
			return
		}
		defer res.Body.Close()

		duration := time.Since(start)

		// Read body
		body, _ := io.ReadAll(res.Body)

		// Color status output
		statusColor := color.New(color.FgGreen).Add(color.Bold)
		if res.StatusCode >= 400 {
			statusColor = color.New(color.FgRed).Add(color.Bold)
		}

		statusColor.Printf("\nStatus: %d (%s)\n\n", res.StatusCode, duration)

		// Pretty print
		if prettyOutput {
		    formatted, err := prettyjson.Format(body)

		    if err == nil {
		        fmt.Println(string(formatted))
		    } else {
		        fmt.Println(string(body))
		    }
		}

	},
}

func init() {
	getCmd.Flags().StringArrayVarP(&headers, "header", "H", []string{}, "Custom request headers")
	getCmd.Flags().BoolVarP(&prettyOutput, "pretty", "p", true, "Pretty-print JSON output")
	rootCmd.AddCommand(getCmd)
}

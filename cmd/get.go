package cmd

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"encoding/json"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [url]",
	Short: "Send a GET request",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: URL required")
			return
		}

		url := args[0]

		start := time.Now()
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("Request Error:", err)
			return
		}
		defer res.Body.Close()

		duration := time.Since(start)
		body, _ := io.ReadAll(res.Body)

		fmt.Printf("\nStatus: %d (%s)\n\n", res.StatusCode, duration)

		var pretty map[string]interface{}
		if json.Unmarshal(body, &pretty) == nil {
			b, _ := json.MarshalIndent(pretty, "", "  ")
			fmt.Println(string(b))
		} else {
			fmt.Println(string(body))
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

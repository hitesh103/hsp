package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"
)

type RequestBuilder struct {
	URL          string
	Method       string
	Headers      map[string]string
	QueryParams  map[string]string
	Body         string
	BodyFormat   string
	PrettyOutput bool
}

var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Build and send HTTP requests interactively (like Postman in terminal)",
	Long:  "Interactive request builder - easiest way to make HTTP requests",
	Run: func(cmd *cobra.Command, args []string) {
		builder := NewRequestBuilder()
		builder.InteractiveFlow()
	},
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		Headers:      make(map[string]string),
		QueryParams:  make(map[string]string),
		Method:       "GET",
		PrettyOutput: true,
	}
}

func (rb *RequestBuilder) InteractiveFlow() {
	reader := bufio.NewReader(os.Stdin)

	// Step 1: Get URL
	rb.PromptURL(reader)

	// Step 2: Get Method
	rb.PromptMethod(reader)

	// Step 3: Add Headers
	rb.PromptHeaders(reader)

	// Step 4: Add Query Params
	rb.PromptQueryParams(reader)

	// Step 5: Handle Body for POST/PUT/PATCH
	if rb.Method == "POST" || rb.Method == "PUT" || rb.Method == "PATCH" {
		rb.PromptBody(reader)
	}

	// Step 6: Pretty print preference
	rb.PromptPrettyPrint(reader)

	// Step 7: Show preview
	rb.ShowPreview()

	// Step 8: Confirm and send
	if rb.ConfirmSend(reader) {
		rb.SendRequest()
	} else {
		fmt.Println("\n❌ Request cancelled")
	}
}

func (rb *RequestBuilder) PromptURL(reader *bufio.Reader) {
	for {
		fmt.Print("\n? URL: ")
		input, _ := reader.ReadString('\n')
		url := strings.TrimSpace(input)

		if url == "" {
			color.Red("✗ URL cannot be empty")
			continue
		}

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			color.Red("✗ URL must start with http:// or https://")
			continue
		}

		rb.URL = url
		break
	}
}

func (rb *RequestBuilder) PromptMethod(reader *bufio.Reader) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	fmt.Println("\n? Method: (default: GET)")
	for i, m := range methods {
		fmt.Printf("  %d) %s\n", i+1, m)
	}

	for {
		fmt.Print("Choose (1-7) or type method: ")
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		if choice == "" {
			rb.Method = "GET"
			color.Green("✓ Method: GET")
			break
		}

		// Check if numeric choice
		if len(choice) == 1 && choice >= "1" && choice <= "7" {
			idx := choice[0] - '1'
			rb.Method = methods[idx]
			color.Green("✓ Method: " + rb.Method)
			break
		}

		// Check if method name
		choice = strings.ToUpper(choice)
		valid := false
		for _, m := range methods {
			if m == choice {
				rb.Method = m
				valid = true
				color.Green("✓ Method: " + rb.Method)
				break
			}
		}

		if !valid {
			color.Red("✗ Invalid method")
			continue
		}
		break
	}
}

func (rb *RequestBuilder) PromptHeaders(reader *bufio.Reader) {
	fmt.Print("\n? Add headers? (y/n): ")
	input, _ := reader.ReadString('\n')

	if strings.ToLower(strings.TrimSpace(input)) != "y" {
		// Auto-set Accept header
		rb.Headers["Accept"] = "application/json"
		return
	}

	for {
		fmt.Print("  Header name (or 'done' to finish): ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		if strings.ToLower(key) == "done" {
			break
		}

		if key == "" {
			color.Red("  ✗ Header name cannot be empty")
			continue
		}

		fmt.Print("  Header value: ")
		value, _ := reader.ReadString('\n')
		value = strings.TrimSpace(value)

		rb.Headers[key] = value
		color.Green(fmt.Sprintf("  ✓ Added: %s: %s", key, value))
	}

	// Auto-set Accept if not present
	if _, exists := rb.Headers["Accept"]; !exists {
		rb.Headers["Accept"] = "application/json"
	}
}

func (rb *RequestBuilder) PromptQueryParams(reader *bufio.Reader) {
	fmt.Print("\n? Add query parameters? (y/n): ")
	input, _ := reader.ReadString('\n')

	if strings.ToLower(strings.TrimSpace(input)) != "y" {
		return
	}

	for {
		fmt.Print("  Parameter name (or 'done' to finish): ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		if strings.ToLower(key) == "done" {
			break
		}

		if key == "" {
			color.Red("  ✗ Parameter name cannot be empty")
			continue
		}

		fmt.Print("  Parameter value: ")
		value, _ := reader.ReadString('\n')
		value = strings.TrimSpace(value)

		rb.QueryParams[key] = value
		color.Green(fmt.Sprintf("  ✓ Added: %s=%s", key, value))
	}
}

func (rb *RequestBuilder) PromptBody(reader *bufio.Reader) {
	fmt.Print("\n? Add request body? (y/n): ")
	input, _ := reader.ReadString('\n')

	if strings.ToLower(strings.TrimSpace(input)) != "y" {
		return
	}

	fmt.Println("  Body format:")
	fmt.Println("  1) JSON")
	fmt.Println("  2) Form data")
	fmt.Println("  3) Raw text")

	for {
		fmt.Print("  Choose (1-3): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			rb.BodyFormat = "json"
			rb.Headers["Content-Type"] = "application/json"
			rb.PromptJSONBody(reader)
			return
		case "2":
			rb.BodyFormat = "form"
			rb.Headers["Content-Type"] = "application/x-www-form-urlencoded"
			rb.PromptFormBody(reader)
			return
		case "3":
			rb.BodyFormat = "raw"
			rb.PromptRawBody(reader)
			return
		default:
			color.Red("  ✗ Invalid choice")
		}
	}
}

func (rb *RequestBuilder) PromptJSONBody(reader *bufio.Reader) {
	fmt.Println("\n  Enter JSON body (press Enter twice when done):")
	var lines []string
	emptyLines := 0

	for emptyLines < 2 {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")

		if line == "" {
			emptyLines++
			continue
		}
		emptyLines = 0
		lines = append(lines, line)
	}

	jsonStr := strings.Join(lines, "\n")
	if jsonStr == "" {
		color.Red("  ✗ Body cannot be empty")
		return
	}

	// Validate JSON
	var obj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		color.Red(fmt.Sprintf("  ✗ Invalid JSON: %v", err))
		rb.PromptJSONBody(reader)
		return
	}

	// Pretty format the JSON
	formatted, _ := prettyjson.Format([]byte(jsonStr))
	rb.Body = string(formatted)
	color.Green("  ✓ JSON body set")
}

func (rb *RequestBuilder) PromptFormBody(reader *bufio.Reader) {
	params := make(url.Values)

	for {
		fmt.Print("  Form field name (or 'done' to finish): ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		if strings.ToLower(key) == "done" {
			break
		}

		if key == "" {
			color.Red("  ✗ Field name cannot be empty")
			continue
		}

		fmt.Print("  Field value: ")
		value, _ := reader.ReadString('\n')
		value = strings.TrimSpace(value)

		params.Add(key, value)
		color.Green(fmt.Sprintf("  ✓ Added: %s=%s", key, value))
	}

	rb.Body = params.Encode()
}

func (rb *RequestBuilder) PromptRawBody(reader *bufio.Reader) {
	fmt.Println("  Enter raw body (press Enter twice when done):")
	var lines []string
	emptyLines := 0

	for emptyLines < 2 {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")

		if line == "" {
			emptyLines++
			continue
		}
		emptyLines = 0
		lines = append(lines, line)
	}

	rb.Body = strings.Join(lines, "\n")
	if rb.Body == "" {
		color.Red("  ✗ Body cannot be empty")
		return
	}
	color.Green("  ✓ Raw body set")
}

func (rb *RequestBuilder) PromptPrettyPrint(reader *bufio.Reader) {
	fmt.Print("\n? Pretty-print response? (y/n, default: y): ")
	input, _ := reader.ReadString('\n')
	choice := strings.ToLower(strings.TrimSpace(input))

	rb.PrettyOutput = choice != "n"
}

func (rb *RequestBuilder) ShowPreview() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("PREVIEW")
	fmt.Println(strings.Repeat("=", 70))

	// Build full URL with query params
	fullURL := rb.URL
	if len(rb.QueryParams) > 0 {
		params := url.Values{}
		for k, v := range rb.QueryParams {
			params.Add(k, v)
		}
		fullURL = rb.URL + "?" + params.Encode()
	}

	// Method and URL
	methodColor := color.New(color.FgCyan, color.Bold)
	methodColor.Printf("%s ", rb.Method)
	fmt.Println(fullURL)

	// Headers
	if len(rb.Headers) > 0 {
		fmt.Println("\nHeaders:")
		for key, value := range rb.Headers {
			fmt.Printf("  %s: %s\n", color.GreenString(key), value)
		}
	}

	// Body
	if rb.Body != "" {
		fmt.Println("\nBody:")
		lines := strings.Split(rb.Body, "\n")
		for _, line := range lines {
			fmt.Printf("  %s\n", line)
		}
	}

	fmt.Println(strings.Repeat("=", 70))
}

func (rb *RequestBuilder) ConfirmSend(reader *bufio.Reader) bool {
	fmt.Print("\n? Send request? (y/n): ")
	input, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(input)) == "y"
}

func (rb *RequestBuilder) SendRequest() {
	// Build full URL with query params
	fullURL := rb.URL
	if len(rb.QueryParams) > 0 {
		params := url.Values{}
		for k, v := range rb.QueryParams {
			params.Add(k, v)
		}
		fullURL = rb.URL + "?" + params.Encode()
	}

	// Create request
	var req *http.Request
	var err error

	if rb.Body != "" {
		req, err = http.NewRequest(rb.Method, fullURL, strings.NewReader(rb.Body))
	} else {
		req, err = http.NewRequest(rb.Method, fullURL, nil)
	}

	if err != nil {
		color.Red(fmt.Sprintf("✗ Request creation failed: %v", err))
		return
	}

	// Add headers
	for key, value := range rb.Headers {
		req.Header.Set(key, value)
	}

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	start := time.Now()

	res, err := client.Do(req)
	if err != nil {
		color.Red(fmt.Sprintf("✗ Request failed: %v", err))
		return
	}
	defer res.Body.Close()

	duration := time.Since(start)

	// Read response
	body, _ := io.ReadAll(res.Body)

	// Print status
	fmt.Println()
	statusColor := color.New(color.FgGreen, color.Bold)
	if res.StatusCode >= 400 {
		statusColor = color.New(color.FgRed, color.Bold)
	}

	statusMsg := rb.GetStatusMessage(res.StatusCode)
	statusColor.Printf("✔ %d %s (%s)\n\n", res.StatusCode, statusMsg, duration)

	// Print response headers
	fmt.Println(color.BlueString("Response Headers:"))
	for key, values := range res.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", color.CyanString(key), value)
		}
	}

	// Print body
	fmt.Println("\n" + color.BlueString("Response Body:"))
	if rb.PrettyOutput {
		formatted, err := prettyjson.Format(body)
		if err == nil {
			fmt.Println(string(formatted))
		} else {
			fmt.Println(string(body))
		}
	} else {
		fmt.Println(string(body))
	}

	// Save to history
	rb.SaveToHistory()
}

func (rb *RequestBuilder) GetStatusMessage(code int) string {
	messages := map[int]string{
		200: "OK",
		201: "Created",
		204: "No Content",
		301: "Moved Permanently",
		302: "Found",
		304: "Not Modified",
		400: "Bad Request",
		401: "Unauthorized",
		403: "Forbidden",
		404: "Not Found",
		500: "Internal Server Error",
		502: "Bad Gateway",
		503: "Service Unavailable",
	}

	if msg, exists := messages[code]; exists {
		return msg
	}

	if code >= 200 && code < 300 {
		return "OK"
	} else if code >= 300 && code < 400 {
		return "Redirect"
	} else if code >= 400 && code < 500 {
		return "Client Error"
	} else if code >= 500 && code < 600 {
		return "Server Error"
	}

	return "Unknown"
}

func (rb *RequestBuilder) SaveToHistory() {
	// Create history directory if it doesn't exist
	historyDir := os.ExpandEnv("$HOME/.hsp/history")
	os.MkdirAll(historyDir, 0755)

	// Create history entry
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	historyFile := fmt.Sprintf("%s/%s_%s.json", historyDir, rb.Method, timestamp)

	historyEntry := map[string]interface{}{
		"timestamp": timestamp,
		"method":    rb.Method,
		"url":       rb.URL,
		"headers":   rb.Headers,
		"params":    rb.QueryParams,
		"body":      rb.Body,
	}

	data, _ := json.MarshalIndent(historyEntry, "", "  ")
	os.WriteFile(historyFile, data, 0644)

	color.Green(fmt.Sprintf("✓ Request saved to history: %s", historyFile))
}

func init() {
	rootCmd.AddCommand(requestCmd)
}

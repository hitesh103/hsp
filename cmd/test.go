package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type TestSuite struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Env         string            `json:"env"`
	Variables  map[string]string `json:"variables"`
	Tests       []TestCase        `json:"tests"`
}

type TestCase struct {
	Name        string         `json:"name"`
	Description string        `json:"description"`
	Request    TestRequest    `json:"request"`
	Assertions []Assertion    `json:"assertions"`
	SaveVar    []SaveVariable `json:"save"`
}

type TestRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

type Assertion struct {
	Type      string `json:"type"`
	Expected  string `json:"expected"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	Max       int    `json:"max"`
	Contains  string `json:"contains"`
	Value     string `json:"value"`
	Threshold int    `json:"threshold"`
}

type SaveVariable struct {
	Var  string `json:"var"`
	Path string `json:"path"`
}

type TestResult struct {
	Name        string
	Passed      bool
	Duration    time.Duration
	StatusCode  int
	ResponseBody string
	ResponseHeaders http.Header
	FailureReason string
	SavedVars map[string]string
}

var (
	testCmd = &cobra.Command{
		Use:   "test",
		Short: "Run and manage API test suites",
		Long:  "Run, list, and create API test suites in JSON format",
	}

	testRunCmd = &cobra.Command{
		Use:   "run <suite.json>",
		Short: "Run a test suite",
		Run:   runTestSuite,
	}

	testListCmd = &cobra.Command{
		Use:   "list",
		Short: "List available test suites",
		Run:   listTestSuites,
	}

	testCreateCmd = &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new test suite interactively",
		Run:   createTestSuite,
	}
)

func init() {
	testCmd.AddCommand(testRunCmd)
	testCmd.AddCommand(testListCmd)
	testCmd.AddCommand(testCreateCmd)
	rootCmd.AddCommand(testCmd)

	testRunCmd.Flags().StringP("env", "e", "", "Environment to use")
	testRunCmd.Flags().BoolP("stop-on-fail", "s", false, "Stop at first failure")
}

func SuitesDir() string {
	home := os.ExpandEnv("$HOME")
	dir := filepath.Join(home, ".hsp", "suites")
	if _, err := os.Stat(dir); err != nil {
		os.MkdirAll(dir, 0755)
	}
	return dir
}

func ResultsDir() string {
	home := os.ExpandEnv("$HOME")
	dir := filepath.Join(home, ".hsp", "test-results")
	if _, err := os.Stat(dir); err != nil {
		os.MkdirAll(dir, 0755)
	}
	return dir
}

func loadTestSuite(path string) (*TestSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var suite TestSuite
	if err := json.Unmarshal(data, &suite); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &suite, nil
}

func getVariables(envName string, suiteVars map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	if envName != "" {
		envVars, err := GetEnv(envName)
		if err != nil {
			return nil, fmt.Errorf("failed to get environment %q: %w", envName, err)
		}
		for k, v := range envVars {
			result[k] = v
		}
	} else {
		activeEnv, err := GetActiveEnv()
		if err == nil {
			for k, v := range activeEnv {
				result[k] = v
			}
		}
	}

	for k, v := range suiteVars {
		result[k] = v
	}

	return result, nil
}

func resolveTestVariables(input string, vars map[string]string) string {
	result := input
	for k, v := range vars {
		result = strings.ReplaceAll(result, "{{"+k+"}}", v)
	}
	return result
}

func runTestSuite(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		color.Red("✗ Please specify a test suite file")
		return
	}

	envFlag, _ := cmd.Flags().GetString("env")
	stopOnFail, _ := cmd.Flags().GetBool("stop-on-fail")

	suitePath := args[0]
	suite, err := loadTestSuite(suitePath)
	if err != nil {
		color.Red("✗ Failed to load test suite: %v", err)
		return
	}

	envName := envFlag
	if envName == "" && suite.Env != "" {
		envName = suite.Env
	}

	vars, err := getVariables(envName, suite.Variables)
	if err != nil {
		color.Red("✗ Failed to get variables: %v", err)
		return
	}

	displayEnv := envName
	if displayEnv == "" {
		displayEnv = "default"
	}

	fmt.Println(color.CyanString("+------------------------------------------------------------------------------+"))
	fmt.Printf(color.CyanString("|  TEST SUITE: %s"), suite.Name)
	fmt.Printf("%s", strings.Repeat(" ", 67-len(suite.Name)-12))
	fmt.Println(color.CyanString(fmt.Sprintf("[%s]       |", displayEnv)))
	fmt.Println(color.CyanString("+------------------------------------------------------------------------------+"))

	var results []TestResult
	hasFailure := false

	for i, tc := range suite.Tests {
		result := runTestCase(tc, vars, i+1)
		results = append(results, result)

		if result.Passed {
			dots := strings.Repeat(".", 50-len(tc.Name))
			fmt.Printf(color.GreenString("|  PASS  %s %s OK"), tc.Name, dots)
			fmt.Fprintf(color.Output, " %7.2fms\n", float64(result.Duration.Milliseconds()))
		} else {
			dots := strings.Repeat(".", 50-len(tc.Name))
			fmt.Printf(color.RedString("|  FAIL  %s %s FAIL"), tc.Name, dots)
			fmt.Fprintf(color.Output, " %7.2fms\n", float64(result.Duration.Milliseconds()))
			hasFailure = true
			if stopOnFail {
				break
			}
		}

		for k, v := range result.SavedVars {
			vars[k] = v
		}
	}

	fmt.Println(color.CyanString("+------------------------------------------------------------------------------+"))

	if hasFailure {
		fmt.Println(color.RedString("|  FAILURES:"))
		for idx, r := range results {
			if !r.Passed {
				fmt.Printf(color.RedString("|  [%d] %s\n"), idx+1, r.Name)
				fmt.Printf(color.RedString("|    Expected: %s\n"), r.FailureReason)
				fmt.Printf(color.RedString("|    Actual:   %d\n"), r.StatusCode)
			}
		}
		fmt.Println(color.CyanString("+------------------------------------------------------------------------------+"))
	}

	passed := 0
	for _, r := range results {
		if r.Passed {
			passed++
		}
	}

	total := len(results)
	if hasFailure {
		fmt.Printf(color.RedString("|  Summary: %d/%d passed                                              [FAIL]"), passed, total)
	} else {
		fmt.Printf(color.GreenString("|  Summary: %d/%d passed                                              [PASS]"), passed, total)
	}
	fmt.Println()
	fmt.Println(color.CyanString("+------------------------------------------------------------------------------+"))

	saveTestResults(suite.Name, results)
}

func runTestCase(tc TestCase, vars map[string]string, index int) TestResult {
	result := TestResult{
		Name:       tc.Name,
		Passed:     true,
		SavedVars:  make(map[string]string),
	}

	url := resolveTestVariables(tc.Request.URL, vars)
	method := tc.Request.Method
	if method == "" {
		method = "GET"
	}

	headers := make(map[string]string)
	for k, v := range tc.Request.Headers {
		headers[k] = resolveTestVariables(v, vars)
	}

	var body io.Reader
	if tc.Request.Body != nil {
		bodyStr, ok := tc.Request.Body.(string)
		if !ok {
			bodyBytes, _ := json.Marshal(tc.Request.Body)
			bodyStr = string(bodyBytes)
		}
		body = strings.NewReader(resolveTestVariables(bodyStr, vars))
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		result.Passed = false
		result.FailureReason = fmt.Sprintf("Request creation failed: %v", err)
		return result
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		result.Passed = false
		result.FailureReason = fmt.Sprintf("Request failed: %v", err)
		return result
	}
	defer resp.Body.Close()

	result.Duration = time.Since(start)
	result.StatusCode = resp.StatusCode
	result.ResponseHeaders = resp.Header

	bodyBytes, _ := io.ReadAll(resp.Body)
	result.ResponseBody = string(bodyBytes)

	for _, assertion := range tc.Assertions {
		assertionResult := checkAssertion(assertion, result)
		if !assertionResult.Passed {
			result.Passed = false
			result.FailureReason = assertionResult.Reason
			break
		}
	}

	if tc.SaveVar != nil && result.Passed {
		for _, sv := range tc.SaveVar {
			value := extractJSONPath(result.ResponseBody, sv.Path)
			if value != "" {
				result.SavedVars[sv.Var] = value
			}
		}
	}

	return result
}

type assertionResult struct {
	Passed bool
	Reason string
}

func checkAssertion(assertion Assertion, result TestResult) assertionResult {
	switch assertion.Type {
	case "status":
		expected := assertion.Expected
		if expected == "" {
			return assertionResult{Passed: false, Reason: "missing expected value"}
		}
		expectedCode := 0
		fmt.Sscanf(expected, "%d", &expectedCode)
		if result.StatusCode != expectedCode {
			return assertionResult{Passed: false, Reason: fmt.Sprintf("status %d", expectedCode)}
		}

	case "body_contains":
		path := assertion.Path
		value := assertion.Value
		if value == "" {
			value = assertion.Contains
		}
		extracted := extractJSONPath(result.ResponseBody, path)
		if !strings.Contains(extracted, value) {
			return assertionResult{Passed: false, Reason: fmt.Sprintf("body does not contain '%s' at %s", value, path)}
		}

	case "header":
		name := assertion.Name
		contains := assertion.Contains
		headerValue := result.ResponseHeaders.Get(name)
		if headerValue == "" {
			return assertionResult{Passed: false, Reason: fmt.Sprintf("header '%s' not found", name)}
		}
		if contains != "" && !strings.Contains(strings.ToLower(headerValue), strings.ToLower(contains)) {
			return assertionResult{Passed: false, Reason: fmt.Sprintf("header '%s' does not contain '%s'", name, contains)}
		}

	case "response_time_ms":
		max := assertion.Max
		if max == 0 {
			max = assertion.Threshold
		}
		if int(result.Duration.Milliseconds()) > max {
			return assertionResult{Passed: false, Reason: fmt.Sprintf("response time %dms > %dms", result.Duration.Milliseconds(), max)}
		}

	case "body_equals":
		path := assertion.Path
		expected := assertion.Expected
		if expected == "" {
			expected = assertion.Value
		}
		extracted := extractJSONPath(result.ResponseBody, path)
		if extracted != expected {
			return assertionResult{Passed: false, Reason: fmt.Sprintf("body at %s = '%s', expected '%s'", path, extracted, expected)}
		}
	}

	return assertionResult{Passed: true}
}

func extractJSONPath(body, path string) string {
	if path == "" {
		return body
	}

	path = strings.TrimPrefix(path, "$.")

	var data interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return ""
	}

	parts := strings.Split(path, ".")
	return navigateJSON(data, parts)
}

func navigateJSON(data interface{}, parts []string) string {
	current := data
	for _, part := range parts {
		if current == nil {
			return ""
		}

		if m, ok := current.(map[string]interface{}); ok {
			current = m[part]
		} else if arr, ok := current.([]interface{}); ok {
			idx := 0
			fmt.Sscanf(part, "%d", &idx)
			if idx >= 0 && idx < len(arr) {
				current = arr[idx]
			} else {
				return ""
			}
		} else {
			return ""
		}
	}

	if s, ok := current.(string); ok {
		return s
	}
	if n, ok := current.(float64); ok {
		return fmt.Sprintf("%v", n)
	}
	if b, ok := current.(bool); ok {
		return fmt.Sprintf("%v", b)
	}

	if bs, err := json.Marshal(current); err == nil {
		return string(bs)
	}

	return fmt.Sprintf("%v", current)
}

func saveTestResults(suiteName string, results []TestResult) {
	resultsDir := ResultsDir()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(resultsDir, fmt.Sprintf("%s_%s.json", suiteName, timestamp))

	data, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile(filename, data, 0644)
}

func listTestSuites(cmd *cobra.Command, args []string) {
	suitesDir := SuitesDir()

	entries, err := os.ReadDir(suitesDir)
	if err != nil {
		color.Red("✗ Failed to read suites directory: %v", err)
		return
	}

	if len(entries) == 0 {
		color.Yellow("No test suites found. Create one with: hsp test create <name>")
		return
	}

	fmt.Println(color.CyanString("Available test suites:"))
	fmt.Println()

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") {
			path := filepath.Join(suitesDir, entry.Name())
			data, _ := os.ReadFile(path)
			var suite TestSuite
			if err := json.Unmarshal(data, &suite); err != nil {
				fmt.Printf("  %s (invalid JSON)\n", entry.Name())
				continue
			}

			fmt.Printf("  %s\n", color.GreenString(entry.Name()))
			fmt.Printf("    Name: %s\n", suite.Name)
			if suite.Description != "" {
				fmt.Printf("    Description: %s\n", suite.Description)
			}
			fmt.Printf("    Tests: %d\n", len(suite.Tests))
			if suite.Env != "" {
				fmt.Printf("    Environment: %s\n", suite.Env)
			}
			fmt.Println()
		}
	}
}

func createTestSuite(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	var suite TestSuite
	suite.Variables = make(map[string]string)
	var input string

	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	if name == "" {
		fmt.Print("? Suite name: ")
		input, _ := reader.ReadString('\n')
		name = strings.TrimSpace(input)
	}
	suite.Name = name

	fmt.Print("? Description: ")
	input, _ = reader.ReadString('\n')
	suite.Description = strings.TrimSpace(input)

	var baseURL string
	fmt.Print("? Base URL (e.g., {{BASE_URL}}/users): ")
	input, _ = reader.ReadString('\n')
	baseURL = strings.TrimSpace(input)

	var method string
	fmt.Print("? Default method (GET/POST/PUT/DELETE): ")
	input, _ = reader.ReadString('\n')
	method = strings.TrimSpace(strings.ToUpper(input))
	if method == "" {
		method = "GET"
	}

	fmt.Print("? Add suite-specific variables? (y/n): ")
	input, _ = reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(input)) == "y" {
		for {
			fmt.Print("  Variable name (or 'done' to finish): ")
			key, _ := reader.ReadString('\n')
			key = strings.TrimSpace(key)
			if strings.ToLower(key) == "done" {
				break
			}
			if key == "" {
				continue
			}
			fmt.Print("  Variable value: ")
			value, _ := reader.ReadString('\n')
			suite.Variables[key] = strings.TrimSpace(value)
		}
	}

	fmt.Println("\n[Add test cases...]")

	for {
		fmt.Print("? Add another test case? (y/n): ")
		input, _ = reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(input)) != "y" {
			break
		}

		var tc TestCase
		tc.Request.Headers = make(map[string]string)

		fmt.Print("  ? Test name: ")
		input, _ = reader.ReadString('\n')
		tc.Name = strings.TrimSpace(input)

		fmt.Print("  ? Test description: ")
		input, _ = reader.ReadString('\n')
		tc.Description = strings.TrimSpace(input)

		fmt.Printf("  ? Method (default: %s): ", method)
		input, _ = reader.ReadString('\n')
		tc.Request.Method = strings.TrimSpace(strings.ToUpper(input))
		if tc.Request.Method == "" {
			tc.Request.Method = method
		}

		fmt.Printf("  ? Path (after base URL): ")
		input, _ = reader.ReadString('\n')
		path := strings.TrimSpace(input)
		tc.Request.URL = baseURL + path

		fmt.Print("  ? Add headers? (y/n): ")
		input, _ = reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(input)) == "y" {
			for {
				fmt.Print("    Header name (or 'done' to finish): ")
				key, _ := reader.ReadString('\n')
				key = strings.TrimSpace(key)
				if strings.ToLower(key) == "done" {
					break
				}
				if key == "" {
					continue
				}
				fmt.Print("    Header value: ")
				value, _ := reader.ReadString('\n')
				tc.Request.Headers[key] = strings.TrimSpace(value)
			}
		}

		if tc.Request.Method == "POST" || tc.Request.Method == "PUT" || tc.Request.Method == "PATCH" {
			fmt.Print("  ? Add request body? (y/n): ")
			input, _ = reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(input)) == "y" {
				fmt.Println("  Enter JSON body (press Enter twice when done):")
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
				body := strings.Join(lines, "\n")
				tc.Request.Body = body
			}
		}

		fmt.Print("  ? Add assertions? (y/n): ")
		input, _ = reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(input)) == "y" {
			for {
				fmt.Print("    Assertion type (status/body_contains/header/response_time_ms, or 'done'): ")
				input, _ = reader.ReadString('\n')
				assertType := strings.TrimSpace(strings.ToLower(input))
				if assertType == "done" {
					break
				}

				var assertion Assertion
				assertion.Type = assertType

				switch assertType {
				case "status":
					fmt.Print("    Expected status code: ")
					input, _ = reader.ReadString('\n')
					assertion.Expected = strings.TrimSpace(input)
				case "body_contains":
					fmt.Print("    JSON path (e.g., $.name): ")
					input, _ = reader.ReadString('\n')
					assertion.Path = strings.TrimSpace(input)
					fmt.Print("    Expected value: ")
					input, _ = reader.ReadString('\n')
					assertion.Value = strings.TrimSpace(input)
				case "header":
					fmt.Print("    Header name: ")
					input, _ = reader.ReadString('\n')
					assertion.Name = strings.TrimSpace(input)
					fmt.Print("    Must contain: ")
					input, _ = reader.ReadString('\n')
					assertion.Contains = strings.TrimSpace(input)
				case "response_time_ms":
					fmt.Print("    Max response time (ms): ")
					input, _ = reader.ReadString('\n')
					fmt.Sscanf(strings.TrimSpace(input), "%d", &assertion.Max)
				}

				if assertion.Type != "" {
					tc.Assertions = append(tc.Assertions, assertion)
				}
			}
		}

		fmt.Print("  ? Save variable from response? (y/n): ")
		input, _ = reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(input)) == "y" {
			for {
				fmt.Print("    Variable name to save as: ")
				input, _ = reader.ReadString('\n')
				varName := strings.TrimSpace(input)
				if varName == "" || strings.ToLower(varName) == "done" {
					break
				}
				fmt.Print("    JSON path to extract: ")
				input, _ = reader.ReadString('\n')
				path := strings.TrimSpace(input)

				tc.SaveVar = append(tc.SaveVar, SaveVariable{
					Var:  varName,
					Path: path,
				})
			}
		}

		suite.Tests = append(suite.Tests, tc)
	}

	data, _ := json.MarshalIndent(suite, "", "  ")
	filename := toFileName(suite.Name) + ".json"
	filepath := filepath.Join(SuitesDir(), filename)
	os.WriteFile(filepath, data, 0644)

	color.Green("\n✓ Test suite created: %s", filepath)
}

func toFileName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return strings.ToLower(re.ReplaceAllString(name, "-"))
}
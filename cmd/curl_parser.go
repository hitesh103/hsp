package cmd

import (
	"errors"
	"strings"
	"unicode"
)

// RequestConfig holds the parsed information from a cURL command
type RequestConfig struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

// ParseCurl parses a cURL command string and returns a RequestConfig
func ParseCurl(input string) (RequestConfig, error) {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "curl") {
		return RequestConfig{}, errors.New("input must start with 'curl'")
	}

	args, err := splitArgs(input)
	if err != nil {
		return RequestConfig{}, err
	}

	config := RequestConfig{
		Method:  "GET",
		Headers: make(map[string]string),
	}

	for i := 1; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-X", "--request":
			if i+1 < len(args) {
				config.Method = strings.ToUpper(args[i+1])
				i++
			}
		case "-H", "--header":
			if i+1 < len(args) {
				headerLine := args[i+1]
				parts := strings.SplitN(headerLine, ":", 2)
				if len(parts) == 2 {
					config.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
				i++
			}
		case "-d", "--data", "--data-raw", "--data-binary":
			if i+1 < len(args) {
				config.Body = args[i+1]
				if config.Method == "GET" {
					config.Method = "POST"
				}
				i++
			}
		default:
			if !strings.HasPrefix(arg, "-") {
				// Assume it's the URL if it doesn't start with -
				// (cURL can have URL anywhere, but usually it's at the end or after flags)
				// We'll take the first non-flag argument that isn't an argument to a flag
				if config.URL == "" {
					config.URL = strings.Trim(arg, "'\"")
				}
			}
		}
	}

	if config.URL == "" {
		return config, errors.New("could not find URL in cURL command")
	}

	return config, nil
}

// splitArgs splits a string into arguments, respecting quotes
func splitArgs(s string) ([]string, error) {
	var args []string
	var current strings.Builder
	inQuote := false
	var quoteChar rune

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if inQuote {
			if r == quoteChar {
				inQuote = false
			} else if r == '\\' && i+1 < len(runes) {
				i++
				current.WriteRune(runes[i])
			} else {
				current.WriteRune(r)
			}
		} else {
			if r == '"' || r == '\'' {
				inQuote = true
				quoteChar = r
			} else if unicode.IsSpace(r) {
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			} else {
				current.WriteRune(r)
			}
		}
	}

	if inQuote {
		return nil, errors.New("unclosed quote")
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args, nil
}

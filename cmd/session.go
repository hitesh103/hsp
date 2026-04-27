package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LastRequestJSON struct {
	URL          string            `json:"url"`
	Method      string            `json:"method"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"params"`
	Body        string            `json:"body"`
	BodyFormat  string            `json:"body_format"`
	CreatedAt   string            `json:"created_at"`
}

func GetLastRequestPath() string {
	homeDir := os.ExpandEnv("$HOME")
	return filepath.Join(homeDir, ".hsp", ".last_request.json")
}

func LoadLastRequest() (*LastRequestJSON, error) {
	path := GetLastRequestPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var lastReq LastRequestJSON
	if err := json.Unmarshal(data, &lastReq); err != nil {
		return nil, err
	}

	return &lastReq, nil
}

func SaveLastRequest(rb *RequestBuilder) error {
	path := GetLastRequestPath()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	lastReq := LastRequestJSON{
		URL:          rb.URL,
		Method:      rb.Method,
		Headers:     rb.Headers,
		QueryParams: rb.QueryParams,
		Body:        rb.Body,
		BodyFormat: rb.BodyFormat,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(lastReq, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}

func LoadLastRequestOrWarn() *LastRequestJSON {
	lastReq, err := LoadLastRequest()
	if err != nil {
		return nil
	}
	return lastReq
}

func MustLoadLastRequest() (*LastRequestJSON, error) {
	lastReq, err := LoadLastRequest()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no previous request found")
		}
		return nil, err
	}
	return lastReq, nil
}

func (rb *RequestBuilder) ApplyLastRequest(lastReq *LastRequestJSON) {
	rb.URL = lastReq.URL
	rb.Method = lastReq.Method
	rb.Headers = lastReq.Headers
	rb.QueryParams = lastReq.QueryParams
	rb.Body = lastReq.Body
	rb.BodyFormat = lastReq.BodyFormat
}
package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadAndSave(t *testing.T) {
	testDir := filepath.Join(os.TempDir(), "hsp-test")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	os.Setenv("HOME", testDir)
	config = nil

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.ActiveEnv != "default" {
		t.Errorf("expected default env, got %q", cfg.ActiveEnv)
	}

	if cfg.Environments["default"] == nil {
		t.Error("default environment should exist")
	}

	if err := SaveConfig(); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	if _, err := os.Stat(ConfigFile()); err != nil {
		t.Errorf("config file should exist after save: %v", err)
	}
}

func TestGetEnvGetActiveEnv(t *testing.T) {
	testDir := filepath.Join(os.TempDir(), "hsp-test2")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	os.Setenv("HOME", testDir)
	config = nil

	env, err := GetActiveEnv()
	if err != nil {
		t.Fatalf("GetActiveEnv failed: %v", err)
	}

	if env == nil {
		t.Error("env should not be nil")
	}

	_, err = GetEnv("nonexistent")
	if err == nil {
		t.Error("should return error for nonexistent env")
	}
}

func TestMaskValue(t *testing.T) {
	tests := []struct {
		key        string
		value     string
		wantMask  bool
	}{
		{"token", "secret", true},
		{"TOKEN", "secret", true},
		{"secret", "secret", true},
		{"SECRET", "secret", true},
		{"password", "secret", true},
		{"PASSWORD", "secret", true},
		{"api_key", "key", true},
		{"AUTH_TOKEN", "token", true},
		{"url", "http://example.com", false},
		{"name", "test", false},
	}

	for _, tt := range tests {
		got := MaskValue(tt.key, tt.value)
		isMasked := got == "***"
		if isMasked != tt.wantMask {
			t.Errorf("MaskValue(%q, %q) = %q, wantMask=%v", tt.key, tt.value, got, tt.wantMask)
		}
	}
}

func TestResolveVariables(t *testing.T) {
	vars := map[string]string{
		"BASE_URL": "https://api.example.com",
		"TOKEN":  "secret123",
	}

	tests := []struct {
		input string
		want  string
	}{
		{"{{BASE_URL}}/users", "https://api.example.com/users"},
		{"no-vars-here", "no-vars-here"},
	}

	for _, tt := range tests {
		got, missing := ResolveVariables(tt.input, vars)
		if got != tt.want {
			t.Errorf("ResolveVariables(%q) = %q, want %q", tt.input, got, tt.want)
		}
		if len(missing) > 0 && tt.input != "{{MISSING}}" {
			t.Errorf("unexpected missing vars: %v", missing)
		}
	}
}

func TestResolveAll(t *testing.T) {
	testDir := filepath.Join(os.TempDir(), "hsp-test3")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	os.Setenv("HOME", testDir)
	config = nil

	vars := map[string]string{
		"BASE_URL": "https://api.example.com",
		"TOKEN":  "abc123",
	}

	req := &RequestBuilder{
		URL:         "{{BASE_URL}}/users",
		Headers:     map[string]string{"Authorization": "Bearer {{TOKEN}}"},
		QueryParams: map[string]string{"format": "json"},
		Body:        `{"url": "{{BASE_URL}}"}`,
	}

	missing := ResolveAll(req, vars)

	if len(missing) > 0 {
		t.Errorf("ResolveAll should not report missing vars, got %v", missing)
	}

	if req.URL != "https://api.example.com/users" {
		t.Errorf("URL = %q, want %q", req.URL, "https://api.example.com/users")
	}

	authVal, exists := req.Headers["Authorization"]
	if !exists || authVal != "Bearer abc123" {
		t.Errorf("Authorization = %q, want %q", authVal, "Bearer abc123")
	}
}
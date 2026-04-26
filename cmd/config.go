package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environments map[string]map[string]string `yaml:"environments"`
	Aliases      map[string]string              `yaml:"aliases"`
	ActiveEnv    string                         `yaml:"activeEnv"`
}

var config *Config

func ConfigDir() string {
	home := os.ExpandEnv("$HOME")
	dir := filepath.Join(home, ".hsp")
	if _, err := os.Stat(dir); err != nil {
		os.MkdirAll(dir, 0755)
	}
	return dir
}

func ConfigFile() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func LoadConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	config = &Config{
		Environments: make(map[string]map[string]string),
		Aliases:      make(map[string]string),
		ActiveEnv:    "default",
	}

	config.Environments["default"] = make(map[string]string)

	path := ConfigFile()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if config.Environments == nil {
		config.Environments = make(map[string]map[string]string)
	}
	if config.Environments["default"] == nil {
		config.Environments["default"] = make(map[string]string)
	}
	if config.Aliases == nil {
		config.Aliases = make(map[string]string)
	}
	if config.ActiveEnv == "" {
		config.ActiveEnv = "default"
	}

	return config, nil
}

func SaveConfig() error {
	if config == nil {
		return fmt.Errorf("no config loaded")
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	path := ConfigFile()
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}

func GetEnv(name string) (map[string]string, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	env, exists := cfg.Environments[name]
	if !exists {
		return nil, fmt.Errorf("environment %q not found", name)
	}

	return env, nil
}

func GetActiveEnv() (map[string]string, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	env, exists := cfg.Environments[cfg.ActiveEnv]
	if !exists {
		return nil, fmt.Errorf("active environment %q not found", cfg.ActiveEnv)
	}

	return env, nil
}

func MaskValue(key, value string) string {
	lowerKey := strings.ToLower(key)
	sensitive := []string{"token", "secret", "password", "key", "auth"}
	for _, s := range sensitive {
		if strings.Contains(lowerKey, s) {
			return "***"
		}
	}
	return value
}
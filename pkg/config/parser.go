package config

import (
	"os"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"gopkg.in/yaml.v3"
)

// ParseFile reads and parses a sqlc.yaml file.
func ParseFile(path string) (*SqlcConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, apperrors.FileNotFoundError(path)
		}
		return nil, apperrors.FileReadError(path, err)
	}

	cfg, err := Parse(data)
	if err != nil {
		return nil, apperrors.ConfigParseError(path, err)
	}

	return cfg, nil
}

// Parse parses YAML data into a SqlcConfig.
func Parse(data []byte) (*SqlcConfig, error) {
	var cfg SqlcConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, apperrors.Wrapf(err, apperrors.ErrConfigParseFailed, "failed to parse YAML")
	}

	return &cfg, nil
}

// LoadOrDefault attempts to load a config file, returns default config if not found.
func LoadOrDefault(path string) (*SqlcConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	return ParseFile(path)
}

// DefaultConfig returns a basic default configuration.
func DefaultConfig() *SqlcConfig {
	return &SqlcConfig{
		Version: "2",
		SQL:     []SQLConfig{},
	}
}

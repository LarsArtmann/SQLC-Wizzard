package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Marshal converts a SqlcConfig to YAML bytes
func Marshal(cfg *SqlcConfig) ([]byte, error) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	return data, nil
}

// WriteFile writes a SqlcConfig to a file
func WriteFile(cfg *SqlcConfig, path string) error {
	data, err := Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// WriteFileFormatted writes a SqlcConfig with better formatting
func WriteFileFormatted(cfg *SqlcConfig, path string) error {
	// Use yaml.Encoder for better control over formatting
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
		defer func() {
		if err := file.Close(); err != nil {
			// Log the close error, but don't override the primary error
			fmt.Printf("warning: failed to close file %s: %v\n", path, err)
		}
	}()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2) // Use 2 spaces for indentation

	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return encoder.Close()
}

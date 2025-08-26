package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const configContent = `
`

// CreateConfigIfNotExists checks and creates config file if it doesn't exist
func CreateConfigIfNotExists(configPath string) error {
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		dir := filepath.Dir(configPath)
		if dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		}

		// File doesn't exist, create and write content
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
		fmt.Printf("Config file created: %s\n", configPath)
	} else if err != nil {
		// Other errors
		return fmt.Errorf("failed to check config file: %w", err)
	} else {
		fmt.Printf("Config file at: %s\n", configPath)
	}
	return nil
}

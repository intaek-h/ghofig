package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetConfigPath returns the path to write config to.
// On macOS, prefers the macOS-specific path if it exists (has override priority).
// Otherwise uses XDG path (cross-platform default).
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Check macOS-specific path first (if on macOS)
	if runtime.GOOS == "darwin" {
		macPath := filepath.Join(home, "Library", "Application Support", "com.mitchellh.ghostty", "config")
		if _, err := os.Stat(macPath); err == nil {
			return macPath, nil
		}
	}

	// Use XDG path (cross-platform)
	xdgPath := getXDGConfigPath(home)
	return xdgPath, nil
}

// getXDGConfigPath returns the XDG config path
func getXDGConfigPath(home string) string {
	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgHome == "" {
		xdgHome = filepath.Join(home, ".config")
	}
	return filepath.Join(xdgHome, "ghostty", "config")
}

// AppendLine appends a line to the config file.
// If the same option already exists, it comments out the old line(s) first.
// Creates the file and parent directories if they don't exist.
func AppendLine(line string) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure parent directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Extract option name from the new line
	optionName := ""
	if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
		optionName = strings.TrimSpace(parts[0])
	}

	// Read existing content and comment out matching options
	data, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var newContent strings.Builder
	if len(data) > 0 {
		lines := strings.Split(string(data), "\n")
		for i, existingLine := range lines {
			trimmed := strings.TrimSpace(existingLine)

			// Check if this line sets the same option (and isn't already commented)
			if optionName != "" && !strings.HasPrefix(trimmed, "#") {
				if parts := strings.SplitN(trimmed, "=", 2); len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					if key == optionName {
						// Comment out this line
						existingLine = "# " + existingLine
					}
				}
			}

			newContent.WriteString(existingLine)
			// Add newline except for last line if it was empty
			if i < len(lines)-1 {
				newContent.WriteString("\n")
			}
		}
	}

	// Ensure trailing newline before appending
	content := newContent.String()
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	// Append the new line
	content += line + "\n"

	// Write the entire file back
	return os.WriteFile(configPath, []byte(content), 0644)
}

// ConfigExists checks if a config file exists at any known location
func ConfigExists() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	// Check macOS path
	if runtime.GOOS == "darwin" {
		macPath := filepath.Join(home, "Library", "Application Support", "com.mitchellh.ghostty", "config")
		if _, err := os.Stat(macPath); err == nil {
			return true
		}
	}

	// Check XDG path
	xdgPath := getXDGConfigPath(home)
	_, err = os.Stat(xdgPath)
	return err == nil
}

// GetValue reads the current value for a config option from the config file.
// Returns empty string if not found.
func GetValue(optionName string) string {
	configPath, err := GetConfigPath()
	if err != nil {
		return ""
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	// Parse config file line by line, looking for the option
	// Later occurrences override earlier ones (Ghostty behavior)
	lines := strings.Split(string(data), "\n")
	var value string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse "option = value" format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			if key == optionName {
				value = strings.TrimSpace(parts[1])
			}
		}
	}

	return value
}

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	configDir      = ".fb"
	configFileName = "config.yaml"
	configDirPerm  = 0700 // User-only access for security (Story 5.1)
	configFilePerm = 0600
)

// Validation error messages
const (
	errAuthKeyRequired   = "auth_key is required in config file"
	errOrgIDRequired     = "org_id is required in config file"
	errUserEmailRequired = "user_email is required in config file"
)

// Config represents the application configuration
type Config struct {
	AuthKey   string `yaml:"auth_key"`
	OrgID     string `yaml:"org_id"`
	UserEmail string `yaml:"user_email"`
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, configDir, configFileName), nil
}

// LoadConfigFromPath reads configuration from a specific path
func LoadConfigFromPath(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, buildMissingConfigError(configPath)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		// Story 5.3: Enhance YAML syntax errors with helpful guidance
		return nil, EnhanceYAMLError(err)
	}

	return &cfg, nil
}

// buildMissingConfigError creates a helpful error message for missing config (Story 5.2)
func buildMissingConfigError(configPath string) error {
	return fmt.Errorf("config file not found at %s\n\n%s",
		configPath,
		GetFirstRunMessage(configPath))
}

// GetFirstRunMessage returns a helpful message for first-time users (Story 5.2)
func GetFirstRunMessage(configPath string) string {
	return fmt.Sprintf(`Welcome to fb! Let's get you set up.

To use this tool, create a configuration file at:
  %s

The configuration file needs these three fields:
  • auth_key  - Your Flow Boards API authentication key
  • org_id    - Your organization identifier
  • user_email - Your email address

Here's an example configuration you can use as a template:

auth_key: your-api-key-here
org_id: your-org-id
user_email: you@example.com

To obtain your API key and org ID, log into Flow Boards and check your
account settings or contact your administrator.

Once you've created the config file, run this command again to see your tickets!`, configPath)
}

// EnhanceYAMLError adds helpful context to YAML parsing errors (Story 5.3)
func EnhanceYAMLError(err error) error {
	return fmt.Errorf(`YAML syntax error in configuration file: %w

Common YAML mistakes to check:
  • Use spaces, not tabs, for indentation
  • Ensure consistent indentation (usually 2 spaces)
  • Check that each field has a colon followed by a space
  • Make sure quotes are properly matched

Here's an example of correct YAML format:

auth_key: your-api-key-here
org_id: your-org-id
user_email: you@example.com

You can check your YAML syntax at: https://www.yamllint.com/`, err)
}

// Validate checks that all required configuration fields are present
func (c *Config) Validate() error {
	if err := c.validateAuthKey(); err != nil {
		return err
	}
	if err := c.validateOrgID(); err != nil {
		return err
	}
	if err := c.validateUserEmail(); err != nil {
		return err
	}
	return nil
}

// validateAuthKey checks if the auth_key field is present
func (c *Config) validateAuthKey() error {
	if c.AuthKey == "" {
		return fmt.Errorf(errAuthKeyRequired)
	}
	return nil
}

// validateOrgID checks if the org_id field is present
func (c *Config) validateOrgID() error {
	if c.OrgID == "" {
		return fmt.Errorf(errOrgIDRequired)
	}
	return nil
}

// validateUserEmail checks if the user_email field is present
func (c *Config) validateUserEmail() error {
	if c.UserEmail == "" {
		return fmt.Errorf(errUserEmailRequired)
	}
	return nil
}

// LoadConfig reads the configuration from ~/.fb/config.yaml
func LoadConfig() (*Config, error) {
	// Story 5.1: Create config directory if it doesn't exist
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	if err := EnsureConfigDirectory(home); err != nil {
		return nil, err
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	cfg, err := LoadConfigFromPath(configPath)
	if err != nil {
		return nil, err
	}

	// Validate required fields (Story 1.3)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// SaveConfig writes the configuration to ~/.fb/config.yaml
func SaveConfig(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	if err := ensureConfigDirectory(configPath); err != nil {
		return err
	}

	data, err := marshalConfig(cfg)
	if err != nil {
		return err
	}

	return writeConfigFile(configPath, data)
}

// ensureConfigDirectory creates the config directory if it doesn't exist
func ensureConfigDirectory(configPath string) error {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, configDirPerm); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return nil
}

// EnsureConfigDirectory creates the ~/.fb directory if it doesn't exist (Story 5.1)
// It ensures the directory has user-only permissions (0700) for security.
// If the directory already exists, it is not modified.
func EnsureConfigDirectory(homeDir string) error {
	cfgDir := filepath.Join(homeDir, configDir)

	if err := validateConfigDirectoryPath(cfgDir); err != nil {
		return err
	}

	return createConfigDirectoryIfNeeded(cfgDir)
}

// validateConfigDirectoryPath checks if the config path exists and is valid
func validateConfigDirectoryPath(cfgDir string) error {
	info, err := os.Stat(cfgDir)
	if os.IsNotExist(err) {
		// Directory doesn't exist, which is fine - we'll create it
		return nil
	}
	if err != nil {
		// Some other stat error occurred
		return fmt.Errorf("unable to check configuration directory: %w", err)
	}

	// Path exists - verify it's a directory
	if !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", cfgDir)
	}

	// Directory exists and is valid - no action needed
	return nil
}

// createConfigDirectoryIfNeeded creates the config directory with secure permissions
func createConfigDirectoryIfNeeded(cfgDir string) error {
	// Check if directory already exists
	if _, err := os.Stat(cfgDir); err == nil {
		// Directory exists, don't modify it
		return nil
	}

	// Create directory with user-only permissions (0700)
	if err := os.Mkdir(cfgDir, configDirPerm); err != nil {
		return fmt.Errorf("unable to create configuration directory at %s - check permissions: %w", cfgDir, err)
	}

	return nil
}

// marshalConfig converts the config struct to YAML bytes
func marshalConfig(cfg *Config) ([]byte, error) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}
	return data, nil
}

// writeConfigFile writes the config data to the specified path
func writeConfigFile(configPath string, data []byte) error {
	if err := os.WriteFile(configPath, data, configFilePerm); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

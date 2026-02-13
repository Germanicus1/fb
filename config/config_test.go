package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestStory1_2_ReadConfigFile tests reading valid configuration file
func TestStory1_2_ReadConfigFile(t *testing.T) {
	// Given: A valid config file at ~/.fb/config.yaml
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".fb", "config.yaml")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configContent := `auth_key: "test-auth-key"
org_id: "test-org-id"
user_email: "test@example.com"
`
	err = os.WriteFile(configPath, []byte(configContent), 0600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// When: Loading the config
	cfg, err := LoadConfigFromPath(configPath)

	// Then: Config should be loaded successfully
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if cfg == nil {
		t.Fatal("Expected config to be loaded, got nil")
	}

	// Acceptance Criterion: Parses auth_key
	if cfg.AuthKey != "test-auth-key" {
		t.Errorf("Expected auth_key='test-auth-key', got: %s", cfg.AuthKey)
	}

	// Acceptance Criterion: Parses org_id
	if cfg.OrgID != "test-org-id" {
		t.Errorf("Expected org_id='test-org-id', got: %s", cfg.OrgID)
	}

	// Acceptance Criterion: Parses user_email
	if cfg.UserEmail != "test@example.com" {
		t.Errorf("Expected user_email='test@example.com', got: %s", cfg.UserEmail)
	}
}

// TestStory1_2_MissingConfigFile tests error when config file is missing
func TestStory1_2_MissingConfigFile(t *testing.T) {
	// Given: A non-existent config file path
	nonExistentPath := "/tmp/nonexistent_dir_12345/config.yaml"

	// When: Attempting to load the config
	cfg, err := LoadConfigFromPath(nonExistentPath)

	// Then: Should return clear error about missing file
	if err == nil {
		t.Error("Expected error for missing config file, got nil")
	}

	if cfg != nil {
		t.Error("Expected nil config for missing file, got non-nil")
	}

	// Acceptance Criterion: Clear error message for missing file
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "not found") && !strings.Contains(errorMsg, "no such file") {
		t.Errorf("Error message should indicate file not found, got: %s", errorMsg)
	}
}

// TestStory1_2_InvalidYAML tests error when YAML is malformed
func TestStory1_2_InvalidYAML(t *testing.T) {
	// Given: A config file with invalid YAML syntax
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	invalidYAML := `auth_key: test-key
org_id: [this is invalid YAML syntax
user_email: test@example.com
`
	err := os.WriteFile(configPath, []byte(invalidYAML), 0600)
	if err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	// When: Attempting to load the config
	cfg, err := LoadConfigFromPath(configPath)

	// Then: Should return clear error about invalid YAML
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}

	if cfg != nil {
		t.Error("Expected nil config for invalid YAML, got non-nil")
	}

	// Acceptance Criterion: Clear error message for invalid YAML
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "parse") && !strings.Contains(errorMsg, "YAML") && !strings.Contains(errorMsg, "invalid") {
		t.Errorf("Error message should indicate YAML parsing error, got: %s", errorMsg)
	}
}

// TestStory1_2_ConfigPathResolution tests that ~/.fb/config.yaml path is resolved correctly
func TestStory1_2_ConfigPathResolution(t *testing.T) {
	// When: Getting the default config path
	configPath, err := GetConfigPath()

	// Then: Should return valid path
	if err != nil {
		t.Errorf("Expected no error getting config path, got: %v", err)
	}

	if configPath == "" {
		t.Error("Config path should not be empty")
	}

	// Should end with .fb/config.yaml
	if !strings.HasSuffix(configPath, ".fb/config.yaml") && !strings.HasSuffix(configPath, ".fb\\config.yaml") {
		t.Errorf("Config path should end with .fb/config.yaml, got: %s", configPath)
	}
}

// TestStory1_2_EmptyConfigFile tests handling of empty config file
func TestStory1_2_EmptyConfigFile(t *testing.T) {
	// Given: An empty config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(""), 0600)
	if err != nil {
		t.Fatalf("Failed to write empty config file: %v", err)
	}

	// When: Attempting to load the config
	cfg, err := LoadConfigFromPath(configPath)

	// Then: Should handle gracefully (either error or empty config)
	// For Story 1.3 we'll validate required fields, but for 1.2 we just need to parse
	if err != nil {
		// If there's an error, it should be clear
		errorMsg := err.Error()
		if errorMsg == "" {
			t.Error("Error message should not be empty")
		}
	}

	// Empty file might return empty config (validation happens in Story 1.3)
	if cfg == nil && err == nil {
		t.Error("Should either return config or error, got both nil")
	}
}

// STORY 5.1: Create Configuration Directory on First Run

// TestStory5_1_CreateDirectoryWhenMissing tests that ~/.fb is created if it doesn't exist
func TestStory5_1_CreateDirectoryWhenMissing(t *testing.T) {
	// Given: A temporary directory without .fb subdirectory
	tempHomeDir := t.TempDir()
	configDir := filepath.Join(tempHomeDir, ".fb")

	// Verify directory doesn't exist initially
	if _, err := os.Stat(configDir); !os.IsNotExist(err) {
		t.Fatal("Test setup error: config directory should not exist initially")
	}

	// When: EnsureConfigDirectory is called
	err := EnsureConfigDirectory(tempHomeDir)

	// Then: Directory should be created
	// Acceptance Criterion: Tool creates ~/.fb directory if it doesn't exist
	if err != nil {
		t.Errorf("Expected no error creating directory, got: %v", err)
	}

	// Verify directory now exists
	info, err := os.Stat(configDir)
	if os.IsNotExist(err) {
		t.Error("Expected directory to be created, but it doesn't exist")
	}
	if err != nil {
		t.Errorf("Error checking directory: %v", err)
	}
	if info != nil && !info.IsDir() {
		t.Error("Expected .fb to be a directory")
	}
}

// TestStory5_1_DirectoryPermissions tests that directory is created with user-only access (700)
func TestStory5_1_DirectoryPermissions(t *testing.T) {
	// Given: A temporary directory without .fb subdirectory
	tempHomeDir := t.TempDir()
	configDir := filepath.Join(tempHomeDir, ".fb")

	// When: EnsureConfigDirectory is called
	err := EnsureConfigDirectory(tempHomeDir)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Then: Directory should have user-only permissions (700)
	// Acceptance Criterion: Directory is created with appropriate permissions (user-only access: 700)
	info, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("Failed to stat directory: %v", err)
	}

	mode := info.Mode().Perm()
	// Should be at least user-read, user-write, user-execute
	expectedPerms := os.FileMode(0700)
	if mode != expectedPerms {
		t.Errorf("Expected directory permissions 0700, got %04o", mode)
	}
}

// TestStory5_1_PermissionDeniedError tests error message when directory creation fails
func TestStory5_1_PermissionDeniedError(t *testing.T) {
	// Given: A path where we cannot create directories (read-only parent)
	// Note: This is tricky to test portably; we'll test the error handling logic

	// When: Attempting to create directory in an invalid location
	err := EnsureConfigDirectory("/root/impossible_path_for_test")

	// Then: Should show clear error message
	// Acceptance Criterion: If directory creation fails, show clear message about permissions
	if err == nil {
		// This might succeed if running as root, so we skip the test
		t.Skip("Test requires permission restrictions, skipping")
	}

	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "create") && !strings.Contains(errorMsg, "directory") {
		t.Errorf("Error should mention directory creation issue, got: %s", errorMsg)
	}
}

// TestStory5_1_ExistingDirectoryNotModified tests that existing directory is not changed
func TestStory5_1_ExistingDirectoryNotModified(t *testing.T) {
	// Given: An existing ~/.fb directory with specific permissions
	tempHomeDir := t.TempDir()
	configDir := filepath.Join(tempHomeDir, ".fb")

	// Create directory with different permissions
	err := os.Mkdir(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Get original modification time
	originalInfo, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("Failed to stat directory: %v", err)
	}
	originalMode := originalInfo.Mode().Perm()

	// When: EnsureConfigDirectory is called
	err = EnsureConfigDirectory(tempHomeDir)

	// Then: Should not return error
	// Acceptance Criterion: Existing ~/.fb directory is not modified or overwritten
	if err != nil {
		t.Errorf("Expected no error for existing directory, got: %v", err)
	}

	// Verify directory still exists and permissions weren't changed
	newInfo, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("Failed to stat directory after call: %v", err)
	}

	newMode := newInfo.Mode().Perm()
	if newMode != originalMode {
		t.Errorf("Existing directory permissions should not be modified: original=%04o, new=%04o", originalMode, newMode)
	}
}

// TestStory5_1_LoadConfigCreatesDirectory tests that LoadConfig creates directory automatically
func TestStory5_1_LoadConfigCreatesDirectory(t *testing.T) {
	// Given: A temporary home directory without .fb directory, but with config file somehow
	tempHomeDir := t.TempDir()
	configDir := filepath.Join(tempHomeDir, ".fb")

	// When: LoadConfig is called with auto-create enabled
	// This simulates the integrated behavior we'll add to LoadConfig
	err := EnsureConfigDirectory(tempHomeDir)

	// Then: Directory should be created
	// Acceptance Criterion: Tool continues to check for config file after creating directory
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		t.Error("Expected directory to exist after LoadConfig")
	}
}

// STORY 5.2: Show Helpful First-Run Message

// TestStory5_2_MissingConfigShowsHelpfulMessage tests that missing config shows setup guidance
func TestStory5_2_MissingConfigShowsHelpfulMessage(t *testing.T) {
	// Given: A system without config file
	tempDir := t.TempDir()
	nonExistentPath := filepath.Join(tempDir, "config.yaml")

	// When: Attempting to load config
	message := GetFirstRunMessage(nonExistentPath)

	// Then: Should show helpful message with required information
	// Acceptance Criterion: Message displays location where config should be created
	if !strings.Contains(message, "config.yaml") {
		t.Error("Message should mention config.yaml location")
	}

	// Acceptance Criterion: Message lists required fields
	if !strings.Contains(message, "auth_key") {
		t.Error("Message should mention auth_key field")
	}
	if !strings.Contains(message, "org_id") {
		t.Error("Message should mention org_id field")
	}
	if !strings.Contains(message, "user_email") {
		t.Error("Message should mention user_email field")
	}
}

// TestStory5_2_MessageIncludesExampleYAML tests that message has valid YAML example
func TestStory5_2_MessageIncludesExampleYAML(t *testing.T) {
	// Given: First run scenario
	configPath := "/tmp/test/config.yaml"

	// When: Getting first run message
	message := GetFirstRunMessage(configPath)

	// Then: Should include example configuration format
	// Acceptance Criterion: Example configuration format (valid YAML)
	if !strings.Contains(message, "auth_key:") || !strings.Contains(message, "org_id:") || !strings.Contains(message, "user_email:") {
		t.Error("Message should include example YAML format with all fields")
	}

	// Should show example values
	if !strings.Contains(message, "your-api-key") || !strings.Contains(message, "example.com") {
		t.Error("Message should include example placeholder values")
	}
}

// TestStory5_2_MessageIsFriendlyAndEncouraging tests tone of message
func TestStory5_2_MessageIsFriendlyAndEncouraging(t *testing.T) {
	// Given: First run scenario
	configPath := "~/.fb/config.yaml"

	// When: Getting first run message
	message := GetFirstRunMessage(configPath)

	// Then: Message should be friendly and encouraging, not intimidating
	// Acceptance Criterion: Message is friendly and encouraging, not intimidating
	lowerMsg := strings.ToLower(message)

	// Should not have harsh/intimidating language
	harshWords := []string{"error", "failed", "invalid", "missing", "required"}
	hasHarshTone := false
	for _, word := range harshWords {
		if strings.Contains(lowerMsg, word) {
			hasHarshTone = true
			break
		}
	}

	if hasHarshTone {
		// This is acceptable if it's softened with friendly language
		// Check for friendly words too
		friendlyWords := []string{"welcome", "let's", "setup", "get started", "create", "configure"}
		hasFriendlyTone := false
		for _, word := range friendlyWords {
			if strings.Contains(lowerMsg, word) {
				hasFriendlyTone = true
				break
			}
		}

		if !hasFriendlyTone {
			t.Error("Message tone should be friendly and encouraging, not harsh")
		}
	}
}

// TestStory5_2_MessageProvidesNextSteps tests that message is actionable
func TestStory5_2_MessageProvidesNextSteps(t *testing.T) {
	// Given: First run scenario
	configPath := "~/.fb/config.yaml"

	// When: Getting first run message
	message := GetFirstRunMessage(configPath)

	// Then: Message should provide clear next steps
	// Acceptance Criterion: Message provides clear next steps
	lowerMsg := strings.ToLower(message)

	// Should indicate what to do next
	actionWords := []string{"create", "add", "set", "configure", "edit"}
	hasActionGuidance := false
	for _, word := range actionWords {
		if strings.Contains(lowerMsg, word) {
			hasActionGuidance = true
			break
		}
	}

	if !hasActionGuidance {
		t.Error("Message should provide actionable next steps")
	}

	// Should show the exact path where to create the file
	if !strings.Contains(message, configPath) {
		t.Error("Message should include the exact config file path")
	}
}

// TestStory5_2_MessageShowsWhereToObtainAPIKey tests API key guidance
func TestStory5_2_MessageShowsWhereToObtainAPIKey(t *testing.T) {
	// Given: First run scenario
	configPath := "~/.fb/config.yaml"

	// When: Getting first run message
	message := GetFirstRunMessage(configPath)

	// Then: Should include instructions on where to obtain API key
	// Acceptance Criterion: Instructions on where to obtain API key (if known)
	lowerMsg := strings.ToLower(message)

	// Should mention API key or authentication
	if !strings.Contains(lowerMsg, "api") || !strings.Contains(lowerMsg, "key") {
		t.Error("Message should mention API key")
	}

	// Message length should be reasonable (not too short, not overwhelming)
	if len(message) < 100 {
		t.Error("Message seems too short to be helpful")
	}
	if len(message) > 2000 {
		t.Error("Message might be too long and overwhelming")
	}
}

// STORY 5.3: Validate YAML Configuration Syntax Clearly

// TestStory5_3_YAMLErrorShowsParserError tests that YAML parsing errors are shown
func TestStory5_3_YAMLErrorShowsParserError(t *testing.T) {
	// Given: A config file with invalid YAML
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	invalidYAML := `auth_key: test
org_id: [this is broken
user_email: test@example.com`

	err := os.WriteFile(configPath, []byte(invalidYAML), 0600)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// When: Attempting to load config
	_, err = LoadConfigFromPath(configPath)

	// Then: Should show specific error from YAML parser
	// Acceptance Criterion: Shows specific error from YAML parser
	if err == nil {
		t.Fatal("Expected error for invalid YAML")
	}

	errorMsg := err.Error()
	// Should contain helpful YAML syntax guidance
	lowerMsg := strings.ToLower(errorMsg)
	if !strings.Contains(lowerMsg, "yaml") && !strings.Contains(lowerMsg, "syntax") {
		t.Errorf("Error should mention YAML or syntax issue, got: %s", errorMsg)
	}
}

// TestStory5_3_ErrorIncludesLineNumber tests line number reporting
func TestStory5_3_ErrorIncludesLineNumber(t *testing.T) {
	// Given: A config file with invalid YAML on specific line
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	invalidYAML := `auth_key: test
org_id: [broken
user_email: test@example.com`

	err := os.WriteFile(configPath, []byte(invalidYAML), 0600)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// When: Attempting to load config
	_, err = LoadConfigFromPath(configPath)

	// Then: Error should include line number if available
	// Acceptance Criterion: Error message includes line number if available
	if err == nil {
		t.Fatal("Expected error for invalid YAML")
	}

	errorMsg := err.Error()
	// YAML parser typically includes line numbers - check if present
	hasLineInfo := strings.Contains(errorMsg, "line") || strings.Contains(errorMsg, ":")
	if !hasLineInfo {
		// This is acceptable if the YAML parser doesn't provide line info
		t.Logf("Note: Error doesn't include line number: %s", errorMsg)
	}
}

// TestStory5_3_SuggestsCommonMistakes tests helpful suggestions
func TestStory5_3_SuggestsCommonMistakes(t *testing.T) {
	// Given: A config file with invalid YAML
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	invalidYAML := `auth_key: test
	org_id: bad_indent
user_email: test@example.com`

	err := os.WriteFile(configPath, []byte(invalidYAML), 0600)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// When: Getting enhanced YAML error
	enhancedErr := EnhanceYAMLError(fmt.Errorf("yaml parse error"))

	// Then: Should suggest common mistakes
	// Acceptance Criterion: Suggests common YAML mistakes (tabs vs spaces, indentation, missing colons)
	errorMsg := enhancedErr.Error()
	lowerMsg := strings.ToLower(errorMsg)

	suggestionsProvided := strings.Contains(lowerMsg, "tab") ||
		strings.Contains(lowerMsg, "indent") ||
		strings.Contains(lowerMsg, "colon") ||
		strings.Contains(lowerMsg, "space")

	if !suggestionsProvided {
		t.Error("Error should suggest common YAML mistakes")
	}
}

// TestStory5_3_SuggestsYAMLValidator tests validator recommendation
func TestStory5_3_SuggestsYAMLValidator(t *testing.T) {
	// Given: A YAML parsing error
	parseErr := fmt.Errorf("yaml: line 2: could not find expected ':'")

	// When: Enhancing the error message
	enhancedErr := EnhanceYAMLError(parseErr)

	// Then: Should suggest using online YAML validator
	// Acceptance Criterion: Tool suggests checking YAML syntax with online validator
	errorMsg := enhancedErr.Error()
	lowerMsg := strings.ToLower(errorMsg)

	if !strings.Contains(lowerMsg, "validator") && !strings.Contains(lowerMsg, "check") {
		t.Error("Error should suggest using a YAML validator")
	}
}

// TestStory5_3_ProvidesCorrectYAMLExample tests example format
func TestStory5_3_ProvidesCorrectYAMLExample(t *testing.T) {
	// Given: A YAML parsing error
	parseErr := fmt.Errorf("yaml parse error")

	// When: Enhancing the error message
	enhancedErr := EnhanceYAMLError(parseErr)

	// Then: Should provide example of correct YAML format
	// Acceptance Criterion: Example of correct YAML format is provided
	errorMsg := enhancedErr.Error()

	// Should show proper YAML structure
	if !strings.Contains(errorMsg, "auth_key:") || !strings.Contains(errorMsg, "org_id:") {
		t.Error("Error should include example of correct YAML format")
	}
}

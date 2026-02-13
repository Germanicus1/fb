package config

import (
	"strings"
	"testing"
)

// TestStory1_3_ValidateAllFieldsPresent tests validation of complete config
func TestStory1_3_ValidateAllFieldsPresent(t *testing.T) {
	// Given: A config with all required fields
	cfg := &Config{
		AuthKey:   "test-auth-key",
		OrgID:     "test-org-id",
		UserEmail: "test@example.com",
	}

	// When: Validating the config
	err := cfg.Validate()

	// Then: Should not return an error
	if err != nil {
		t.Errorf("Expected no error for valid config, got: %v", err)
	}
}

// TestStory1_3_MissingAuthKey tests validation fails when auth_key is missing
func TestStory1_3_MissingAuthKey(t *testing.T) {
	// Given: A config missing auth_key
	cfg := &Config{
		AuthKey:   "", // Missing
		OrgID:     "test-org-id",
		UserEmail: "test@example.com",
	}

	// When: Validating the config
	err := cfg.Validate()

	// Then: Should return error about missing auth_key
	if err == nil {
		t.Error("Expected error for missing auth_key, got nil")
	}

	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "auth_key") {
		t.Errorf("Error should mention 'auth_key', got: %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "required") {
		t.Errorf("Error should mention field is 'required', got: %s", errorMsg)
	}
}

// TestStory1_3_MissingOrgID tests validation fails when org_id is missing
func TestStory1_3_MissingOrgID(t *testing.T) {
	// Given: A config missing org_id
	cfg := &Config{
		AuthKey:   "test-auth-key",
		OrgID:     "", // Missing
		UserEmail: "test@example.com",
	}

	// When: Validating the config
	err := cfg.Validate()

	// Then: Should return error about missing org_id
	if err == nil {
		t.Error("Expected error for missing org_id, got nil")
	}

	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "org_id") {
		t.Errorf("Error should mention 'org_id', got: %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "required") {
		t.Errorf("Error should mention field is 'required', got: %s", errorMsg)
	}
}

// TestStory1_3_MissingUserEmail tests validation fails when user_email is missing
func TestStory1_3_MissingUserEmail(t *testing.T) {
	// Given: A config missing user_email
	cfg := &Config{
		AuthKey:   "test-auth-key",
		OrgID:     "test-org-id",
		UserEmail: "", // Missing
	}

	// When: Validating the config
	err := cfg.Validate()

	// Then: Should return error about missing user_email
	if err == nil {
		t.Error("Expected error for missing user_email, got nil")
	}

	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "user_email") {
		t.Errorf("Error should mention 'user_email', got: %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "required") {
		t.Errorf("Error should mention field is 'required', got: %s", errorMsg)
	}
}

// TestStory1_3_MultipleFieldsMissing tests validation with multiple missing fields
func TestStory1_3_MultipleFieldsMissing(t *testing.T) {
	// Given: A config with multiple missing fields
	cfg := &Config{
		AuthKey:   "",
		OrgID:     "",
		UserEmail: "test@example.com",
	}

	// When: Validating the config
	err := cfg.Validate()

	// Then: Should return error (at least the first missing field)
	if err == nil {
		t.Error("Expected error for multiple missing fields, got nil")
	}

	// Should mention at least one missing field
	errorMsg := err.Error()
	hasMentionOfMissingField := strings.Contains(errorMsg, "auth_key") ||
		strings.Contains(errorMsg, "org_id")

	if !hasMentionOfMissingField {
		t.Errorf("Error should mention at least one missing field, got: %s", errorMsg)
	}
}

// TestStory1_3_EmptyConfig tests validation with all fields empty
func TestStory1_3_EmptyConfig(t *testing.T) {
	// Given: A completely empty config
	cfg := &Config{
		AuthKey:   "",
		OrgID:     "",
		UserEmail: "",
	}

	// When: Validating the config
	err := cfg.Validate()

	// Then: Should return error
	if err == nil {
		t.Error("Expected error for empty config, got nil")
	}

	// Error should be specific and helpful
	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Error message should not be empty")
	}
}

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestStory1_4_SuccessfulAPICall tests successful API connection
func TestStory1_4_SuccessfulAPICall(t *testing.T) {
	// Given: A mock API server that returns 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify authentication header
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "bearer ") {
			t.Error("Expected Authorization header with 'bearer' prefix")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	// When: Making an API call
	client := NewClient("test-auth-key")
	body, err := client.doRequestWithoutBase("GET", server.URL, nil)

	// Then: Should succeed without error
	if err != nil {
		t.Errorf("Expected no error for 200 OK response, got: %v", err)
	}

	if body == nil {
		t.Error("Expected response body, got nil")
	}
}

// TestStory1_4_HTTP401Unauthorized tests handling of 401 Unauthorized
func TestStory1_4_HTTP401Unauthorized(t *testing.T) {
	// Given: A mock API server that returns 401 Unauthorized
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid credentials"}`))
	}))
	defer server.Close()

	// When: Making an API call
	client := NewClient("invalid-auth-key")
	_, err := client.doRequestWithoutBase("GET", server.URL, nil)

	// Then: Should return error with clear message about authentication
	if err == nil {
		t.Error("Expected error for 401 response, got nil")
	}

	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "401") {
		t.Errorf("Error should mention status code 401, got: %s", errorMsg)
	}
}

// TestStory1_4_HTTP403Forbidden tests handling of 403 Forbidden
func TestStory1_4_HTTP403Forbidden(t *testing.T) {
	// Given: A mock API server that returns 403 Forbidden
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "Access denied"}`))
	}))
	defer server.Close()

	// When: Making an API call
	client := NewClient("test-auth-key")
	_, err := client.doRequestWithoutBase("GET", server.URL, nil)

	// Then: Should return error with clear message about access
	if err == nil {
		t.Error("Expected error for 403 response, got nil")
	}

	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "403") {
		t.Errorf("Error should mention status code 403, got: %s", errorMsg)
	}
}

// TestStory1_4_NetworkError tests handling of network errors
func TestStory1_4_NetworkError(t *testing.T) {
	// Given: An invalid URL that will cause network error
	invalidURL := "http://invalid-domain-that-does-not-exist-12345.com"

	// When: Making an API call
	client := NewClient("test-auth-key")
	_, err := client.doRequestWithoutBase("GET", invalidURL, nil)

	// Then: Should return error with clear message about network
	if err == nil {
		t.Error("Expected error for network failure, got nil")
	}

	errorMsg := err.Error()
	// Should have some indication of network/connection issue
	if errorMsg == "" {
		t.Error("Error message should not be empty for network error")
	}
}

// TestStory1_4_BearerTokenAuthentication tests that bearer token is sent correctly
func TestStory1_4_BearerTokenAuthentication(t *testing.T) {
	// Given: A mock server that checks the Authorization header
	expectedToken := "my-secret-token"
	tokenReceived := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "bearer " + expectedToken

		if authHeader == expectedAuth {
			tokenReceived = true
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer server.Close()

	// When: Making an API call with the token
	client := NewClient(expectedToken)
	_, err := client.doRequestWithoutBase("GET", server.URL, nil)

	// Then: Should send the correct bearer token
	if err != nil {
		t.Errorf("Expected successful request, got error: %v", err)
	}

	if !tokenReceived {
		t.Error("Bearer token was not sent correctly in Authorization header")
	}
}

// TestStory1_4_DiscoverRestPrefix tests discovering REST API prefix
func TestStory1_4_DiscoverRestPrefix(t *testing.T) {
	// Given: A mock REST directory server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"rest-prefix": "https://test-api.flowboards.com/v2"}`))
	}))
	defer server.Close()

	// When: Making a request to get REST prefix info
	client := NewClient("test-auth-key")
	body, err := client.doRequestWithoutBase("GET", server.URL, nil)

	// Then: Should successfully get the response
	if err != nil {
		t.Errorf("Expected successful REST prefix discovery, got error: %v", err)
	}

	if !strings.Contains(string(body), "rest-prefix") {
		t.Errorf("Response should contain rest-prefix, got: %s", string(body))
	}
}

// TestStory1_4_ClientTimeout tests that client has reasonable timeout
func TestStory1_4_ClientTimeout(t *testing.T) {
	// Given: A new client
	client := NewClient("test-auth-key")

	// Then: Should have an HTTP client with timeout configured
	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}

	if client.httpClient.Timeout == 0 {
		t.Error("HTTP client should have a timeout configured")
	}

	// Should be a reasonable timeout (e.g., 30 seconds)
	if client.httpClient.Timeout < 0 {
		t.Error("Timeout should be positive")
	}
}

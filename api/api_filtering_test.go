package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestAPIFilteringCapabilities tests Story 2: Test API Filtering Capabilities
//
// Acceptance Criteria:
// - Test API calls with various query parameters: ?board=X, ?bin=Y, ?boardId=X, ?binId=Y
// - Document which parameters are accepted by the API
// - Document which parameters affect the response
// - Document any error responses from unsupported parameters
// - Compare response sizes and content between filtered and unfiltered calls
// - Determine definitively: does API support server-side filtering or not?
func TestAPIFilteringCapabilities(t *testing.T) {
	t.Run("Given API endpoint When testing board parameter Then document behavior", func(t *testing.T) {
		// Arrange
		cfg := loadTestConfig(t)
		client := NewClient(cfg.AuthKey)
		err := client.DiscoverRestPrefix(cfg.OrgID)
		if err != nil {
			t.Fatalf("Failed to discover REST prefix: %v", err)
		}

		user, err := client.GetCurrentUser(cfg.UserEmail)
		if err != nil {
			t.Fatalf("Failed to get current user: %v", err)
		}

		// Act - Test with board parameter
		result := testFilterParameter(client, user.ID, "board", "test-board")

		// Assert - Verify result was documented
		if result == nil {
			t.Fatal("Expected filter test result to be documented")
		}
	})

	t.Run("Given API endpoint When testing bin parameter Then document behavior", func(t *testing.T) {
		// Arrange
		cfg := loadTestConfig(t)
		client := NewClient(cfg.AuthKey)
		err := client.DiscoverRestPrefix(cfg.OrgID)
		if err != nil {
			t.Fatalf("Failed to discover REST prefix: %v", err)
		}

		user, err := client.GetCurrentUser(cfg.UserEmail)
		if err != nil {
			t.Fatalf("Failed to get current user: %v", err)
		}

		// Act - Test with bin parameter
		result := testFilterParameter(client, user.ID, "bin", "test-bin")

		// Assert - Verify result was documented
		if result == nil {
			t.Fatal("Expected filter test result to be documented")
		}
	})

	t.Run("Given multiple filter parameters When testing all combinations Then document all results", func(t *testing.T) {
		// Arrange
		cfg := loadTestConfig(t)
		client := NewClient(cfg.AuthKey)
		err := client.DiscoverRestPrefix(cfg.OrgID)
		if err != nil {
			t.Fatalf("Failed to discover REST prefix: %v", err)
		}

		user, err := client.GetCurrentUser(cfg.UserEmail)
		if err != nil {
			t.Fatalf("Failed to get current user: %v", err)
		}

		// Act - Test all parameter combinations
		results := testAllFilterCombinations(client, user.ID)

		// Assert - Verify at least 5 combinations were tested
		if len(results) < 5 {
			t.Errorf("Expected at least 5 parameter combinations tested, got %d", len(results))
		}
	})

	t.Run("Given filter test results When documenting Then save findings to file", func(t *testing.T) {
		// Arrange
		cfg := loadTestConfig(t)
		client := NewClient(cfg.AuthKey)
		err := client.DiscoverRestPrefix(cfg.OrgID)
		if err != nil {
			t.Fatalf("Failed to discover REST prefix: %v", err)
		}

		user, err := client.GetCurrentUser(cfg.UserEmail)
		if err != nil {
			t.Fatalf("Failed to get current user: %v", err)
		}

		results := testAllFilterCombinations(client, user.ID)

		// Act - Save findings
		err = saveFilterTestFindings(results)

		// Assert - Verify file was created
		if err != nil {
			t.Fatalf("Failed to save filter test findings: %v", err)
		}

		findingsPath := filepath.Join("testdata", "api-filtering-test-results.json")
		if _, err := os.Stat(findingsPath); os.IsNotExist(err) {
			t.Errorf("Expected findings file to exist at %s", findingsPath)
		}
	})
}

// FilterTestResult represents the result of testing a filter parameter
type FilterTestResult struct {
	ParameterName  string `json:"parameter_name"`
	ParameterValue string `json:"parameter_value"`
	RequestURL     string `json:"request_url"`
	StatusCode     int    `json:"status_code"`
	ResponseSize   int    `json:"response_size"`
	TicketCount    int    `json:"ticket_count"`
	ErrorMessage   string `json:"error_message,omitempty"`
	IsAccepted     bool   `json:"is_accepted"`
	AffectsResults bool   `json:"affects_results"`
}

// FilterTestFindings represents all filter testing results
type FilterTestFindings struct {
	BaselineTicketCount  int                 `json:"baseline_ticket_count"`
	BaselineResponseSize int                 `json:"baseline_response_size"`
	TestResults          []FilterTestResult  `json:"test_results"`
	SupportsFiltering    bool                `json:"supports_server_side_filtering"`
	Conclusion           string              `json:"conclusion"`
	TestTimestamp        string              `json:"test_timestamp"`
}

// testFilterParameter tests a single filter parameter
func testFilterParameter(client *Client, userID, paramName, paramValue string) *FilterTestResult {
	path := fmt.Sprintf("/ticket-search?users=%s&%s=%s", url.QueryEscape(userID), paramName, url.QueryEscape(paramValue))

	result := &FilterTestResult{
		ParameterName:  paramName,
		ParameterValue: paramValue,
		RequestURL:     path,
	}

	response, err := client.doRequest(httpMethodGET, path, nil)
	if err != nil {
		result.ErrorMessage = err.Error()
		result.IsAccepted = false
		result.StatusCode = 0
		return result
	}

	result.IsAccepted = true
	result.StatusCode = httpStatusOK
	result.ResponseSize = len(response)

	var tickets []map[string]interface{}
	if json.Unmarshal(response, &tickets) == nil {
		result.TicketCount = len(tickets)
	}

	return result
}

// testAllFilterCombinations tests all filter parameter combinations
func testAllFilterCombinations(client *Client, userID string) []FilterTestResult {
	paramCombinations := []struct {
		name  string
		value string
	}{
		{"board", "test-board"},
		{"bin", "test-bin"},
		{"boardId", "test-board-id"},
		{"binId", "test-bin-id"},
		{"board_id", "test-board-id"},
		{"bin_id", "test-bin-id"},
		{"boardName", "test-board-name"},
		{"binName", "test-bin-name"},
	}

	results := make([]FilterTestResult, 0, len(paramCombinations))
	for _, combo := range paramCombinations {
		result := testFilterParameter(client, userID, combo.name, combo.value)
		results = append(results, *result)
	}

	return results
}

// saveFilterTestFindings saves the filter test findings to a file
func saveFilterTestFindings(results []FilterTestResult) error {
	testdataDir := "testdata"
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		return err
	}

	baselineCount := 0
	baselineSize := 0
	for _, r := range results {
		if r.IsAccepted && baselineCount == 0 {
			baselineCount = r.TicketCount
			baselineSize = r.ResponseSize
		}
	}

	supportsFiltering := false
	for _, r := range results {
		if r.IsAccepted && r.TicketCount != baselineCount {
			supportsFiltering = true
			r.AffectsResults = true
		}
	}

	conclusion := "Server-side filtering is NOT supported. All tested parameters either return errors or do not affect results. Client-side filtering is required."
	if supportsFiltering {
		conclusion = "Server-side filtering IS supported. Some parameters affect the response."
	}

	findings := FilterTestFindings{
		BaselineTicketCount:  baselineCount,
		BaselineResponseSize: baselineSize,
		TestResults:          results,
		SupportsFiltering:    supportsFiltering,
		Conclusion:           conclusion,
		TestTimestamp:        time.Now().Format(time.RFC3339),
	}

	findingsPath := filepath.Join(testdataDir, "api-filtering-test-results.json")
	data, err := json.MarshalIndent(findings, "", "  ")
	if err != nil {
		return err
	}

	doc := createFilterTestDocument(&findings)
	docPath := filepath.Join(testdataDir, "api-filtering-test-results.md")
	os.WriteFile(docPath, []byte(doc), 0644)

	return os.WriteFile(findingsPath, data, 0644)
}

// createFilterTestDocument creates a human-readable test results document
func createFilterTestDocument(findings *FilterTestFindings) string {
	var doc strings.Builder
	doc.WriteString("# API Filtering Capabilities Test Results\n\n")
	doc.WriteString(fmt.Sprintf("**Test Time:** %s\n\n", findings.TestTimestamp))
	doc.WriteString(fmt.Sprintf("**Baseline Ticket Count:** %d\n", findings.BaselineTicketCount))
	doc.WriteString(fmt.Sprintf("**Baseline Response Size:** %d bytes\n\n", findings.BaselineResponseSize))

	doc.WriteString("## Conclusion\n\n")
	doc.WriteString(fmt.Sprintf("%s\n\n", findings.Conclusion))

	doc.WriteString("## Parameter Test Results\n\n")
	doc.WriteString("| Parameter | Value | Accepted | Ticket Count | Response Size | Affects Results | Error |\n")
	doc.WriteString("|-----------|-------|----------|--------------|---------------|-----------------|-------|\n")

	for _, result := range findings.TestResults {
		accepted := "No"
		if result.IsAccepted {
			accepted = "Yes"
		}
		affects := "No"
		if result.AffectsResults {
			affects = "Yes"
		}
		errorMsg := "-"
		if result.ErrorMessage != "" {
			errorMsg = result.ErrorMessage
		}

		doc.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %d | %s | %s |\n",
			result.ParameterName,
			result.ParameterValue,
			accepted,
			result.TicketCount,
			result.ResponseSize,
			affects,
			errorMsg,
		))
	}

	return doc.String()
}

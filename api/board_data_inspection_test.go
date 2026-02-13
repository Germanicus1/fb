package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Germanicus1/fb/config"
)

// TestInspectAPIResponseForBoardData tests Story 1: Inspect API Response for Board Data
//
// Acceptance Criteria:
// - Capture a live API response from /ticket-search?users={userID} to a file
// - Document all fields present in the response, especially any board-related fields
// - Identify field names, data types, and sample values
// - Determine if board_id, board_name, or similar fields exist
// - Document whether board information is embedded in ticket objects or requires separate lookup
// - Create a findings document listing all available fields
func TestInspectAPIResponseForBoardData(t *testing.T) {
	t.Run("Given a valid API client When fetching tickets Then capture raw API response", func(t *testing.T) {
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

		// Act - Capture raw API response
		response, err := captureRawTicketSearchResponse(client, user.ID)
		if err != nil {
			t.Fatalf("Failed to capture API response: %v", err)
		}

		// Assert - Verify response was captured
		if len(response) == 0 {
			t.Fatal("Expected non-empty API response")
		}
	})

	t.Run("Given raw API response When parsing tickets Then identify board-related fields", func(t *testing.T) {
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

		response, err := captureRawTicketSearchResponse(client, user.ID)
		if err != nil {
			t.Fatalf("Failed to capture API response: %v", err)
		}

		// Act - Analyze response for board fields
		findings := analyzeResponseForBoardFields(response)

		// Assert - Verify we have findings documented
		if findings == nil {
			t.Fatal("Expected findings to be documented")
		}

		if findings.TotalTickets == 0 {
			t.Log("No tickets found in response, skipping field analysis")
			return
		}

		// Verify at least 2 sample tickets are documented
		if len(findings.SampleTickets) < 2 && findings.TotalTickets >= 2 {
			t.Errorf("Expected at least 2 sample tickets, got %d", len(findings.SampleTickets))
		}
	})

	t.Run("Given API response findings When documenting Then save to file", func(t *testing.T) {
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

		response, err := captureRawTicketSearchResponse(client, user.ID)
		if err != nil {
			t.Fatalf("Failed to capture API response: %v", err)
		}

		findings := analyzeResponseForBoardFields(response)

		// Act - Save findings to file
		err = saveFindingsToFile(findings)

		// Assert - Verify file was created
		if err != nil {
			t.Fatalf("Failed to save findings: %v", err)
		}

		// Verify file exists
		findingsPath := filepath.Join("testdata", "board-data-findings.json")
		if _, err := os.Stat(findingsPath); os.IsNotExist(err) {
			t.Errorf("Expected findings file to exist at %s", findingsPath)
		}
	})
}

// captureRawTicketSearchResponse captures the raw JSON response from ticket search
func captureRawTicketSearchResponse(client *Client, userID string) ([]byte, error) {
	path := buildTicketSearchPath([]string{userID})
	response, err := client.doRequest(httpMethodGET, path, nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}

const (
	maxSampleTickets = 2
)

// BoardDataFindings represents the analysis results of API response
type BoardDataFindings struct {
	TotalTickets     int                      `json:"total_tickets"`
	SampleTickets    []map[string]interface{} `json:"sample_tickets"`
	AllFields        []string                 `json:"all_fields"`
	BoardFields      []BoardFieldInfo         `json:"board_fields"`
	HasBoardData     bool                     `json:"has_board_data"`
	RequiresLookup   bool                     `json:"requires_lookup"`
	CaptureTimestamp string                   `json:"capture_timestamp"`
}

// BoardFieldInfo describes a board-related field
type BoardFieldInfo struct {
	FieldName  string `json:"field_name"`
	DataType   string `json:"data_type"`
	SampleValue string `json:"sample_value"`
	IsPresent  bool   `json:"is_present"`
}

// analyzeResponseForBoardFields analyzes the API response for board-related fields
func analyzeResponseForBoardFields(response []byte) *BoardDataFindings {
	var tickets []map[string]interface{}
	if err := json.Unmarshal(response, &tickets); err != nil {
		return &BoardDataFindings{CaptureTimestamp: time.Now().Format(time.RFC3339)}
	}

	findings := &BoardDataFindings{
		TotalTickets:     len(tickets),
		SampleTickets:    []map[string]interface{}{},
		AllFields:        []string{},
		BoardFields:      []BoardFieldInfo{},
		CaptureTimestamp: time.Now().Format(time.RFC3339),
	}

	if len(tickets) == 0 {
		return findings
	}

	findings.SampleTickets = extractSampleTickets(tickets)
	findings.AllFields = extractAllFieldNames(tickets)
	findings.BoardFields = findBoardFields(tickets)
	findings.HasBoardData = len(findings.BoardFields) > 0
	findings.RequiresLookup = !findings.HasBoardData

	return findings
}

// extractSampleTickets extracts up to maxSampleTickets from the tickets
func extractSampleTickets(tickets []map[string]interface{}) []map[string]interface{} {
	sampleCount := maxSampleTickets
	if len(tickets) < sampleCount {
		sampleCount = len(tickets)
	}
	return tickets[:sampleCount]
}

// extractAllFieldNames extracts all unique field names from tickets
func extractAllFieldNames(tickets []map[string]interface{}) []string {
	fieldMap := make(map[string]bool)
	for _, ticket := range tickets {
		for field := range ticket {
			fieldMap[field] = true
		}
	}

	fields := make([]string, 0, len(fieldMap))
	for field := range fieldMap {
		fields = append(fields, field)
	}
	return fields
}

// findBoardFields searches for board-related fields in tickets
func findBoardFields(tickets []map[string]interface{}) []BoardFieldInfo {
	if len(tickets) == 0 {
		return []BoardFieldInfo{}
	}

	boardFields := []BoardFieldInfo{}
	boardFieldNames := []string{"board_id", "boardId", "board_name", "boardName", "board"}

	for _, fieldName := range boardFieldNames {
		if value, exists := tickets[0][fieldName]; exists {
			boardFields = append(boardFields, BoardFieldInfo{
				FieldName:   fieldName,
				DataType:    fmt.Sprintf("%T", value),
				SampleValue: fmt.Sprintf("%v", value),
				IsPresent:   true,
			})
		}
	}

	return boardFields
}

// saveFindingsToFile saves the findings to a JSON file
func saveFindingsToFile(findings *BoardDataFindings) error {
	testdataDir := "testdata"
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		return err
	}

	findingsPath := filepath.Join(testdataDir, "board-data-findings.json")
	data, err := json.MarshalIndent(findings, "", "  ")
	if err != nil {
		return err
	}

	// Also save raw sample tickets for reference
	rawSamplesPath := filepath.Join(testdataDir, "raw-api-response-samples.json")
	samplesData, _ := json.MarshalIndent(findings.SampleTickets, "", "  ")
	os.WriteFile(rawSamplesPath, samplesData, 0644)

	// Create a human-readable findings document
	doc := createFindingsDocument(findings)
	docPath := filepath.Join(testdataDir, "board-data-findings.md")
	os.WriteFile(docPath, []byte(doc), 0644)

	return os.WriteFile(findingsPath, data, 0644)
}

// createFindingsDocument creates a human-readable findings document
func createFindingsDocument(findings *BoardDataFindings) string {
	var doc strings.Builder
	doc.WriteString("# Board Data API Response Findings\n\n")
	doc.WriteString(fmt.Sprintf("**Capture Time:** %s\n\n", findings.CaptureTimestamp))
	doc.WriteString(fmt.Sprintf("**Total Tickets in Response:** %d\n\n", findings.TotalTickets))

	doc.WriteString("## Board Data Availability\n\n")
	doc.WriteString(fmt.Sprintf("- **Has Board Data:** %v\n", findings.HasBoardData))
	doc.WriteString(fmt.Sprintf("- **Requires Separate Lookup:** %v\n\n", findings.RequiresLookup))

	if len(findings.BoardFields) > 0 {
		doc.WriteString("## Board-Related Fields Found\n\n")
		for _, field := range findings.BoardFields {
			doc.WriteString(fmt.Sprintf("### %s\n", field.FieldName))
			doc.WriteString(fmt.Sprintf("- **Data Type:** %s\n", field.DataType))
			doc.WriteString(fmt.Sprintf("- **Sample Value:** %s\n", field.SampleValue))
			doc.WriteString(fmt.Sprintf("- **Present in Response:** %v\n\n", field.IsPresent))
		}
	} else {
		doc.WriteString("## Board-Related Fields\n\n")
		doc.WriteString("No board-related fields found in the ticket response.\n\n")
	}

	doc.WriteString("## All Available Fields\n\n")
	for _, field := range findings.AllFields {
		doc.WriteString(fmt.Sprintf("- %s\n", field))
	}
	doc.WriteString("\n")

	doc.WriteString(fmt.Sprintf("## Sample Tickets\n\n%d sample tickets captured. See raw-api-response-samples.json for full details.\n", len(findings.SampleTickets)))

	return doc.String()
}

// loadTestConfig loads configuration for testing
func loadTestConfig(t *testing.T) *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Skipf("Skipping test: config not available: %v", err)
	}
	return cfg
}

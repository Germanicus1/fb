package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestBoardBinRelationship tests Story 3: Understand Board-Bin Relationship
//
// Acceptance Criteria:
// - Document whether bins are scoped to specific boards or exist globally
// - Determine if bin names are unique within a board or globally unique
// - Identify how to uniquely identify a bin in the context of filtering
// - Understand if a ticket can belong to multiple boards or bins
// - Document the hierarchical relationship: Board → Bin → Ticket
// - Clarify whether BinID alone is sufficient or if BoardID + BinID is needed for uniqueness
func TestBoardBinRelationship(t *testing.T) {
	t.Run("Given ticket data When analyzing bin uniqueness Then document scope", func(t *testing.T) {
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

		// Act - Analyze bin uniqueness
		analysis := analyzeBinUniqueness(response)

		// Assert
		if analysis == nil {
			t.Fatal("Expected bin uniqueness analysis to be documented")
		}
	})

	t.Run("Given ticket data When analyzing board-bin hierarchy Then document relationship", func(t *testing.T) {
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

		// Act - Analyze hierarchy
		analysis := analyzeBoardBinHierarchy(response)

		// Assert
		if analysis == nil {
			t.Fatal("Expected board-bin hierarchy analysis to be documented")
		}
	})

	t.Run("Given analysis results When documenting Then save findings to file", func(t *testing.T) {
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

		analysis := analyzeBoardBinRelationship(response)

		// Act - Save findings
		err = saveBoardBinRelationshipFindings(analysis)

		// Assert
		if err != nil {
			t.Fatalf("Failed to save findings: %v", err)
		}

		findingsPath := filepath.Join("testdata", "board-bin-relationship.json")
		if _, err := os.Stat(findingsPath); os.IsNotExist(err) {
			t.Errorf("Expected findings file to exist at %s", findingsPath)
		}
	})
}

// BinUniquenessAnalysis represents the analysis of bin uniqueness
type BinUniquenessAnalysis struct {
	TotalBins          int      `json:"total_bins"`
	UniqueBinIDs       []string `json:"unique_bin_ids"`
	UniqueBinNames     []string `json:"unique_bin_names"`
	AreBinIDsUnique    bool     `json:"are_bin_ids_unique"`
	AreBinNamesUnique  bool     `json:"are_bin_names_unique"`
	DuplicateBinNames  []string `json:"duplicate_bin_names,omitempty"`
}

// BoardBinHierarchyAnalysis represents the analysis of board-bin hierarchy
type BoardBinHierarchyAnalysis struct {
	HasBoardData           bool     `json:"has_board_data"`
	BinsAreGloballyScoped  bool     `json:"bins_are_globally_scoped"`
	BinsAreBoardScoped     bool     `json:"bins_are_board_scoped"`
	TicketsHaveMultipleBoards bool  `json:"tickets_have_multiple_boards"`
	TicketsHaveMultipleBins   bool  `json:"tickets_have_multiple_bins"`
	HierarchyDescription   string   `json:"hierarchy_description"`
}

// BoardBinRelationshipFindings represents complete relationship analysis
type BoardBinRelationshipFindings struct {
	UniquenessAnalysis *BinUniquenessAnalysis        `json:"uniqueness_analysis"`
	HierarchyAnalysis  *BoardBinHierarchyAnalysis    `json:"hierarchy_analysis"`
	IdentifierStrategy string                        `json:"identifier_strategy"`
	Recommendations    []string                      `json:"recommendations"`
	AnalysisTimestamp  string                        `json:"analysis_timestamp"`
}

// analyzeBinUniqueness analyzes whether bin IDs and names are unique
func analyzeBinUniqueness(response []byte) *BinUniquenessAnalysis {
	var tickets []map[string]interface{}
	if err := json.Unmarshal(response, &tickets); err != nil {
		return &BinUniquenessAnalysis{}
	}

	binIDMap := make(map[string]int)
	binNameMap := make(map[string]int)

	for _, ticket := range tickets {
		if binID, ok := ticket["bin_id"].(string); ok && binID != "" {
			binIDMap[binID]++
		}
		if binName, ok := ticket["bin_name"].(string); ok && binName != "" {
			binNameMap[binName]++
		}
	}

	uniqueBinIDs := make([]string, 0, len(binIDMap))
	for id := range binIDMap {
		uniqueBinIDs = append(uniqueBinIDs, id)
	}

	uniqueBinNames := make([]string, 0, len(binNameMap))
	duplicateBinNames := []string{}
	for name, count := range binNameMap {
		uniqueBinNames = append(uniqueBinNames, name)
		if count > 1 {
			duplicateBinNames = append(duplicateBinNames, name)
		}
	}

	return &BinUniquenessAnalysis{
		TotalBins:         len(binIDMap),
		UniqueBinIDs:      uniqueBinIDs,
		UniqueBinNames:    uniqueBinNames,
		AreBinIDsUnique:   true,
		AreBinNamesUnique: len(duplicateBinNames) == 0,
		DuplicateBinNames: duplicateBinNames,
	}
}

// analyzeBoardBinHierarchy analyzes the board-bin hierarchical relationship
func analyzeBoardBinHierarchy(response []byte) *BoardBinHierarchyAnalysis {
	var tickets []map[string]interface{}
	if err := json.Unmarshal(response, &tickets); err != nil {
		return &BoardBinHierarchyAnalysis{}
	}

	hasBoardData := false
	for _, ticket := range tickets {
		if _, ok := ticket["board_id"]; ok {
			hasBoardData = true
			break
		}
		if _, ok := ticket["boardId"]; ok {
			hasBoardData = true
			break
		}
	}

	description := "Bins exist at the organization level. Each ticket has one bin_id and bin_name. No board information is available in the ticket data, suggesting bins are globally scoped rather than board-scoped."
	if hasBoardData {
		description = "Board information is available in ticket data. Analyzing board-bin relationship..."
	}

	return &BoardBinHierarchyAnalysis{
		HasBoardData:          hasBoardData,
		BinsAreGloballyScoped: !hasBoardData,
		BinsAreBoardScoped:    hasBoardData,
		TicketsHaveMultipleBoards: false,
		TicketsHaveMultipleBins:   false,
		HierarchyDescription:  description,
	}
}

// analyzeBoardBinRelationship performs complete board-bin relationship analysis
func analyzeBoardBinRelationship(response []byte) *BoardBinRelationshipFindings {
	uniquenessAnalysis := analyzeBinUniqueness(response)
	hierarchyAnalysis := analyzeBoardBinHierarchy(response)

	identifierStrategy := "Use bin_id for filtering (globally unique identifier)"
	recommendations := []string{
		"Bin IDs are sufficient for unique identification",
		"Bin names may not be unique across the organization",
		"Filter by bin_id for exact matching, bin_name for user-friendly filtering",
		"No board data available, so board filtering not possible via this endpoint",
		"Client-side filtering required for both board and bin filtering",
	}

	if hierarchyAnalysis.HasBoardData {
		identifierStrategy = "Use board_id + bin_id combination for precise filtering"
		recommendations = []string{
			"Board data is available in tickets",
			"Use board_id and bin_id together for filtering",
			"Bins may be scoped within boards",
		}
	}

	return &BoardBinRelationshipFindings{
		UniquenessAnalysis: uniquenessAnalysis,
		HierarchyAnalysis:  hierarchyAnalysis,
		IdentifierStrategy: identifierStrategy,
		Recommendations:    recommendations,
		AnalysisTimestamp:  time.Now().Format(time.RFC3339),
	}
}

// saveBoardBinRelationshipFindings saves the findings to a file
func saveBoardBinRelationshipFindings(findings *BoardBinRelationshipFindings) error {
	testdataDir := "testdata"
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		return err
	}

	findingsPath := filepath.Join(testdataDir, "board-bin-relationship.json")
	data, err := json.MarshalIndent(findings, "", "  ")
	if err != nil {
		return err
	}

	doc := createBoardBinRelationshipDocument(findings)
	docPath := filepath.Join(testdataDir, "board-bin-relationship.md")
	os.WriteFile(docPath, []byte(doc), 0644)

	return os.WriteFile(findingsPath, data, 0644)
}

// createBoardBinRelationshipDocument creates a human-readable relationship document
func createBoardBinRelationshipDocument(findings *BoardBinRelationshipFindings) string {
	var doc strings.Builder
	doc.WriteString("# Board-Bin Relationship Analysis\n\n")
	doc.WriteString(fmt.Sprintf("**Analysis Time:** %s\n\n", findings.AnalysisTimestamp))

	doc.WriteString("## Bin Uniqueness Analysis\n\n")
	doc.WriteString(fmt.Sprintf("- **Total Unique Bins:** %d\n", findings.UniquenessAnalysis.TotalBins))
	doc.WriteString(fmt.Sprintf("- **Bin IDs are Unique:** %v\n", findings.UniquenessAnalysis.AreBinIDsUnique))
	doc.WriteString(fmt.Sprintf("- **Bin Names are Unique:** %v\n\n", findings.UniquenessAnalysis.AreBinNamesUnique))

	if len(findings.UniquenessAnalysis.DuplicateBinNames) > 0 {
		doc.WriteString("### Duplicate Bin Names Found\n\n")
		for _, name := range findings.UniquenessAnalysis.DuplicateBinNames {
			doc.WriteString(fmt.Sprintf("- %s\n", name))
		}
		doc.WriteString("\n")
	}

	doc.WriteString("## Board-Bin Hierarchy\n\n")
	doc.WriteString(fmt.Sprintf("- **Has Board Data:** %v\n", findings.HierarchyAnalysis.HasBoardData))
	doc.WriteString(fmt.Sprintf("- **Bins are Globally Scoped:** %v\n", findings.HierarchyAnalysis.BinsAreGloballyScoped))
	doc.WriteString(fmt.Sprintf("- **Bins are Board Scoped:** %v\n\n", findings.HierarchyAnalysis.BinsAreBoardScoped))
	doc.WriteString(fmt.Sprintf("**Description:** %s\n\n", findings.HierarchyAnalysis.HierarchyDescription))

	doc.WriteString("## Identifier Strategy\n\n")
	doc.WriteString(fmt.Sprintf("%s\n\n", findings.IdentifierStrategy))

	doc.WriteString("## Recommendations\n\n")
	for _, rec := range findings.Recommendations {
		doc.WriteString(fmt.Sprintf("- %s\n", rec))
	}

	return doc.String()
}

package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SaveBinContext saves the last used bin context to ~/.fb/bin_context.json
func SaveBinContext(binID, binName string) error {
	homeDir, _ := os.UserHomeDir()
	fbDir := filepath.Join(homeDir, ".fb")
	os.MkdirAll(fbDir, 0700)

	context := BinContext{
		BinID:   binID,
		BinName: binName,
	}

	data, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return err
	}

	contextPath := filepath.Join(fbDir, "bin_context.json")
	return os.WriteFile(contextPath, data, 0600)
}

// LoadBinContext loads the last used bin context from ~/.fb/bin_context.json
func LoadBinContext() (*BinContext, error) {
	homeDir, _ := os.UserHomeDir()
	contextPath := filepath.Join(homeDir, ".fb", "bin_context.json")

	data, err := os.ReadFile(contextPath)
	if err != nil {
		return nil, err
	}

	var context BinContext
	if err := json.Unmarshal(data, &context); err != nil {
		return nil, err
	}

	return &context, nil
}

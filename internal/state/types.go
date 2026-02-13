// Package state manages persistent application state stored in ~/.fb/.
// It handles checkout state and bin context using JSON file storage.
package state

// CheckoutState represents the persisted checkout state
type CheckoutState struct {
	TicketID     string `json:"ticket_id"`
	TicketName   string `json:"ticket_name"`
	BinID        string `json:"bin_id"`
	BinName      string `json:"bin_name"`
	CheckedOutAt string `json:"checked_out_at"`
}

// BinContext represents the last used bin
type BinContext struct {
	BinID   string `json:"bin_id"`
	BinName string `json:"bin_name"`
}

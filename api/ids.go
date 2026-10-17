package api

import (
	"crypto/rand"
	"math/big"
)

const (
	idLength   = 17
	idAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// GenerateID returns a client-side entity ID: 17 random [A-Za-z0-9] characters.
// Used for tickets and comments.
func GenerateID() string {
	out := make([]byte, idLength)
	max := big.NewInt(int64(len(idAlphabet)))
	for i := range out {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic("crypto/rand unavailable: " + err.Error())
		}
		out[i] = idAlphabet[n.Int64()]
	}
	return string(out)
}

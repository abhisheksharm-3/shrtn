package util

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

// HashString creates a hash of the input string
func HashString(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Trim padding characters and take first 8 characters
	hash = strings.TrimRight(hash, "=")
	if len(hash) > 8 {
		hash = hash[:8]
	}

	return hash
}

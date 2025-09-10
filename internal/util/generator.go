package util

import (
	"math/rand"
	"time"
)

const (
	// The set of characters to use for generating the short code.
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// The length of the generated short code.
	codeLength = 8
)

// seededRand is a random number generator seeded with the current time.
// We create it as a package-level variable to ensure it's seeded only once.
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateShortCode creates a random string of a fixed length.
func GenerateShortCode() string {
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

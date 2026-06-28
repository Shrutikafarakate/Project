package utils

import (
	"math/rand"
	"strings"
	"time"
)

// GenerateShortCode generates a random alphanumeric string for the shortcode.
func GenerateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 2
	rand.Seed(time.Now().UnixNano())

	var shortCode strings.Builder
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		shortCode.WriteByte(charset[randomIndex])
	}
	return shortCode.String()
}

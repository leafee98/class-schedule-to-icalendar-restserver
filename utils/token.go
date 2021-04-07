package utils

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateToken will return random uuid removed dashes
func GenerateToken() string {
	u, _ := uuid.NewRandom()
	var uuid string = u.String()

	// remove dashes in token
	var builder strings.Builder
	builder.Grow(len(uuid))
	for _, c := range uuid {
		if c != '-' {
			builder.WriteRune(c)
		}
	}
	return builder.String()
}

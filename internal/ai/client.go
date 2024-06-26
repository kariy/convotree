package ai

import (
	"strings"

	"convotree/internal/core"

	"golang.org/x/exp/rand"
)

// TODO: This should be replaced with a real AI service
func GetAIResponse(_ []core.Exchange) (string, error) {
	
	// Generate a random string of up to 10 words
	words := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew", "kiwi", "lemon"}
	n := rand.Intn(10) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString(words[rand.Intn(len(words))])
		sb.WriteString(" ")
	}
	return strings.TrimSpace(sb.String()), nil
}
package ui

import (
	"strings"
	"testing"
)

func TestSlugFromPrompt(t *testing.T) {
	tests := []struct {
		name   string
		prompt string
		want   string
	}{
		{"simple", "add a login form", "add-a-login-form"},
		{"strips punctuation and apostrophes", "Let's add a table for 123", "lets-add-a-table-for-123"},
		{"lowercases", "Fix The Bug", "fix-the-bug"},
		{"caps leading words", "one two three four five six seven eight", "one-two-three-four-five-six"},
		{"collapses extra whitespace", "  refactor   the\tparser  ", "refactor-the-parser"},
		{"drops symbol-only words", "update @@@ handler", "update-handler"},
		{"empty prompt", "", "session"},
		{"symbols only", "!!! @#$ ***", "session"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slugFromPrompt(tt.prompt); got != tt.want {
				t.Errorf("slugFromPrompt(%q) = %q, want %q", tt.prompt, got, tt.want)
			}
		})
	}
}

func TestSlugFromPromptRespectsMaxLength(t *testing.T) {
	long := strings.Repeat("abcdefghij", 10) // 100 chars, one word
	got := slugFromPrompt(long)
	if len(got) > MaxNameLength {
		t.Errorf("slug length %d exceeds MaxNameLength %d: %q", len(got), MaxNameLength, got)
	}
	if got == "" {
		t.Error("slug should not be empty for a long word")
	}
}

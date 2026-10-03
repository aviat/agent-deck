package ui

import "strings"

// indentLines prefixes every line of s with prefix. Used to indent the
// multi-line prompt textarea so each rendered row aligns under its field label.
func indentLines(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// promptSlugMaxWords bounds how many of the prompt's leading words feed the
// auto-derived title, so a long prompt still yields a short, glanceable name.
const promptSlugMaxWords = 6

// slugFromPrompt derives a short kebab-case session title from the initial
// prompt's first words: lowercased, non-alphanumeric characters dropped, words
// joined with hyphens (e.g. "Let's add a table for 123" -> "lets-add-a-table-for-123").
//
// It is the offline fallback used when the user launches a session from the
// prompt field without typing a Name. The result is length-capped to
// MaxNameLength and never empty — a prompt with no usable characters yields
// "session". The title is left unlocked by the caller so the Claude
// session-name sync can later replace this slug with a real, conversation-derived
// title.
func slugFromPrompt(prompt string) string {
	var words []string
	for _, field := range strings.Fields(prompt) {
		var b strings.Builder
		for _, r := range strings.ToLower(field) {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				b.WriteRune(r)
			}
		}
		if b.Len() == 0 {
			continue
		}
		words = append(words, b.String())
		if len(words) >= promptSlugMaxWords {
			break
		}
	}

	slug := strings.Join(words, "-")
	if len(slug) > MaxNameLength {
		slug = strings.Trim(slug[:MaxNameLength], "-")
	}
	if slug == "" {
		return "session"
	}
	return slug
}

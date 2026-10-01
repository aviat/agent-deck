package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func typePrompt(d *NewDialog, s string) *NewDialog {
	for _, r := range s {
		d, _ = d.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	return d
}

// TestNewDialog_PromptFirstNameAutocompletes types into the prompt and checks the
// Name field live-fills with the slug — and that editing Name stops the autofill.
func TestNewDialog_PromptFirstNameAutocompletes(t *testing.T) {
	d := NewNewDialog()
	d.promptFirst = true
	d.ShowInGroup("default", "default", "", nil, "")
	d.rebuildFocusTargets()

	d = typePrompt(d, "fix bug")
	if got := d.nameInput.Value(); got != "fix-bug" {
		t.Fatalf("autocompleted name = %q, want %q", got, "fix-bug")
	}
	if !d.IsAutoName() {
		t.Fatal("IsAutoName should be true while Name is autofilled")
	}

	// Move focus to Name and type: autofill must stop and the user's edit stays.
	d.focusIndex = d.indexOf(focusName)
	d.updateFocus()
	d = typePrompt(d, "x")
	if d.IsAutoName() {
		t.Fatal("IsAutoName should be false after the user edits Name")
	}
	nameAfterEdit := d.nameInput.Value()

	// Further prompt edits must no longer touch the Name.
	d.focusIndex = d.indexOf(focusPrompt)
	d.updateFocus()
	d = typePrompt(d, "ging")
	if got := d.nameInput.Value(); got != nameAfterEdit {
		t.Fatalf("Name changed after manual edit: got %q, want %q", got, nameAfterEdit)
	}
}

// TestNewDialog_DefaultLayoutUnchanged guards the toggle-off default: the dialog
// leads with Name (not the prompt) and never exposes the prompt focus target, so
// the existing flow is untouched when [ui] prompt_first is not set.
func TestNewDialog_DefaultLayoutUnchanged(t *testing.T) {
	d := NewNewDialog()
	if d.promptFirst {
		t.Fatal("promptFirst should default to false")
	}
	d.ShowInGroup("default", "default", "", nil, "")
	if got := d.currentTarget(); got != focusName {
		t.Fatalf("default focus = %v, want focusName", got)
	}
	for _, target := range d.focusTargets {
		if target == focusPrompt {
			t.Fatal("focusPrompt must not appear in the default layout")
		}
	}
}

// TestNewDialog_PromptFirstLayout covers the opt-in layout: prompt is focused
// first, Name is optional (validates with a prompt alone), and the title/branch
// derive from the prompt.
func TestNewDialog_PromptFirstLayout(t *testing.T) {
	d := NewNewDialog()
	d.promptFirst = true
	d.ShowInGroup("default", "default", "", nil, "")
	d.rebuildFocusTargets()

	if len(d.focusTargets) == 0 || d.focusTargets[0] != focusPrompt {
		t.Fatalf("focusTargets[0] = %v, want focusPrompt", d.focusTargets)
	}
	if got := d.currentTarget(); got != focusPrompt {
		t.Fatalf("default focus = %v, want focusPrompt", got)
	}
	// Enter on the prompt row submits (handled by home.go), not locally.
	if d.shouldHandleEnterLocally() {
		t.Fatal("Enter on focusPrompt should submit, not be handled locally")
	}

	// Empty prompt + empty name is refused.
	if msg := d.Validate(); msg == "" {
		t.Fatal("Validate should refuse an empty prompt and empty name")
	}

	// A prompt with no name validates, and the effective name is the slug.
	d.pathInput.SetValue("/tmp")
	d.promptInput.SetValue("add a users table")
	if msg := d.Validate(); msg != "" {
		t.Fatalf("Validate with prompt and no name = %q, want valid", msg)
	}
	if got := d.EffectiveName(); got != "add-a-users-table" {
		t.Fatalf("EffectiveName() = %q, want %q", got, "add-a-users-table")
	}
}

// TestNewDialog_PromptFirstRemoteCarriesPrompt checks that a prompt-first remote
// create delivers the prompt as the claude startup query and titles the session
// from the prompt slug (regression for the silent-drop finding).
func TestNewDialog_PromptFirstRemoteCarriesPrompt(t *testing.T) {
	d := NewNewDialog()
	d.promptFirst = true
	d.ShowInGroup("default", "default", "", nil, "")
	d.rebuildFocusTargets()
	// Force a claude selection so the startup-query path is exercised.
	selected := false
	for i, cmd := range d.presetCommands {
		if d.toolKind(cmd) == "claude" {
			d.commandCursor = i
			selected = true
			break
		}
	}
	if !selected {
		t.Skip("no claude preset available in this environment")
	}
	d.pathInput.SetValue("/tmp")
	d.promptInput.SetValue("add a users table")

	opts, errMsg := d.GetRemoteCreateOptions()
	if errMsg != "" {
		t.Fatalf("GetRemoteCreateOptions error = %q, want none", errMsg)
	}
	if opts.StartQuery != "add a users table" {
		t.Fatalf("StartQuery = %q, want the prompt text", opts.StartQuery)
	}
	if opts.Title != "add-a-users-table" {
		t.Fatalf("Title = %q, want the prompt slug", opts.Title)
	}
}

// TestNewDialog_PromptFirstBranchFromPrompt reproduces the "branch name required
// for worktree" bug: with a worktree enabled and no Name typed, the branch must
// be derived from the prompt rather than left empty.
func TestNewDialog_PromptFirstBranchFromPrompt(t *testing.T) {
	d := NewNewDialog()
	d.promptFirst = true
	d.ShowInGroup("default", "default", "", nil, "")
	d.rebuildFocusTargets()

	d.worktreeEnabled = true
	d.branchAutoSet = true
	d.branchPrefix = "feature/"
	d.promptInput.SetValue("fix the parser")
	d.autoBranchFromName()

	if got := d.branchInput.Value(); got != "feature/fix-the-parser" {
		t.Fatalf("branch = %q, want %q", got, "feature/fix-the-parser")
	}
}

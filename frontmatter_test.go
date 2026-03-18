package main

import (
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	input := []byte("---\ntitle: Hello\ndescription: World\n---\nBody here")
	fm, body := ParseFrontmatter(input)

	if fm.Title != "Hello" {
		t.Errorf("Title = %q, want Hello", fm.Title)
	}
	if fm.Description != "World" {
		t.Errorf("Description = %q, want World", fm.Description)
	}
	if string(body) != "Body here" {
		t.Errorf("Body = %q, want 'Body here'", string(body))
	}
}

func TestParseFrontmatterNoDelimiter(t *testing.T) {
	input := []byte("Just markdown, no frontmatter")
	fm, body := ParseFrontmatter(input)

	if fm.Title != "" {
		t.Errorf("Title = %q, want empty", fm.Title)
	}
	if string(body) != string(input) {
		t.Errorf("Body should be full input")
	}
}

func TestParseFrontmatterExtras(t *testing.T) {
	input := []byte("---\ntitle: Test\ncustom_key: custom_val\n---\nBody")
	fm, _ := ParseFrontmatter(input)

	if fm.Extra == nil {
		t.Fatal("Extra should not be nil")
	}
	if fm.Extra["custom_key"] != "custom_val" {
		t.Errorf("Extra[custom_key] = %v, want custom_val", fm.Extra["custom_key"])
	}
}

func TestParseFrontmatterDateAndDraft(t *testing.T) {
	input := []byte("---\ntitle: My Post\ndate: 2026-03-18\ndraft: true\n---\nBody")
	fm, _ := ParseFrontmatter(input)

	if fm.Date != "2026-03-18" {
		t.Errorf("Date = %q, want 2026-03-18", fm.Date)
	}
	if !fm.Draft {
		t.Error("Draft = false, want true")
	}
	// date and draft should NOT appear in Extra
	if fm.Extra != nil {
		if _, ok := fm.Extra["date"]; ok {
			t.Error("date should not be in Extra")
		}
		if _, ok := fm.Extra["draft"]; ok {
			t.Error("draft should not be in Extra")
		}
	}
}

func TestParseDateFormats(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
	}{
		{"2026-03-18", true},
		{"2026-03-18 14:30", true},
		{"2026-03-18 14:30:00", true},
		{"2026-03-18T14:30", true},
		{"2026-03-18T14:30:00", true},
		{"18-03-2026", false},
		{"March 18, 2026", false},
		{"", false},
	}

	for _, tt := range tests {
		_, ok := ParseDate(tt.input)
		if ok != tt.ok {
			t.Errorf("ParseDate(%q) ok=%v, want %v", tt.input, ok, tt.ok)
		}
	}
}

func TestParseDateTimestampOrdering(t *testing.T) {
	// String comparison of ISO timestamps must sort correctly
	a := "2026-02-25 10:00"
	b := "2026-02-25 18:00"
	if a >= b {
		t.Errorf("expected %q < %q", a, b)
	}

	// Date-only vs timestamp
	c := "2026-02-25"
	d := "2026-03-18"
	if c >= d {
		t.Errorf("expected %q < %q", c, d)
	}
}

func TestParseFrontmatterCRLF(t *testing.T) {
	input := []byte("---\r\ntitle: CRLF\r\n---\r\nBody")
	fm, body := ParseFrontmatter(input)

	if fm.Title != "CRLF" {
		t.Errorf("Title = %q, want CRLF", fm.Title)
	}
	if string(body) != "Body" {
		t.Errorf("Body = %q, want 'Body'", string(body))
	}
}

func TestParseFrontmatterUnclosed(t *testing.T) {
	input := []byte("---\ntitle: Unclosed\nno closing delimiter")
	fm, body := ParseFrontmatter(input)

	if fm.Title != "" {
		t.Errorf("Title = %q, want empty (unclosed)", fm.Title)
	}
	if string(body) != string(input) {
		t.Errorf("Body should be full input when unclosed")
	}
}

func TestTitleFromFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"02-agents.md", "Agents"},
		{"quickstart.md", "Quickstart"},
		{"1-intro.md", "Intro"},
		{"001-advanced-topics.md", "Advanced Topics"},
		{"no-prefix.md", "No Prefix"},
		{"index.md", "Index"},
	}

	for _, tt := range tests {
		got := TitleFromFilename(tt.input)
		if got != tt.want {
			t.Errorf("TitleFromFilename(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestTitleFromDir(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"01-guide", "Guide"},
		{"1-intro", "Intro"},
		{"reference", "Reference"},
		{"001-advanced-topics", "Advanced Topics"},
	}

	for _, tt := range tests {
		got := TitleFromDir(tt.input)
		if got != tt.want {
			t.Errorf("TitleFromDir(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

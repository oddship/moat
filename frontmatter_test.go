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

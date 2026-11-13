package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportMarkdownPreservesExamples(t *testing.T) {
	resolver := newPageResolver([]Page{{RelPath: "01-guide/01-start.md", Frontmatter: Frontmatter{Title: "Start"}}}, "/project")
	body := []byte("# Guide\n\nSee [[Start#install|install now]] and [[Missing]].\n\n`[[Start]]`\n\n```md\n[[Start]]\n```\n\n    [[Start]]\n\n\\[[Start]]\n\n| Option | Meaning |\n| --- | --- |\n| A | B |\n\n<div>custom HTML</div>\n")
	got := string(exportMarkdown(body, resolver))
	for _, want := range []string{
		"[install now](</project/guide/start/#install>) and Missing.",
		"`[[Start]]`", "```md\n[[Start]]\n```", "    [[Start]]",
		"\\[[Start]]", "| A | B |", "<div>custom HTML</div>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("export missing %q:\n%s", want, got)
		}
	}
}

func TestExportMarkdownPreservesLeadingIndentation(t *testing.T) {
	body := "    [[Start]]\n"
	got := string(exportMarkdown([]byte(body), newPageResolver(nil, "")))
	if got != body {
		t.Fatalf("indented code changed: %q", got)
	}
}

func TestExportMarkdownNestedContexts(t *testing.T) {
	resolver := newPageResolver([]Page{{RelPath: "start.md", Frontmatter: Frontmatter{Title: "Start"}}}, "")
	for _, body := range []string{
		"> See [[Start]]\n", "- See [[Start]] and [[Start|again]]\n",
		"| Page |\n| --- |\n| [[Start]] |\n", "See [[Start]].\n\nSee [[Start]] again.\n",
	} {
		got := string(exportMarkdown([]byte(body), resolver))
		if strings.Contains(got, "[[Start") || !strings.Contains(got, "[Start](</start/>)") {
			t.Errorf("unresolved link in %q: %q", body, got)
		}
	}
}

func TestAgentExportsStaticHosting(t *testing.T) {
	for _, base := range []string{"", "/project"} {
		t.Run("base="+base, func(t *testing.T) {
			src, dst := t.TempDir(), t.TempDir()
			files := map[string]string{
				"index.md":              "---\ntitle: Home\n---\n# Home\n\n[[Start]]\n",
				"01-guide/01-start.md":  "---\ntitle: Start\ndescription: Install the project\nurl: custom/start\n---\n# Start\n\n{{< note >}}Keep this content.{{< /note >}}\n\n" + strings.Repeat("long content ", 250),
				"01-guide/index.md":     "# Guide\n",
				"_shortcodes/note.html": "<aside>{{ .Inner }}</aside>",
				"draft.md":              "---\ndraft: true\n---\nSECRET DRAFT",
				".hidden.md":            "SECRET HIDDEN",
				"_private.md":           "SECRET PRIVATE",
			}
			for path, body := range files {
				fullPath := filepath.Join(src, path)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(fullPath, []byte(body), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			cfg := Config{SiteName: "Project", BasePath: base, Agents: AgentsConfig{Full: true}}
			if err := Build(src, dst, cfg); err != nil {
				t.Fatal(err)
			}
			handler := http.FileServer(http.Dir(dst))
			if base != "" {
				handler = http.StripPrefix(base, handler)
			}
			server := httptest.NewServer(handler)
			defer server.Close()
			read := func(path string) string {
				t.Helper()
				response, err := http.Get(server.URL + base + path)
				if err != nil {
					t.Fatal(err)
				}
				defer response.Body.Close()
				if response.StatusCode != http.StatusOK {
					t.Fatalf("GET %s: %s", path, response.Status)
				}
				body, err := io.ReadAll(response.Body)
				if err != nil {
					t.Fatal(err)
				}
				return string(body)
			}
			index := read("/llms.txt")
			for _, path := range []string{"/index.md", "/guide/index.md", "/custom/start/index.md"} {
				if !strings.Contains(index, "(<"+base+path+">)") {
					t.Errorf("catalog missing %s:\n%s", path, index)
				}
				read(path)
			}
			export := read("/custom/start/index.md")
			if strings.Contains(export, "title:") || strings.Contains(export, "{{<") || !strings.Contains(export, "<aside><p>Keep this content.</p>") || !strings.Contains(export, strings.Repeat("long content ", 240)) {
				t.Fatalf("export lost content or leaked source metadata: %s", export)
			}
			full := read("/llms-full.txt")
			if !strings.Contains(full, export) || !strings.Contains(full, "Source: <"+base+"/custom/start/index.md>") {
				t.Fatal("full bundle missing complete page or URL")
			}
			for _, body := range []string{index, full} {
				if strings.Contains(body, "SECRET") || strings.Contains(body, "draft") || strings.Contains(body, "hidden") || strings.Contains(body, "private") {
					t.Fatal("excluded content leaked into exports")
				}
			}
			html := read("/custom/start/")
			for _, want := range []string{
				`type="text/markdown" href="` + base + `/custom/start/index.md"`,
				`rel="describedby" href="` + base + `/llms.txt"`,
				`data-markdown-url="` + base + `/custom/start/index.md"`,
			} {
				if !strings.Contains(html, want) {
					t.Errorf("HTML missing %s", want)
				}
			}
			// Rebuilding with exports disabled must remove all previous artifacts.
			disabled := false
			cfg.Agents.Enabled = &disabled
			if err := Build(src, dst, cfg); err != nil {
				t.Fatal(err)
			}
			for _, path := range []string{"llms.txt", "llms-full.txt", "index.md", "custom/start/index.md"} {
				if _, err := os.Stat(filepath.Join(dst, path)); !os.IsNotExist(err) {
					t.Errorf("disabled export remains: %s (%v)", path, err)
				}
			}
			if strings.Contains(read("/"), "data-copy-page") || strings.Contains(read("/"), `type="text/markdown"`) {
				t.Fatal("disabled export controls remain")
			}
		})
	}
}

func TestAgentExportsDefaultsAndCustomLayout(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	for name, body := range map[string]string{
		"index.md":     "# Home",
		"_layout.html": `<a href="{{ .MarkdownURL }}">Markdown</a><a href="{{ .LLMsURL }}">Catalog</a>{{ .Content }}`,
		"config.toml":  "[agents]\nfull = true\n",
	} {
		if err := os.WriteFile(filepath.Join(src, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	cfg, err := LoadConfig(filepath.Join(src, "config.toml"))
	if err != nil || !cfg.AgentsEnabled() || !cfg.Agents.Full {
		t.Fatalf("agents config: %#v, %v", cfg.Agents, err)
	}
	if err := Build(src, dst, Config{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dst, "llms-full.txt")); !os.IsNotExist(err) {
		t.Fatal("full bundle should be opt-in")
	}
	data, err := os.ReadFile(filepath.Join(dst, "index.html"))
	if err != nil || !strings.Contains(string(data), `href="/index.md"`) || !strings.Contains(string(data), `href="/llms.txt"`) {
		t.Fatalf("custom layout missing export URLs: %s (%v)", data, err)
	}
}

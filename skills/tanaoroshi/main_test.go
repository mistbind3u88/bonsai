package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestRefsKeepsSameRefFromDifferentSources(t *testing.T) {
	input := `{
  "owner/repo": {
    "issues": [
      {"number": 1, "body": "blocks owner/other#10"},
      {"number": 2, "body": "also blocks owner/other#10"}
    ],
    "prs": []
  }
}`
	path := filepath.Join(t.TempDir(), "collect.json")
	if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() {
		refs([]string{path})
	})

	var entries []struct {
		Source string `json:"source"`
		Ref    string `json:"ref"`
	}
	if err := json.Unmarshal([]byte(out), &entries); err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 source/ref entries, got %d: %s", len(entries), out)
	}
	if entries[0].Source == entries[1].Source {
		t.Fatalf("expected different sources, got %#v", entries)
	}
	if entries[0].Ref != "owner/other#10" || entries[1].Ref != "owner/other#10" {
		t.Fatalf("unexpected refs: %#v", entries)
	}
}

func TestExtractRefsIgnoresURLFragments(t *testing.T) {
	refs := extractRefs("see https://github.com/owner/repo/pull/123#456 and [owner/other#9](https://example.com)", "owner/repo")

	want := []string{"owner/repo#123", "owner/other#9"}
	if len(refs) != len(want) {
		t.Fatalf("expected %v, got %v", want, refs)
	}
	for i := range want {
		if refs[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, refs)
		}
	}
}

func TestParsePaginatedListFlattensPages(t *testing.T) {
	out := []byte(`[
  [{"body":"first","created_at":"2026-06-01T00:00:00Z"}],
  [{"body":"second","created_at":"2026-06-02T00:00:00Z"}]
]`)
	items, err := parsePaginatedList(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0]["body"] != "first" || items[1]["body"] != "second" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestIsGhNotFoundRequiresHTTP404(t *testing.T) {
	if !isGhNotFound([]byte("gh: Not Found (HTTP 404)")) {
		t.Fatal("expected GitHub 404 output to be treated as not found")
	}
	if isGhNotFound([]byte("body contains Not Found but no status")) {
		t.Fatal("expected non-404 output not to be treated as not found")
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = original
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

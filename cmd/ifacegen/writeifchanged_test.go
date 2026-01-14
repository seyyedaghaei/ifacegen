package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteIfChanged_NoChange(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "iface_gen.go")

	content := []byte("hello\n")
	if err := os.WriteFile(p, content, 0o644); err != nil {
		t.Fatalf("write initial file: %v", err)
	}

	wrote, err := writeIfChanged(p, content)
	if err != nil {
		t.Fatalf("writeIfChanged returned error: %v", err)
	}
	if wrote {
		t.Fatalf("expected wrote=false when content is unchanged")
	}

	after, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read after: %v", err)
	}
	if string(after) != string(content) {
		t.Fatalf("file content changed unexpectedly")
	}
}

func TestWriteIfChanged_ChangeAndCreate(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "iface_gen.go")

	// Create when file doesn't exist.
	wrote, err := writeIfChanged(p, []byte("v1\n"))
	if err != nil {
		t.Fatalf("writeIfChanged create returned error: %v", err)
	}
	if !wrote {
		t.Fatalf("expected wrote=true when file doesn't exist")
	}

	// Update when content differs.
	wrote, err = writeIfChanged(p, []byte("v2\n"))
	if err != nil {
		t.Fatalf("writeIfChanged update returned error: %v", err)
	}
	if !wrote {
		t.Fatalf("expected wrote=true when content differs")
	}

	after, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read after: %v", err)
	}
	if string(after) != "v2\n" {
		t.Fatalf("unexpected file content: %q", string(after))
	}
}

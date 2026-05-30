package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGolden(t *testing.T) {
	merged, locales, keys, err := loadTranslations([]string{"example"}, io.Discard)
	if err != nil {
		t.Fatalf("loadTranslations: %v", err)
	}

	outDir := t.TempDir()
	if err := writeKeysFile(filepath.Join(outDir, "keys.go"), keys); err != nil {
		t.Fatalf("writeKeysFile: %v", err)
	}
	if err := writeCatalogFile(filepath.Join(outDir, "catalog.go"), keys, locales, merged); err != nil {
		t.Fatalf("writeCatalogFile: %v", err)
	}

	for _, name := range []string{"catalog.go", "keys.go"} {
		got, err := os.ReadFile(filepath.Join(outDir, name))
		if err != nil {
			t.Fatalf("reading generated %s: %v", name, err)
		}
		want, err := os.ReadFile(filepath.Join("testdata/golden", name))
		if err != nil {
			t.Fatalf("reading golden %s: %v", name, err)
		}
		if string(got) != string(want) {
			t.Errorf("%s mismatch:\n--- want ---\n%s\n--- got ---\n%s", name, want, got)
		}
	}
}

func TestDuplicateKey(t *testing.T) {
	_, _, _, err := loadTranslations([]string{"testdata/dup"}, io.Discard)
	if err == nil {
		t.Fatal("expected error on duplicate key")
	}
	if !strings.Contains(err.Error(), "duplicate key") {
		t.Errorf("expected 'duplicate key' in error, got: %v", err)
	}
}

func TestMissingLocaleFallback(t *testing.T) {
	var buf bytes.Buffer
	_, _, _, err := loadTranslations([]string{"testdata/partial"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	warn := buf.String()
	if !strings.Contains(warn, "WARN:") {
		t.Errorf("expected WARN in output, got: %s", warn)
	}
	if !strings.Contains(warn, "falling back to en") {
		t.Errorf("expected 'falling back to en', got: %s", warn)
	}
}

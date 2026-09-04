package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateExportPathEmpty(t *testing.T) {
	if err := ValidateExportPath(""); err != nil {
		t.Fatalf(
			"ValidateExportPath(\"\") returned error: %v",
			err,
		)
	}
}

func TestValidateExportPathExistingWritableDirectory(t *testing.T) {
	exportDir := t.TempDir()

	if err := ValidateExportPath(exportDir); err != nil {
		t.Fatalf(
			"expected writable directory to pass, got: %v",
			err,
		)
	}

	testFile := filepath.Join(
		exportDir,
		".certmole-write-test",
	)

	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Fatalf(
			"expected write-test file to be removed, stat error: %v",
			err,
		)
	}
}

func TestValidateExportPathExistingFile(t *testing.T) {
	root := t.TempDir()

	filePath := filepath.Join(root, "results.csv")

	if err := os.WriteFile(
		filePath,
		[]byte("test"),
		0o644,
	); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	err := ValidateExportPath(filePath)

	if err == nil {
		t.Fatal("expected error for existing file")
	}
}

func TestValidateExportPathNonexistentDirectory(t *testing.T) {
	root := t.TempDir()

	exportDir := filepath.Join(
		root,
		"does-not-exist",
	)

	if err := ValidateExportPath(exportDir); err != nil {
		t.Fatalf(
			"expected nonexistent export directory to pass, got: %v",
			err,
		)
	}
}

func TestValidateExportPathNestedNonexistentDirectory(t *testing.T) {
	root := t.TempDir()

	exportDir := filepath.Join(
		root,
		"one",
		"two",
		"three",
	)

	if err := ValidateExportPath(exportDir); err != nil {
		t.Fatalf(
			"expected nonexistent nested directory to pass, got: %v",
			err,
		)
	}
}

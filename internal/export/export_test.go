package export

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"certmole/internal/crawler"
	"certmole/internal/parser"
)

func TestExportCSVToExistingDirectory(t *testing.T) {
	exportDir := t.TempDir()

	result := testScanResult()

	err := ExportCSV(exportDir, result)
	if err != nil {
		t.Fatalf("ExportCSV() returned error: %v", err)
	}

	exportFile := filepath.Join(exportDir, "certmole-results.csv")

	assertCSV(t, exportFile, result)
}

func TestExportCSVToNewDirectory(t *testing.T) {
	root := t.TempDir()
	exportDir := filepath.Join(root, "nested", "exports")

	result := testScanResult()

	err := ExportCSV(exportDir, result)
	if err != nil {
		t.Fatalf("ExportCSV() returned error: %v", err)
	}

	exportFile := filepath.Join(exportDir, "certmole-results.csv")

	assertCSV(t, exportFile, result)
}

func TestExportCSVToExplicitFilename(t *testing.T) {
	root := t.TempDir()

	exportFile := filepath.Join(
		root,
		"nested",
		"results.csv",
	)

	result := testScanResult()

	err := ExportCSV(exportFile, result)
	if err != nil {
		t.Fatalf("ExportCSV() returned error: %v", err)
	}

	assertCSV(t, exportFile, result)
}

func TestExportCSVToExistingFile(t *testing.T) {
	root := t.TempDir()

	exportFile := filepath.Join(root, "results.csv")

	if err := os.WriteFile(
		exportFile,
		[]byte("old content"),
		0o644,
	); err != nil {
		t.Fatalf("create existing file: %v", err)
	}

	result := testScanResult()

	err := ExportCSV(exportFile, result)
	if err != nil {
		t.Fatalf("ExportCSV() returned error: %v", err)
	}

	assertCSV(t, exportFile, result)
}

func TestExportCSVInvalidParent(t *testing.T) {
	result := testScanResult()

	exportFile := filepath.Join(
		t.TempDir(),
		"file-that-is-not-a-directory",
		"results.csv",
	)

	parent := filepath.Dir(exportFile)

	if err := os.WriteFile(
		parent,
		[]byte("not a directory"),
		0o644,
	); err != nil {
		t.Fatalf("create invalid parent: %v", err)
	}

	err := ExportCSV(exportFile, result)

	if err == nil {
		t.Fatal("expected ExportCSV() to return an error")
	}
}

func assertCSV(
	t *testing.T,
	path string,
	result *crawler.ScanResult,
) {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open CSV: %v", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			t.Errorf("close CSV: %v", err)
		}
	}()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("read CSV: %v", err)
	}

	expectedRows := len(result.Results) + 1

	if len(rows) != expectedRows {
		t.Fatalf(
			"expected %d rows, got %d",
			expectedRows,
			len(rows),
		)
	}

	expectedHeader := []string{
		"FILE PATH",
		"TYPE",
		"STATUS",
		"EXPIRES",
	}

	if len(rows[0]) != len(expectedHeader) {
		t.Fatalf(
			"expected %d header fields, got %d",
			len(expectedHeader),
			len(rows[0]),
		)
	}

	for i := range expectedHeader {
		if rows[0][i] != expectedHeader[i] {
			t.Errorf(
				"header field %d: expected %q, got %q",
				i,
				expectedHeader[i],
				rows[0][i],
			)
		}
	}

	for i, item := range result.Results {
		row := rows[i+1]

		expected := []string{
			item.Path,
			item.Asset.Type,
			item.Asset.Status,
			item.Asset.Expires,
		}

		for j := range expected {
			if row[j] != expected[j] {
				t.Errorf(
					"row %d field %d: expected %q, got %q",
					i,
					j,
					expected[j],
					row[j],
				)
			}
		}
	}
}

func testScanResult() *crawler.ScanResult {
	return &crawler.ScanResult{
		Results: []crawler.Result{
			{
				Path: "/tmp/example.pem",
				Asset: parser.Asset{
					Type:    "PEM Cert",
					Status:  "✅ Valid",
					Expires: "2027-01-01",
				},
			},
			{
				Path: "/tmp/private.key",
				Asset: parser.Asset{
					Type:    "⚠️ PRIVATE KEY",
					Status:  "EXPOSED KEY",
					Expires: "N/A",
				},
			},
		},
		CertificateCount: 1,
		PrivateKeyCount:  1,
		ValidCount:       1,
		ExpiredCount:     0,
	}
}

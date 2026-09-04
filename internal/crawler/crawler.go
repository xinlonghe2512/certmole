package crawler

import (
	"os"
	"path/filepath"
	"strings"

	"certmole/internal/parser"
)

type Result struct {
	Path  string
	Asset parser.Asset
}

type ScanResult struct {
	Results          []Result
	CertificateCount int
	PrivateKeyCount  int
	ValidCount       int
	ExpiredCount     int
}

var scanExtensions = []string{
	".crt",
	".pem",
	".der",
	".cer",
	".key",
}

func Run(scanDir string) (*ScanResult, error) {
	result := &ScanResult{}

	err := filepath.WalkDir(scanDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() || !hasScanExtension(path) {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		assets, _, _ := parser.AuditData(data)

		for _, asset := range assets {
			result.Results = append(result.Results, Result{
				Path:  path,
				Asset: asset,
			})

			switch asset.Type {
			case "PEM Cert", "DER Cert":
				result.CertificateCount++

				switch asset.Status {
				case "✅ Valid":
					result.ValidCount++

				case "❌ EXPIRED":
					result.ExpiredCount++
				}

			case "⚠️ PRIVATE KEY":
				result.PrivateKeyCount++
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func hasScanExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))

	for _, extension := range scanExtensions {
		if ext == extension {
			return true
		}
	}

	return false
}

package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"certmole/internal/crawler"
)

func ExportCSV(exportPath string, result *crawler.ScanResult) (err error) {
	info, err := os.Stat(exportPath)

	if err == nil && info.IsDir() {
		exportPath = filepath.Join(exportPath, "certmole-results.csv")
	} else if os.IsNotExist(err) {
		// Treat a path without a .csv extension as an export directory.
		if filepath.Ext(exportPath) == "" {
			if err := os.MkdirAll(exportPath, 0o755); err != nil {
				return fmt.Errorf("create export directory: %w", err)
			}

			exportPath = filepath.Join(exportPath, "certmole-results.csv")
		} else {
			// It looks like an explicit filename.
			parentDir := filepath.Dir(exportPath)

			if err := os.MkdirAll(parentDir, 0o755); err != nil {
				return fmt.Errorf("create export directory: %w", err)
			}
		}
	} else if err != nil {
		return fmt.Errorf("check export path: %w", err)
	}

	file, err := os.Create(exportPath)
	if err != nil {
		return fmt.Errorf("create CSV file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close CSV file: %w", closeErr)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{
		"FILE PATH",
		"TYPE",
		"STATUS",
		"EXPIRES",
	}); err != nil {
		return fmt.Errorf("write CSV header: %w", err)
	}

	for _, item := range result.Results {
		if err := writer.Write([]string{
			item.Path,
			item.Asset.Type,
			item.Asset.Status,
			item.Asset.Expires,
		}); err != nil {
			return fmt.Errorf("write CSV row: %w", err)
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush CSV: %w", err)
	}

	return nil
}

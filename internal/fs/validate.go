package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

func ValidateExportPath(exportPath string) error {
	if exportPath == "" {
		return nil
	}

	info, err := os.Stat(exportPath)

	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf(
				"export path is not a directory: %s",
				exportPath,
			)
		}

		testFile := filepath.Join(
			exportPath,
			".certmole-write-test",
		)

		file, err := os.Create(testFile)
		if err != nil {
			return fmt.Errorf(
				"cannot write to export directory %q: %w",
				exportPath,
				err,
			)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("close write test file: %w", err)
		}

		if err := os.Remove(testFile); err != nil {
			return fmt.Errorf("remove write test file: %w", err)
		}

		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf(
			"cannot access export path %q: %w",
			exportPath,
			err,
		)
	}

	parentDir := filepath.Dir(exportPath)

	if _, err := os.Stat(parentDir); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf(
				"cannot access parent directory %q: %w",
				parentDir,
				err,
			)
		}
	}

	return nil
}

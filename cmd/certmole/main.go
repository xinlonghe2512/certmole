package main

import (
	"fmt"
	"os"

	"certmole/internal/crawler"
	"certmole/internal/export"
	"certmole/internal/fs"
)

func main() {
	options := parseArgs()

	if err := fs.ValidateExportPath(options.ExportPath); err != nil {
		printError(err)
		os.Exit(1)
	}

	printScanBanner()

	printScanConfiguration(options)

	confirmed, err := printScanConfirmation()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	if !confirmed {
		printCancelled()
		return
	}

	printScanHeader(options.ScanDir)

	result, err := crawler.Run(options.ScanDir)
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	for _, item := range result.Results {
		printScanRow(
			item.Path,
			item.Asset.Type,
			item.Asset.Status,
			item.Asset.Expires,
		)
	}

	printScanResult(result)

	if options.ExportPath != "" {
		if err := export.ExportCSV(options.ExportPath, result); err != nil {
			printError(err)
			os.Exit(1)
		}

		fmt.Printf(
			"\nResults exported to %s\n",
			options.ExportPath,
		)
	}
}

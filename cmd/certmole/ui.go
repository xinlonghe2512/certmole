package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"certmole/internal/crawler"
)

const (
	reset = "\033[0m"

	bold = "\033[1m"

	cyan   = "\033[36m" // application/UI
	green  = "\033[32m" // valid
	yellow = "\033[33m" // warning/private key
	red    = "\033[31m" // expired/error
	gray   = "\033[90m" // secondary information
	white  = "\033[37m" // normal text
)

func printHelp() {
	fmt.Println(`Scan your machine for certificates and export the results

Usage: certmole [options...]

Options:
  -d, --directory <path>   Directory to recursively scan
  -e, --export <file>      Export scan results to a CSV file
  -v, --version            Show version number and quit
  -h, --help               Show help menu and quit

Examples:
  certmole --directory /etc/ssl
  certmole --directory /etc/ssl --export results.csv
  certmole --help
  certmole --version`)
}

func printScanBanner() {
	fmt.Println()
	fmt.Printf("%s%s>_ Certmole CLI (%s)%s\n", bold, green, version, reset)
	fmt.Printf("%s──────────────────────────────────%s\n", cyan, reset)
	fmt.Println()
	fmt.Printf("%s%sCertificate and private-key discovery tool.%s\n", bold, gray, reset)
	fmt.Printf("%s%sLocal. Single Binary. Zero-dependency.%s\n", bold, gray, reset)
	fmt.Println()
}

func printScanConfiguration(options *ParseArgsOptions) {
	fmt.Printf(
		"%sScan configuration%s\n",
		cyan,
		reset,
	)

	fmt.Printf(
		"  %s(1) Scan Directory :%s %s\n",
		gray,
		reset,
		options.ScanDir,
	)

	if options.ScanDir == "/" {
		fmt.Printf(
			"   %s⚠  This will recursively scan the entire filesystem.%s\n",
			yellow,
			reset,
		)
	}

	fmt.Println()

	if options.ExportPath != "" {
		fmt.Printf(
			"  %s(2) Export Path    :%s %s\n",
			gray,
			reset,
			options.ExportPath,
		)

		if _, err := os.Stat(options.ExportPath); os.IsNotExist(err) {
			fmt.Printf(
				"   %s⚠  The export path does not exist and will be created.%s\n",
				yellow,
				reset,
			)
		}
	} else {
		fmt.Printf(
			"  %s(2) Export Path    :%s NIL\n",
			gray,
			reset,
		)
	}

	fmt.Printf("\n\n")
}

func printScanConfirmation() (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Begin scanning? [y/N]: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		switch strings.ToLower(strings.TrimSpace(input)) {
		case "y", "yes":
			return true, nil

		case "", "n", "no":
			return false, nil

		default:
			fmt.Printf(
				"  %sPlease choose 'y' or 'n'.%s\n\n",
				gray,
				reset,
			)
		}
	}
}

func printScanHeader(scanDir string) {
	fmt.Printf(
		"\nScanning '%s' ...\n\n",
		scanDir,
	)

	fmt.Printf(
		"%-50s | %-15s | %-12s | %s\n",
		"FILE PATH",
		"TYPE",
		"STATUS",
		"EXPIRES",
	)

	fmt.Println(strings.Repeat("-", 110))
}

func printScanRow(path, assetType, status, expires string) {
	displayPath := path

	if len(displayPath) > 50 {
		displayPath = "..." + displayPath[len(displayPath)-47:]
	}

	fmt.Printf(
		"%-50s | %-15s | %-12s | %s\n",
		displayPath,
		assetType,
		status,
		expires,
	)
}

func printScanResult(result *crawler.ScanResult) {
	fmt.Println()
	fmt.Println(strings.Repeat("─", 110))
	fmt.Println()

	fmt.Println("Scan complete.")
	fmt.Println()

	fmt.Printf(
		"  Certificates found : %d\n",
		result.CertificateCount,
	)

	fmt.Printf(
		"  Private keys found : %d\n",
		result.PrivateKeyCount,
	)

	fmt.Println()

	fmt.Printf(
		"  Valid certificates : %s%d%s\n",
		green,
		result.ValidCount,
		reset,
	)

	fmt.Printf(
		"  Expired certificates: %s%d%s\n",
		red,
		result.ExpiredCount,
		reset,
	)

	fmt.Println()

	findings := result.ExpiredCount + result.PrivateKeyCount

	if findings == 0 {
		fmt.Printf(
			"  %s✓ No security findings detected.%s\n",
			green,
			reset,
		)
	} else {
		fmt.Printf(
			"  %s⚠ Security findings : %d%s\n",
			yellow,
			findings,
			reset,
		)
	}

	fmt.Println()
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
}

func printCancelled() {
	fmt.Println("\nScan aborted.")
}

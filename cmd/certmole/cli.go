package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type ParseArgsOptions struct {
	ScanDir    string
	ExportPath string
}

func parseArgs() *ParseArgsOptions {
	flags := flag.NewFlagSet("certmole", flag.ContinueOnError)

	// Suppress the default usage output when parsing fails.
	flags.SetOutput(io.Discard)

	var scanDir string
	var exportPath string
	var versionFlag bool

	// Directory
	flags.StringVar(
		&scanDir,
		"directory",
		"",
		"Directory to recursively scan",
	)

	flags.StringVar(
		&scanDir,
		"d",
		"",
		"",
	)

	// Export
	flags.StringVar(
		&exportPath,
		"export",
		"",
		"Export scan results to a CSV file",
	)

	flags.StringVar(
		&exportPath,
		"e",
		"",
		"",
	)

	// Version
	flags.BoolVar(
		&versionFlag,
		"version",
		false,
		"Show version number and quit",
	)

	flags.BoolVar(
		&versionFlag,
		"v",
		false,
		"",
	)

	// Parse arguments.
	err := flags.Parse(os.Args[1:])

	if err == flag.ErrHelp {
		flags.SetOutput(os.Stdout)
		printHelp()
		// fmt.Fprintf(os.Stdout, "Usage: certmole [options...]\n\n")
		// flags.PrintDefaults()
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "certmole: %v\n", err)
		fmt.Fprintln(
			os.Stderr,
			"certmole: try 'certmole --help' for more information.",
		)
		os.Exit(2)
	}

	// Handle version.
	if versionFlag {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	// Directory is required.
	if scanDir == "" {
		fmt.Fprintln(
			os.Stderr,
			"certmole: -d or --directory is required.",
		)
		fmt.Fprintln(
			os.Stderr,
			"certmole: try 'certmole --help' for more information.",
		)
		os.Exit(2)
	}

	return &ParseArgsOptions{
		ScanDir:    scanDir,
		ExportPath: exportPath,
	}
}

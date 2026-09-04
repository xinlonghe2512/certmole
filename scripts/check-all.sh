#!/usr/bin/env bash

set -euo pipefail

echo "==> Checking formatting..."
if [ -n "$(gofmt -l .)" ]; then
    echo "ERROR: Go files are not formatted."
    gofmt -l .
    exit 1
fi

echo "==> Verifying modules..."
go mod verify

echo "==> Running tests..."
go test ./...

echo "==> Running vet..."
go vet ./...

echo "==> Running golangci-lint..."
golangci-lint run

echo "==> Running govulncheck..."
govulncheck ./...

echo "==> Checking shell scripts..."
shellcheck scripts/*.sh

echo "==> Building..."
go build -o certmole ./cmd/certmole

echo "==> Smoke testing..."
printf 'n\n' | ./certmole --directory .

rm -f certmole

echo
echo "All checks passed."
.PHONY: build run test fmt vet lint vuln tidy check clean release

# Compile the certmole binary for the current platform
build:
	go build -o dist/certmole ./cmd/certmole

# Build and run certmole against the current directory
run: build
	./certmole .

# Run the Go test suite
test:
	go test ./...

# Format all Go source files using gofumpt
fmt:
	gofumpt -w .

# Run Go's built-in static analysis checks
vet:
	go vet ./...

# Run the project's configured linters
lint:
	golangci-lint run

# Scan project dependencies and source code for known vulnerabilities
vuln:
	govulncheck ./...

# Synchronize and clean up Go module dependencies
tidy:
	go mod tidy

# Run all project validation and integrity checks
check:
	./scripts/check-all.sh

# Remove locally built binaries and release artifacts
clean:
	rm -f certmole
	rm -rf dist

# Validate the project and build binaries for all supported platforms
release: check
	./scripts/build-all.sh
#!/usr/bin/env bash

# Exit on errors
set -e

BINARY_NAME="certmole"
OUTPUT_DIR="dist/releases"

# Clean previous distribution builds
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Target Operating Systems and Architectures (GOOS/GOARCH)
OS_TARGETS=(
    "linux/amd64"   # Linux AMD64
    "linux/arm64"   # Linux ARM64
    "windows/amd64" # Windows AMD64
)

echo "==> Starting cross-compilation for ${BINARY_NAME}..."

for OS_TARGET in "${OS_TARGETS[@]}"; do
    IFS="/" read -r -a ARR <<< "$OS_TARGET"
    OS="${ARR[0]}"
    ARCH="${ARR[1]}"
    
    OUTPUT_NAME="${OUTPUT_DIR}/${BINARY_NAME}-${OS}-${ARCH}"
    
    if [ "$OS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo "==> Building for ${OS} (${ARCH})..."
    # CGO_ENABLED=0 ensures a static binary with no external C library dependencies
    GOOS="$OS" GOARCH="$ARCH" CGO_ENABLED=0 go build -o "$OUTPUT_NAME" ./cmd/certmole
done

echo "==> Generating SHA-256 checksums..."

(
    cd "$OUTPUT_DIR"
    sha256sum certmole-* > checksums.txt
)

echo "==> System binaries compiled successfully inside /${OUTPUT_DIR}!"
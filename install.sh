#!/usr/bin/env bash

set -euo pipefail

REPO_OWNER="xinlonghe2512"
REPO_NAME="certmole"
BINARY_NAME="certmole"

info() {
    printf '==> %s\n' "$*"
}

success() {
    printf '%s\n' "$*"
}

fail() {
    printf '\nERROR: %s\n' "$*" >&2
    exit 1
}

# Check dependencies
command -v curl >/dev/null 2>&1 || fail "curl is required."
command -v tar >/dev/null 2>&1 || fail "tar is required."

# Detect operating system
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"

case "$OS" in
    linux)
        ;;
    *)
        fail "Unsupported operating system: $OS"
        ;;
esac

# Detect architecture
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        fail "Unsupported processor architecture: $ARCH"
        ;;
esac

info "Detected platform: $OS-$ARCH"

# Resolve latest release version
RELEASE_TAG="$(
    curl -fsSL \
        "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" |
        sed -n 's/.*"tag_name": "\(.*\)".*/\1/p'
)"

[ -n "$RELEASE_TAG" ] ||
    fail "Could not determine the latest Certmole release."

# Git tags use the "v" prefix, while artifact names do not.
RELEASE_VERSION="${RELEASE_TAG#v}"

info "Resolved version: $RELEASE_VERSION"

# Target release asset
TARGET_ARCHIVE="${BINARY_NAME}-${RELEASE_VERSION}-${OS}-${ARCH}.tar.gz"

# GitHub release download URL
DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${RELEASE_TAG}/${TARGET_ARCHIVE}"

# Determine installation directory
if [ "${EUID:-$(id -u)}" -eq 0 ]; then
    INSTALL_DIR="/usr/local/bin"
else
    INSTALL_DIR="${HOME}/.local/bin"
fi

INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}"

# Check currently installed version
CURRENT_VERSION="not installed"

if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    CURRENT_VERSION="$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")"
fi

if [ "$CURRENT_VERSION" = "$RELEASE_VERSION" ]; then
    success "Certmole CLI $CURRENT_VERSION is already up to date."
    exit 0
elif [ "$CURRENT_VERSION" = "not installed" ]; then
    info "Certmole CLI is not installed."
else
    info "Updating Certmole CLI from $CURRENT_VERSION to $RELEASE_VERSION"
fi

# Create installation directory
if [ ! -d "$INSTALL_DIR" ]; then
    info "Creating installation directory..."

    mkdir -p "$INSTALL_DIR" ||
        fail "Could not create $INSTALL_DIR"

    success "Created $INSTALL_DIR"
fi

# Create temporary files/directories
TMP_FILE="$(mktemp)"
TMP_DIR="$(mktemp -d)"

trap 'rm -f "$TMP_FILE"; rm -rf "$TMP_DIR"' EXIT

# Download archive
if ! curl -fLsS "$DOWNLOAD_URL" -o "$TMP_FILE"; then
    fail "Failed to download $TARGET_ARCHIVE."
fi

# Extract archive
if ! tar -xzf "$TMP_FILE" -C "$TMP_DIR"; then
    fail "Failed to extract $TARGET_ARCHIVE."
fi

# Verify extracted binary
if [ ! -f "$TMP_DIR/$BINARY_NAME" ]; then
    fail "Binary '$BINARY_NAME' was not found in the archive."
fi

# Install
info "Installing Certmole CLI $RELEASE_VERSION to $INSTALL_PATH"

chmod +x "$TMP_DIR/$BINARY_NAME" ||
    fail "Could not make binary executable."

mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_PATH" ||
    fail "Could not install to $INSTALL_PATH."

# Verify installation
info "Verifying installation..."

if ! "$INSTALL_PATH" --help >/dev/null 2>&1; then
    fail "Binary was installed but could not be executed."
fi

INSTALLED_VERSION="$("$INSTALL_PATH" --version 2>/dev/null)" ||
    fail "Could not determine installed version."

success "Certmole CLI $INSTALLED_VERSION installed successfully."

printf '\n'

# Check whether it is already on PATH
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    printf 'Get started by running either:\n\n'
    printf '  certmole --help\n'
    printf '  certmole --directory .\n'
    printf '\n'
else
    printf 'Note: %s is not currently in your PATH.\n\n' "$INSTALL_DIR"

    case "${SHELL:-}" in
        */fish)
            printf 'For fish:\n\n'
            printf '  fish_add_path %s\n' "$INSTALL_DIR"
            ;;
        *)
            printf 'For bash/zsh:\n\n'
            printf '  export PATH="%s:$PATH"\n' "$INSTALL_DIR"
            ;;
    esac

    printf '\nThen run:\n\n'
    printf '  certmole --help\n'
fi

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

TARGET_BINARY="${BINARY_NAME}-${OS}-${ARCH}"

# Determine installation directory
if [ "${EUID:-$(id -u)}" -eq 0 ]; then
    INSTALL_DIR="/usr/local/bin"
else
    INSTALL_DIR="${HOME}/.local/bin"
fi

INSTALL_DIR="${HOME}/.local/bin"
INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}"

# GitHub's latest-release asset endpoint.
DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/latest/download/${TARGET_BINARY}"

# Resolve latest release version
RELEASE_VERSION="$(
    curl -fsSL \
        "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" |
        sed -n 's/.*"tag_name": "\(.*\)".*/\1/p'
)"

[ -n "$RELEASE_VERSION" ] || fail "Could not determine the latest Certmole version."

RELEASE_VERSION="${RELEASE_VERSION#v}"

CURRENT_VERSION="not installed"

if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    CURRENT_VERSION="$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")"
fi

if [ "$CURRENT_VERSION" = "$RELEASE_VERSION" ]; then
    success "Certmole CLI $CURRENT_VERSION is already up to date."
    exit 0
elif [ "$CURRENT_VERSION" = "not installed" ]; then
    info "Certmole CLI is not installed."
    info "Installing Certmole CLI $RELEASE_VERSION"
else
    info "Updating Certmole CLI from $CURRENT_VERSION to $RELEASE_VERSION"
fi

info "Resolved version: $RELEASE_VERSION"
info "Detected platform: $OS-$ARCH"

# Create installation directory
if [ ! -d "$INSTALL_DIR" ]; then
    info "Creating installation directory..."
    mkdir -p "$INSTALL_DIR" || fail "Could not create $INSTALL_DIR"
    success "Created $INSTALL_DIR"
fi

# Download to a temporary file first
TMP_FILE="$(mktemp)"
trap 'rm -f "$TMP_FILE"' EXIT

if ! curl -fLsS "$DOWNLOAD_URL" -o "$TMP_FILE"; then
    fail "Failed to download ${TARGET_BINARY}."
fi

# Install
info "Installing standalone package to $INSTALL_PATH/$TARGET_BINARY"

chmod +x "$TMP_FILE" || fail "Could not make binary executable."

mv "$TMP_FILE" "$INSTALL_PATH" || fail "Could not install to $INSTALL_PATH."

# Verify
info "Verifying installation..."

if ! "$INSTALL_PATH" --help >/dev/null 2>&1; then
    fail "Binary was installed but could not be executed."
fi

VERSION="$("$INSTALL_PATH" --version)"

success "Certmole CLI $VERSION installed successfully."
printf '\n'

# Check whether it is already on PATH
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    printf 'Get started by running either:
  certmole --help
  certmole --directory .

'
else
    printf 'Note: %s is not currently in your PATH.

' "$INSTALL_DIR"

    case "${SHELL:-}" in
        */fish)
            printf 'For fish:
  fish_add_path %s

' "$INSTALL_DIR"
            ;;
        *)
            printf 'For bash/zsh:
  export PATH="%s:$PATH"

' "$INSTALL_DIR"
            ;;
    esac

    printf 'Then run:
  certmole --help

'
fi
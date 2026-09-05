FROM mcr.microsoft.com/devcontainers/go:1.26-trixie

USER root

# Install general development/build tooling
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  bash \
  build-essential \
  ca-certificates \
  curl \
  gh \
  jq \
  ripgrep \
  fd \
  fzf \
  bat \
  make \
  shellcheck \
  tar \
  unzip \
  wget \
  zip \
  && rm -rf /var/lib/apt/lists/*

USER vscode

ARG GOLANGCI_LINT_VERSION=v2.8.0

# Install latest Go core development tools
RUN go install golang.org/x/tools/gopls@latest \
  && go install golang.org/x/tools/cmd/goimports@latest \
  && go install mvdan.cc/gofumpt@latest \
  && go install honnef.co/go/tools/cmd/staticcheck@latest \
  && go install github.com/go-delve/delve/cmd/dlv@latest \
  && go install golang.org/x/vuln/cmd/govulncheck@latest

# Install golangci-lint
RUN curl -sSfL https://golangci-lint.run/install.sh \
  | sh -s -- -b "$(go env GOPATH)/bin" "${GOLANGCI_LINT_VERSION}"

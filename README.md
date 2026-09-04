<div align="center">
  <a href="https://github.com/xinlonghe2512/certmole/actions"><img src="https://github.com/xinlonghe2512/certmole/actions/workflows/ci.yaml/badge.svg?branch=main" alt="Build Status"></a>
  <a href="http://github.com/xinlonghe2512/certmole/releases"><img src="https://img.shields.io/github/v/tag/xinlonghe2512/certmole" alt="Version"></a>
  <a href="https://github.com/xinlonghe2512/certmole/releases"><img src="https://img.shields.io/github/v/tag/xinlonghe2512/certmole?label=version&color=blue" alt="Version"></a>
  <a href="https://github.com/xinlonghe2512/certmole/blob/main/LICENSE"><img src="https://img.shields.io/github/license/xinlonghe2512/certmole" alt="License"></a>
  <a href="https://github.com/xinlonghe2512/certmole/graphs/contributors"><img src="https://img.shields.io/github/contributors/xinlonghe2512/certmole" alt="Contributors"></a>
  <a href="https://github.com/xinlonghe2512/certmole/stargazers"><img src="https://img.shields.io/github/stars/xinlonghe2512/certmole?style=flat" alt="Stars"></a>
</div>

# Certmole

Certmole is a command-line tool made for System Engineers or DevSecOps personnel for discovering cryptographic assets across filesystems. It recursively scans directories for X.509 certificates and private keys, identifies certificate expiration status, and reports exposed private keys.

Built entirely in Go, Certmole compiles into a single standalone binary with no runtime dependencies, making it suitable for local security audits, infrastructure maintenance, container environments, and CI/CD pipelines.

## Features

- **Recursive filesystem scanning** — Walks a target directory recursively to discover supported cryptographic files.

- **X.509 certificate discovery** — Detects certificates in both PEM and DER formats.

- **Certificate expiration detection** — Identifies valid and expired certificates and reports their expiration dates.

- **Private-key detection** — Detects PEM-encoded private keys and explicitly flags them as exposed cryptographic assets.

- **Permission-tolerant scanning** — Skips files and directories that cannot be accessed instead of terminating the entire scan.

- **CSV export** — Exports discovered assets and their metadata to a CSV file for reporting and further analysis.

- **Zero runtime dependencies** — Compiles into a standalone Go binary. No Python, Node.js, Java, or other runtime is required.

- **Cross-platform builds** — Supports Linux and Windows on common AMD64 and ARM64 architectures.

## Installation

### Quickstart

#### Linux

Install Certmole using the installation script:

```shell
curl -fsSL https://raw.githubusercontent.com/xinlonghe2512/certmole/main/install.sh | sh
```

#### Windows

```shell
powershell -ExecutionPolicy ByPass -c "irm https://raw.githubusercontent.com/xinlonghe2512/certmole/main/install.ps1 | iex"
```

### Using Git

Alternatively, you can "git clone" this repository to any directory and run install script.

```shell
git clone --depth 1 https://github.com/xinlonghe2512/certmole.git

chmod +x install.sh

./install.sh

```

> [!IMPORTANT]
> Only download Certmole from the official repository or its published release artifacts. Review installation scripts before executing them in security-sensitive environments.

### Binary releases

You can download the official certmole binaries from the releases page.

- [https://github.com/xinlonghe2512/certmole/releases](https://github.com/xinlonghe2512/certmole/releases
)

## Usage

Certmole requires a target directory to scan.

### Scan the current directory

```shell
certmole --directory .
```

### Scan a specific target directory

```shell
certmole --directory /etc/nginx/certs
```

> [!IMPORTANT]
Ensure that you have the necessary directory permissions to access the target directory before running Certmole.

### Scan and export results to CSV file (specifying a filename)

```shell
certmole --directory /etc/nginx/certs -export ./reports/certificates-report.csv
```

### Scan and export results to CSV file (without specifying a filename)

```shell
certmole --directory /etc/nginx/certs --export ./reports/
```

> [!NOTE]
When not specifying a filename, the default filename `certmole-result` will be used. In this example, you can find your result at `./reports/certmole-result.csv` .

## Supported File Types

Certmole currently examines files with the following certificate extensions.

| Extension | Format |
| :--- | :----: |
| .crt | PEM-encoded certificate or private key |
| .pem | DER-encoded certificate |
| .cer | Certificate |
| .key | Private Key |

> [!NOTE]
The parser identifies the actual cryptographic content rather than relying solely on the file extension.

## Contributing

Contributions, bug reports, documentation improvements, and feature proposals are welcome.

Please keep pull requests focused on a single change and ensure that all CI checks pass before requesting review.

For larger changes, open an issue first to discuss the proposed approach.

## Security

If you discover a security vulnerability in Certmole, please follow the security reporting instructions in [SECURITY.md](SECURITY).  rather than opening a public issue.

## License

Certmole is licensed under the [MIT License](LICENSE).

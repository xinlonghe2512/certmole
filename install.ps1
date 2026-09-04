$ErrorActionPreference = "Stop"

$RepoOwner = "xinlonghe2512"
$RepoName = "certmole"
$BinaryName = "certmole"

function Info {
    param([string]$Message)
    Write-Host "==> $Message"
}

function Success {
    param([string]$Message)
    Write-Host $Message
}

function Fail {
    param([string]$Message)
    Write-Host ""
    Write-Error $Message
    exit 1
}

# Detect processor architecture
$Arch = $env:PROCESSOR_ARCHITECTURE.ToLower()

switch ($Arch) {
    "amd64" {
        $Arch = "amd64"
    }

    "arm64" {
        $Arch = "arm64"
    }

    default {
        Fail "Unsupported Windows processor architecture: $Arch"
    }
}

$TargetBinary = "${BinaryName}-windows-${Arch}.exe"

# Installation directory
$InstallDir = Join-Path $env:LOCALAPPDATA "Programs\certmole"
$InstallPath = Join-Path $InstallDir "${BinaryName}.exe"

# GitHub latest release API
$ApiUrl = "https://api.github.com/repos/${RepoOwner}/${RepoName}/releases/latest"

# Resolve latest release version
try {
    $ReleaseInfo = Invoke-RestMethod -Uri $ApiUrl -Method Get
}
catch {
    Fail "Failed to communicate with the GitHub Releases API: $($_.Exception.Message)"
}

$ReleaseVersion = $ReleaseInfo.tag_name

if ([string]::IsNullOrWhiteSpace($ReleaseVersion)) {
    Fail "Could not determine the latest Certmole version."
}

$ReleaseVersion = $ReleaseVersion -replace '^v', ''

# Determine currently installed version
$CurrentVersion = "not installed"

if (Test-Path $InstallPath) {
    try {
        $CurrentVersion = (& $InstallPath --version 2>$null).Trim()

        if ([string]::IsNullOrWhiteSpace($CurrentVersion)) {
            $CurrentVersion = "unknown"
        }
    }
    catch {
        $CurrentVersion = "unknown"
    }
}

# Check version
if ($CurrentVersion -eq $ReleaseVersion) {
    Success "Certmole CLI $CurrentVersion is already up to date."
    exit 0
}
elseif ($CurrentVersion -eq "not installed") {
    Info "Certmole CLI is not installed."
    Info "Installing Certmole CLI $ReleaseVersion"
}
else {
    Info "Updating Certmole CLI from $CurrentVersion to $ReleaseVersion"
}

Info "Resolved version: $ReleaseVersion"
Info "Detected platform: $OS-$Arch"

# Create installation directory
if (-not (Test-Path $InstallDir)) {
    Info "Creating installation directory..."

    try {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    catch {
        Fail "Could not create $InstallDir"
    }

    Success "Created $InstallDir"
}

# Locate release asset
$Asset = $ReleaseInfo.assets |
    Where-Object { $_.name -eq $TargetBinary } |
    Select-Object -First 1

if (-not $Asset) {
    Fail "Failed to locate a release matching: $TargetBinary"
}

$DownloadUrl = $Asset.browser_download_url

# Download
Info "Installing standalone package to $InstallPath"

try {
    Invoke-WebRequest `
        -Uri $DownloadUrl `
        -OutFile $InstallPath `
        -UseBasicParsing
}
catch {
    Fail "Failed to download ${TargetBinary}: $($_.Exception.Message)"
}

# Verify
Info "Verifying installation..."

try {
    $InstalledVersion = (& $InstallPath --version 2>$null).Trim()

    if ($LASTEXITCODE -ne 0) {
        Fail "Binary was installed but could not be executed."
    }
}
catch {
    Fail "Binary was installed but could not be executed."
}

Success "Certmole CLI $InstalledVersion installed successfully."

Write-Host ""

# Add installation directory to user PATH
$UserPath = [Environment]::GetEnvironmentVariable(
    "Path",
    [EnvironmentVariableTarget]::User
)

$PathEntries = @()

if (-not [string]::IsNullOrWhiteSpace($UserPath)) {
    $PathEntries = $UserPath -split ";" |
        Where-Object {
            -not [string]::IsNullOrWhiteSpace($_)
        }
}

$PathExists = $PathEntries |
    Where-Object {
        $_.TrimEnd("\") -ieq $InstallDir.TrimEnd("\")
    }

if (-not $PathExists) {
    $NewPath = ($PathEntries + $InstallDir) -join ";"

    [Environment]::SetEnvironmentVariable(
        "Path",
        $NewPath,
        [EnvironmentVariableTarget]::User
    )
}

# Quickstart
if ($PathExists) {
    Write-Host "Get started by running either:"
    Write-Host ""
    Write-Host "  certmole --help"
    Write-Host ""
    Write-Host "  certmole --directory ."
}
else {
    Write-Host "Note: $InstallDir was added to your user PATH."
    Write-Host ""
    Write-Host "Restart your PowerShell session, then run:"
    Write-Host ""
    Write-Host "  certmole --help"
}
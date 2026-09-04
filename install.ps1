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

# Detect operating system
$OS = "windows"

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

Info "Detected platform: $OS-$Arch"

# Installation directory
$InstallDir = Join-Path $env:LOCALAPPDATA "Programs\certmole"
$InstallPath = Join-Path $InstallDir "${BinaryName}.exe"

# GitHub latest release API
$ApiUrl = "https://api.github.com/repos/${RepoOwner}/${RepoName}/releases/latest"

# Resolve latest release
try {
    $ReleaseInfo = Invoke-RestMethod -Uri $ApiUrl -Method Get
}
catch {
    Fail "Failed to communicate with the GitHub Releases API: $($_.Exception.Message)"
}

$ReleaseTag = $ReleaseInfo.tag_name

if ([string]::IsNullOrWhiteSpace($ReleaseTag)) {
    Fail "Could not determine the latest Certmole version."
}

# Convert v0.1.1 -> 0.1.1
$ReleaseVersion = $ReleaseTag -replace '^v', ''

Info "Resolved version: $ReleaseVersion"

# Determine currently installed version
$CurrentVersion = "not installed"

if (Test-Path $InstallPath) {
    try {
        $CurrentVersionOutput = (& $InstallPath --version 2>$null).Trim()

        if ([string]::IsNullOrWhiteSpace($CurrentVersionOutput)) {
            $CurrentVersion = "unknown"
        }
        else {
            # Extract semantic version from output such as:
            #   0.1.1
            #   v0.1.1
            #   certmole v0.1.1
            #   certmole version 0.1.1
            $VersionMatch = [regex]::Match(
                $CurrentVersionOutput,
                'v?(\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?)'
            )

            if ($VersionMatch.Success) {
                $CurrentVersion = $VersionMatch.Groups[1].Value
            }
            else {
                $CurrentVersion = $CurrentVersionOutput
            }
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
}
else {
    Info "Updating Certmole CLI from $CurrentVersion to $ReleaseVersion"
}

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

# Locate versioned Windows release asset
$TargetArchive = "${BinaryName}-${ReleaseVersion}-windows-${Arch}.zip"

$Asset = $ReleaseInfo.assets |
    Where-Object { $_.name -eq $TargetArchive } |
    Select-Object -First 1

if (-not $Asset) {
    Fail "Failed to locate release asset: $TargetArchive"
}

$DownloadUrl = $Asset.browser_download_url

# Temporary download/extraction directory
$TempDir = Join-Path $env:TEMP "certmole-install-$([Guid]::NewGuid())"
$TempArchive = Join-Path $TempDir $TargetArchive
$ExtractDir = Join-Path $TempDir "extracted"

try {
    New-Item -ItemType Directory -Path $TempDir -Force | Out-Null
    New-Item -ItemType Directory -Path $ExtractDir -Force | Out-Null
}
catch {
    Fail "Could not create temporary installation directory."
}

try {
    # Download
    Invoke-WebRequest `
        -Uri $DownloadUrl `
        -OutFile $TempArchive `
        -UseBasicParsing

    # Extract
    Expand-Archive `
        -Path $TempArchive `
        -DestinationPath $ExtractDir `
        -Force

    # Locate binary
    $ExtractedBinary = Join-Path $ExtractDir "${BinaryName}.exe"

    if (-not (Test-Path $ExtractedBinary)) {
        Fail "Binary '$BinaryName.exe' was not found in the archive."
    }

    # Install
    Info "Installing Certmole CLI $ReleaseVersion to $InstallPath"

    Copy-Item `
        -Path $ExtractedBinary `
        -Destination $InstallPath `
        -Force

    # Verify
    Info "Verifying installation..."

    try {
        $InstalledVersionOutput = (& $InstallPath --version 2>$null).Trim()

        if ($LASTEXITCODE -ne 0) {
            Fail "Binary was installed but could not be executed."
        }

        if ([string]::IsNullOrWhiteSpace($InstalledVersionOutput)) {
            Fail "Binary was installed but did not return a version."
        }

        $VersionMatch = [regex]::Match(
            $InstalledVersionOutput,
            'v?(\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?)'
        )

        if ($VersionMatch.Success) {
            $InstalledVersion = $VersionMatch.Groups[1].Value
        }
        else {
            $InstalledVersion = $InstalledVersionOutput
        }
    }
    catch {
        Fail "Binary was installed but could not be executed."
    }

    Success "Certmole CLI $InstalledVersion installed successfully."
}
catch {
    Fail "Installation failed: $($_.Exception.Message)"
}
finally {
    # Clean up temporary files
    if (Test-Path $TempDir) {
        Remove-Item `
            -Path $TempDir `
            -Recurse `
            -Force `
            -ErrorAction SilentlyContinue
    }
}

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

    $PathWasAdded = $true
}
else {
    $PathWasAdded = $false
}

# Quickstart
Write-Host ""

if ($PathWasAdded) {
    Write-Host "Note: $InstallDir was added to your user PATH."
    Write-Host ""
    Write-Host "Restart your PowerShell session, then run:"
    Write-Host ""
    Write-Host "  certmole --help"
    Write-Host ""
    Write-Host "  certmole --directory ."
}
else {
    Write-Host "Get started by running either:"
    Write-Host ""
    Write-Host "  certmole --help"
    Write-Host ""
    Write-Host "  certmole --directory ."
}

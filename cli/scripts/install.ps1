$ErrorActionPreference = "Stop"

$Repo   = "aifunc-dev/aifunc"
$Binary = "aifn"
$InstallDir = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { "$env:LOCALAPPDATA\Programs\aifn" }

# Resolve latest version if not specified
$Version = $env:VERSION
if (-not $Version) {
    $release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
    $Version = $release.tag_name
}

if (-not $Version) {
    Write-Error "Failed to resolve latest version."
    exit 1
}

$Tag = $Version.TrimStart("v")

# Detect arch
$Arch = if ([System.Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Error "32-bit systems are not supported."
    exit 1
}

$Filename = "${Binary}_${Tag}_windows_${Arch}.zip"
$Url = "https://github.com/$Repo/releases/download/$Version/$Filename"

$TmpDir = [System.IO.Path]::Combine([System.IO.Path]::GetTempPath(), [System.IO.Path]::GetRandomFileName())
New-Item -ItemType Directory -Path $TmpDir | Out-Null

try {
    Write-Host "Downloading $Binary $Version (windows/$Arch)..."
    $ZipPath = Join-Path $TmpDir $Filename
    Invoke-WebRequest -Uri $Url -OutFile $ZipPath -UseBasicParsing

    Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir | Out-Null
    }

    Copy-Item (Join-Path $TmpDir "${Binary}.exe") (Join-Path $InstallDir "${Binary}.exe") -Force

    # Add to PATH for current user if not already present
    $UserPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        [System.Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
        Write-Host "Added $InstallDir to your PATH."
        Write-Host "Restart your terminal for the PATH change to take effect."
    }

    Write-Host ""
    Write-Host "$Binary $Version installed to $InstallDir\${Binary}.exe"
    Write-Host "Run '$Binary --version' to verify."
} finally {
    Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue
}

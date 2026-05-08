[CmdletBinding()]
param(
    [switch]$Race,
    [switch]$Vet,
    [switch]$Coverage
)

$ErrorActionPreference = "Stop"

function Invoke-Step {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Name,
        [Parameter(Mandatory = $true)]
        [scriptblock]$Step
    )

    Write-Host "==> $Name"
    & $Step
    if ($LASTEXITCODE -ne 0) {
        throw "$Name failed with exit code $LASTEXITCODE"
    }
}

function Resolve-GoCommand {
    $goCommand = Get-Command go -ErrorAction SilentlyContinue
    if ($goCommand) {
        return $goCommand.Source
    }

    $defaultWindowsGo = "C:\Program Files\Go\bin\go.exe"
    if (Test-Path $defaultWindowsGo) {
        return $defaultWindowsGo
    }

    throw "Go toolchain not found. Install Go 1.22 or newer, add it to PATH, or install it at $defaultWindowsGo."
}

$go = Resolve-GoCommand

$goVersion = & $go env GOVERSION
if ($LASTEXITCODE -ne 0) {
    throw "failed to read Go version"
}
Write-Host "Using $goVersion"

Invoke-Step "Download modules" { & $go mod download }

if ($Coverage) {
    Invoke-Step "Run tests with coverage" { & $go test ./... -coverprofile coverage.out }
    Invoke-Step "Show coverage summary" { & $go tool cover -func coverage.out }
} else {
    Invoke-Step "Run tests" { & $go test ./... }
}

if ($Race) {
    Invoke-Step "Run race tests" { & $go test -race ./... }
}

if ($Vet) {
    Invoke-Step "Run go vet" { & $go vet ./... }
}

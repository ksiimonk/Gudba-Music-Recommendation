param(
    [Parameter(Position = 0)]
    [ValidateSet("help", "doctor", "migrate", "reset-db", "migration-status", "test", "backend", "frontend", "frontend-build", "api-smoke", "check")]
    [string]$Command = "help"
)

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$Root = (Resolve-Path (Join-Path $ScriptDir "..")).Path
$BackendDir = Join-Path $Root "backend"
$FrontendDir = Join-Path $Root "frontend"
$env:GOTELEMETRY = "off"
$env:GOTELEMETRYDIR = Join-Path $BackendDir "tmp\telemetry"

function Show-Help {
    Write-Host "Usage: .\scripts\dev.ps1 <command>"
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  doctor            Check local tools and Postgres reachability"
    Write-Host "  migrate           Create database if needed and apply backend migrations + seed"
    Write-Host "  reset-db          Drop all project DB objects, then apply migrations + seed"
    Write-Host "  migration-status  Show applied/pending migrations"
    Write-Host "  test              Run backend tests"
    Write-Host "  backend           Start backend API on HTTP_PORT from backend/.env"
    Write-Host "  frontend          Start Vite dev server"
    Write-Host "  frontend-build    Build frontend"
    Write-Host "  api-smoke         Call health/tracks/playlists endpoints on a running backend"
    Write-Host "  check             Run doctor, tests, frontend build, and reset-db"
}

function Set-GoCache {
    $env:GOTELEMETRY = "off"
    $env:GOTELEMETRYDIR = Join-Path $BackendDir "tmp\telemetry"
    $env:GOCACHE = Join-Path $BackendDir "tmp\gocache"
    $env:GOMODCACHE = Join-Path $BackendDir "tmp\gomodcache"
    New-Item -ItemType Directory -Force -Path $env:GOTELEMETRYDIR | Out-Null
    New-Item -ItemType Directory -Force -Path $env:GOCACHE | Out-Null
    New-Item -ItemType Directory -Force -Path $env:GOMODCACHE | Out-Null
}

function Invoke-BackendGo {
    param([string[]]$GoArgs)

    Set-GoCache
    Push-Location $BackendDir
    try {
        & go @GoArgs
        if ($LASTEXITCODE -ne 0) {
            throw "go $($GoArgs -join ' ') failed with exit code $LASTEXITCODE"
        }
    }
    finally {
        Pop-Location
    }
}

function Invoke-FrontendNpm {
    param([string[]]$NpmArgs)

    Push-Location $FrontendDir
    try {
        & npm @NpmArgs
        if ($LASTEXITCODE -ne 0) {
            throw "npm $($NpmArgs -join ' ') failed with exit code $LASTEXITCODE"
        }
    }
    finally {
        Pop-Location
    }
}

function Get-EnvValue {
    param([string]$Name)

    $envFile = Join-Path $BackendDir ".env"
    if (-not (Test-Path $envFile)) {
        return $null
    }

    $line = Get-Content $envFile | Where-Object { $_ -match "^$Name=" } | Select-Object -First 1
    if (-not $line) {
        return $null
    }

    return $line.Substring($Name.Length + 1).Trim()
}

function Get-BackendBaseUrl {
    $port = Get-EnvValue "HTTP_PORT"
    if (-not $port) {
        $port = "8080"
    }

    return "http://127.0.0.1:$port"
}

function Test-TcpPort {
    param(
        [string]$HostName,
        [int]$Port,
        [int]$TimeoutMs = 1000
    )

    $client = [System.Net.Sockets.TcpClient]::new()
    try {
        $connect = $client.BeginConnect($HostName, $Port, $null, $null)
        if (-not $connect.AsyncWaitHandle.WaitOne($TimeoutMs, $false)) {
            return $false
        }

        $client.EndConnect($connect)
        return $true
    }
    catch {
        return $false
    }
    finally {
        $client.Close()
    }
}

function Invoke-Doctor {
    Write-Host "Project root: $Root"

    try {
        $goVersion = & go version 2>$null
        if ($LASTEXITCODE -eq 0) { Write-Host "Go:   $goVersion" } else { Write-Host "Go:   not found" }
    }
    catch { Write-Host "Go:   not found" }

    try { Write-Host "Node: $(& node --version)" } catch { Write-Host "Node: not found" }
    try { Write-Host "npm:  $(& npm --version)" } catch { Write-Host "npm:  not found" }

    $pgUrl = Get-EnvValue "PG_URL"
    if (-not $pgUrl) {
        Write-Host "Postgres: PG_URL is missing in backend/.env"
        return
    }

    try {
        $uri = [System.Uri]$pgUrl
        $hostName = $uri.Host
        $port = $uri.Port
        $database = $uri.AbsolutePath.TrimStart("/")

        if (Test-TcpPort -HostName $hostName -Port $port) {
            Write-Host "Postgres: reachable at ${hostName}:${port}, database '${database}'"
        }
        else {
            Write-Host "Postgres: not reachable at ${hostName}:${port}"
        }
    }
    catch {
        Write-Host "Postgres: PG_URL could not be parsed"
    }
}

function Invoke-ApiSmoke {
    $baseUrl = Get-BackendBaseUrl
    Write-Host "Checking API at $baseUrl"

    $health = Invoke-RestMethod "$baseUrl/healthz"
    if ($health.message -ne "ok") {
        throw "health check failed"
    }
    Write-Host "healthz: ok"

    $tracks = Invoke-RestMethod "$baseUrl/api/v1/tracks"
    if (-not $tracks.tracks -or $tracks.tracks.Count -lt 1) {
        throw "tracks endpoint returned no demo tracks"
    }
    Write-Host "tracks: $($tracks.tracks.Count)"

    $track = Invoke-RestMethod "$baseUrl/api/v1/tracks/1"
    if (-not $track.track.id) {
        throw "track detail endpoint failed"
    }
    Write-Host "track detail: $($track.track.title)"

    $playlists = Invoke-RestMethod "$baseUrl/api/v1/playlists"
    if (-not $playlists.playlists -or $playlists.playlists.Count -lt 1) {
        throw "playlists endpoint returned no demo playlists"
    }
    Write-Host "playlists: $($playlists.playlists.Count)"

    $playlist = Invoke-RestMethod "$baseUrl/api/v1/playlists/1"
    if (-not $playlist.playlist.id) {
        throw "playlist detail endpoint failed"
    }
    Write-Host "playlist detail: $($playlist.playlist.name)"
}

switch ($Command) {
    "help" {
        Show-Help
    }
    "doctor" {
        Invoke-Doctor
    }
    "migrate" {
        Invoke-BackendGo @("run", "./cmd/migrate", "up")
    }
    "reset-db" {
        Invoke-BackendGo @("run", "./cmd/migrate", "reset")
    }
    "migration-status" {
        Invoke-BackendGo @("run", "./cmd/migrate", "status")
    }
    "test" {
        Invoke-BackendGo @("test", "./cmd/...", "./config", "./internal/...")
    }
    "backend" {
        Invoke-BackendGo @("run", "./cmd/app")
    }
    "frontend" {
        Invoke-FrontendNpm @("run", "dev")
    }
    "frontend-build" {
        Invoke-FrontendNpm @("run", "build")
    }
    "api-smoke" {
        Invoke-ApiSmoke
    }
    "check" {
        Invoke-Doctor
        Invoke-BackendGo @("test", "./cmd/...", "./config", "./internal/...")
        Invoke-FrontendNpm @("run", "build")
        Invoke-BackendGo @("run", "./cmd/migrate", "reset")
        Write-Host ""
        Write-Host "Next:"
        Write-Host "  1. .\scripts\dev.ps1 backend"
        Write-Host "  2. In another terminal: .\scripts\dev.ps1 api-smoke"
    }
}

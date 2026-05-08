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
    Write-Host "  api-smoke         Full user flow smoke test (health, tracks, playlists, auth, onboarding, recs, events, stats)"
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
    Write-Host ""

    function Check-Step {
        param([scriptblock]$Block, [string]$Label)
        try {
            & $Block
            Write-Host "  [OK] $Label"
        }
        catch {
            Write-Host "  [FAIL] $Label : $_"
            throw "api-smoke failed at: $Label"
        }
    }

    # ---- Unauthenticated endpoints ----

    Check-Step -Label "GET /healthz" -Block {
        $health = Invoke-RestMethod "$baseUrl/healthz"
        if ($health.message -ne "ok") { throw "health check failed" }
    }

    Check-Step -Label "GET /api/v1/genres" -Block {
        $genres = Invoke-RestMethod "$baseUrl/api/v1/genres"
        if (-not $genres.genres -or $genres.genres.Count -lt 1) { throw "genres endpoint failed" }
    }

    Check-Step -Label "GET /api/v1/artists" -Block {
        $artists = Invoke-RestMethod "$baseUrl/api/v1/artists"
        if (-not $artists.artists -or $artists.artists.Count -lt 1) { throw "artists endpoint failed" }
    }

    Check-Step -Label "GET /api/v1/tracks" -Block {
        $tracks = Invoke-RestMethod "$baseUrl/api/v1/tracks"
        if (-not $tracks.tracks -or $tracks.tracks.Count -lt 1) { throw "tracks endpoint failed" }
    }

    Check-Step -Label "GET /api/v1/tracks/1" -Block {
        $track = Invoke-RestMethod "$baseUrl/api/v1/tracks/1"
        if (-not $track.track.id) { throw "track detail endpoint failed" }
    }

    Check-Step -Label "GET /api/v1/playlists" -Block {
        $playlists = Invoke-RestMethod "$baseUrl/api/v1/playlists"
        if (-not $playlists.playlists -or $playlists.playlists.Count -lt 1) { throw "playlists endpoint failed" }
    }

    Check-Step -Label "GET /api/v1/playlists/1" -Block {
        $playlist = Invoke-RestMethod "$baseUrl/api/v1/playlists/1"
        if (-not $playlist.playlist.id) { throw "playlist detail endpoint failed" }
    }

    # ---- User registration + auth flow ----

    $email = "smoke$(Get-Random -Min 10000 -Max 99999)@test.local"
    $password = "smoketest123"
    $creds = @{ email = $email; password = $password } | ConvertTo-Json

    Check-Step -Label "POST /api/v1/auth/register" -Block {
        $register = Invoke-RestMethod "$baseUrl/api/v1/auth/register" -Method Post -Body $creds -ContentType "application/json"
        if (-not $register.user.id) { throw "register endpoint failed" }
    }

    Check-Step -Label "POST /api/v1/auth/login" -Block {
        $script:login = Invoke-RestMethod "$baseUrl/api/v1/auth/login" -Method Post -Body $creds -ContentType "application/json"
        if (-not $login.token) { throw "login endpoint failed" }
    }
    $token = $login.token
    $authHeader = @{ Authorization = "Bearer $token" }

    Check-Step -Label "GET /api/v1/me" -Block {
        $me = Invoke-RestMethod "$baseUrl/api/v1/me" -Headers $authHeader
        if ($me.user.email -ne $email) { throw "/me returned wrong email" }
    }

    # ---- Onboarding ----

    $onboarding = @{ genre_ids = @(1, 6); artist_ids = @(1); track_ids = @(); contexts = @("focus") } | ConvertTo-Json

    Check-Step -Label "POST /api/v1/me/onboarding" -Block {
        $onbResult = Invoke-RestMethod "$baseUrl/api/v1/me/onboarding" -Method Post -Body $onboarding -ContentType "application/json" -Headers $authHeader
        if ($onbResult.message -ne "preferences saved") { throw "onboarding endpoint failed" }
    }

    Check-Step -Label "GET /api/v1/me/profile" -Block {
        $profile = Invoke-RestMethod "$baseUrl/api/v1/me/profile" -Headers $authHeader
        if ($profile.profile.favorite_genre_ids.Count -lt 1) { throw "profile endpoint returned no genres" }
    }

    # ---- Recommendations ----

    Check-Step -Label "GET /api/v1/recommendations/tracks" -Block {
        $recs = Invoke-RestMethod "$baseUrl/api/v1/recommendations/tracks?limit=5" -Headers $authHeader
        if (-not $recs.recommendations -or $recs.recommendations.Count -lt 1) { throw "track recommendations returned no results" }
    }

    # ---- Events ----

    $trackId = $recs.recommendations[0].track.id

    Check-Step -Label "POST /api/v1/tracks/$trackId/like" -Block {
        $likeResult = Invoke-RestMethod "$baseUrl/api/v1/tracks/$trackId/like" -Method Post -Headers $authHeader
        if ($likeResult.message -ne "event saved") { throw "track like endpoint failed" }
    }

    Check-Step -Label "POST /api/v1/tracks/$trackId/play" -Block {
        $playResult = Invoke-RestMethod "$baseUrl/api/v1/tracks/$trackId/play" -Method Post -Headers $authHeader
        if ($playResult.message -ne "event saved") { throw "track play endpoint failed" }
    }

    # ---- Playlist recommendations ----

    Check-Step -Label "GET /api/v1/recommendations/playlists" -Block {
        $plRecs = Invoke-RestMethod "$baseUrl/api/v1/recommendations/playlists?limit=3" -Headers $authHeader
    }

    # ---- Admin stats ----

    Check-Step -Label "GET /api/v1/admin/stats" -Block {
        $stats = Invoke-RestMethod "$baseUrl/api/v1/admin/stats" -Headers $authHeader
        if (-not $stats.users_count) { throw "admin stats endpoint failed" }
    }

    Check-Step -Label "GET /api/v1/admin/recommendation-metrics" -Block {
        $metrics = Invoke-RestMethod "$baseUrl/api/v1/admin/recommendation-metrics" -Headers $authHeader
    }

    # ---- Seed user login ----

    Check-Step -Label "Seed user login" -Block {
        $seedCreds = @{ email = "seed@music.local"; password = "plain-password" } | ConvertTo-Json
        $seedLogin = Invoke-RestMethod "$baseUrl/api/v1/auth/login" -Method Post -Body $seedCreds -ContentType "application/json"
        if (-not $seedLogin.token) { throw "seed user login failed" }
    }

    Write-Host ""
    Write-Host "smoke OK: all endpoints passed"
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

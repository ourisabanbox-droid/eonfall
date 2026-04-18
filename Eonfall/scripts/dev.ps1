param(
    [Parameter(Mandatory = $true)]
    [ValidateSet("up", "down", "logs", "ps", "migrate-up", "migrate-down", "seed", "run", "tidy", "fmt", "test", "build")]
    [string]$Command
)

$DbUrl = "postgres://postgres:J3rusa%2F3m@localhost:5433/eonfall?sslmode=disable"
$DockerDbUrl = "postgres://postgres:J3rusa%2F3m@host.docker.internal:5433/eonfall?sslmode=disable"
$ProjectPath = (Get-Location).Path
$env:DATABASE_URL = $DbUrl
$env:REDIS_URL = "redis://localhost:6379"
$env:HTTP_PORT = "8080"
$env:TICK_RATE_MS = "1000"

switch ($Command) {
    "up" {
        docker compose up -d
    }
    "down" {
        docker compose down
    }
    "logs" {
        docker compose logs -f
    }
    "ps" {
        docker compose ps
    }
    "migrate-up" {
        docker run --rm `
          -v "${ProjectPath}:/work" `
          -w /work `
          migrate/migrate `
          -path ./migrations `
          -database $DockerDbUrl `
          up
    }
    "migrate-down" {
        docker run --rm `
          -v "${ProjectPath}:/work" `
          -w /work `
          migrate/migrate `
          -path ./migrations `
          -database $DockerDbUrl `
          down 1
    }
    "seed" {
        Get-Content .\seed_dev.sql | docker exec -i eonfall-postgres psql -U postgres -d eonfall
    }
    "run" {
        go run .\cmd\server
    }
    "tidy" {
        go mod tidy
    }
    "fmt" {
        go fmt .\...
    }
    "test" {
        go test .\...
    }
    "build" {
        New-Item -ItemType Directory -Force -Path .\bin | Out-Null
        go build -o .\bin\project-eonfall.exe .\cmd\server
    }
}
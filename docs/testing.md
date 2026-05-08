# Go Testing

This repository targets Go 1.22 or newer.

## Local test commands

On Windows PowerShell:

```powershell
.\scripts\test.ps1
```

Optional checks:

```powershell
.\scripts\test.ps1 -Vet
.\scripts\test.ps1 -Race
.\scripts\test.ps1 -Coverage
```

On shells with `make`:

```sh
make test
make test-vet
make test-race
make test-cover
```

## Continuous integration

GitHub Actions runs `go mod download`, `go test ./...`, and `go vet ./...` on Ubuntu and Windows for pushes to `main` or `master` and for pull requests.

The smoke test in `internal/cli/root_test.go` verifies that the root CLI command still exposes the expected command surface without requiring an OpenRouter API key or user runtime data.

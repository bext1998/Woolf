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

## Focused smoke coverage

The smoke test in `internal/cli/root_test.go` verifies that the root CLI command still exposes the expected command surface without requiring an OpenRouter API key or user runtime data.

The `agents` command tests cover listing and showing built-in roles, adding a custom YAML role into the configured agents directory, loading it through the registry, and deleting it again.

The `start` command tests use a fake chat client so the orchestration path can be exercised without calling OpenRouter. These tests cover both the successful session path and the error path where agent/API failures must surface as command errors.

The config tests verify environment overrides for `OPENROUTER_API_KEY` and `WOOLF_SESSIONS_DIR`, plus API key masking behavior used by user-facing configuration output.

The OpenRouter client tests verify SSE parsing, missing API key handling, HTTP status to Woolf error-code mapping, `Retry-After` handling for 429 responses, and bounded retry behavior for 5xx responses. These tests use local HTTP test servers and do not call OpenRouter.

The orchestrator tests verify session persistence, cancellation behavior, stream error handling, and context propagation between agents. Stream errors should persist the skipped response, mark the session as `error`, and emit an error event instead of being reported as a completed run. Context builder tests verify that draft content, session summaries, user interventions, focus ranges, role prompt metadata, stance tags, and previous agent responses are included in the messages sent to the chat client.

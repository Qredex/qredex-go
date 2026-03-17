# Contributing to the Qredex Go SDK

Thank you for your interest in contributing. This document describes how to contribute code, tests, and documentation to the Qredex Go SDK.

---

## Before You Start

- This SDK covers the **Qredex Integrations API only**.  Do not add support for Merchant API (`/api/v1/merchant/**`) or Internal API (`/api/v1/internal/**`) endpoints.
- Changes to the public API surface are **product decisions**.  Open an issue before beginning work on any new exported type, function, or behaviour.
- The canonical flow (IIT → PIT → paid order → refund) must remain correct after every change.

---

## Development Requirements

- Go 1.21 or later
- No external test frameworks — use the standard library `testing` package only
- No external runtime dependencies — the SDK uses the standard library only

---

## Getting the Code

```bash
git clone https://github.com/qredex/sdk-go.git
cd sdk-go
go mod download
```

---

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- All exported types and functions must have a GoDoc comment
- Use `context.Context` as the first argument on every networked operation
- Prefer explicit over implicit; avoid hidden global state
- Do not log secrets, tokens, or Authorization header values

---

## Running Tests

```bash
# Unit tests
go test ./...

# With race detector (required before submitting a PR)
go test -race ./...

# Formatting check
go fmt ./...
git diff --exit-code

# Vet
go vet ./...
```

Live integration tests are excluded by default. See [README.md](README.md) for how to run them.

---

## Test Requirements

- Every bug fix must include a regression test that reproduces the exact failure scenario.
- New behaviour must include both a happy-path and at least one failure-path test.
- Use `FakeTransport` (defined in `fake_transport_test.go`) to mock HTTP responses.  Do not mock `net/http` internals directly.
- Tests that call the real Qredex API must use the `live` build tag and be excluded from the default test run.

---

## Submitting a Pull Request

1. Fork the repository and create a branch: `git checkout -b fix/my-fix`
2. Make your changes on that branch.
3. Ensure all tests pass with the race detector enabled.
4. Ensure `go fmt ./...` produces no diff.
5. Update `CHANGELOG.md` under `[Unreleased]` with a brief description of your change.
6. Open a pull request against `main`.

PR descriptions should explain:
- **What** changed and **why**
- Any public API additions or removals
- How to test the change manually if applicable

---

## Changelog Format

This project follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

Add your entry under the `[Unreleased]` heading using one of these categories:

- **Added** — new features or endpoints
- **Changed** — changes to existing behaviour
- **Fixed** — bug fixes
- **Deprecated** — features that will be removed in a future version
- **Removed** — features that have been removed
- **Security** — security-related improvements

---

## Naming Rules

Preserve canonical Qredex terminology in all code, comments, and documentation:

| Term | Do NOT rename to |
|------|-----------------|
| IIT (Influence Intent Token) | `session_token`, `click_token` |
| PIT (Purchase Intent Token) | `cart_token`, `checkout_token` |
| `token_integrity` | `token_valid`, `integrity_check` |
| `integrity_reason` | `failure_reason`, `check_reason` |
| `resolution_status` | `attribution_status`, `status` |
| OrderAttribution | `OrderRecord`, `AttributionRecord` |

---

## Licence

By contributing you agree that your contributions will be licensed under the Apache License 2.0, the same licence used by this project.

---

## Contact

- **Email**: os@qredex.com
- **Website**: https://qredex.com

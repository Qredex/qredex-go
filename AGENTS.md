<!--
     ▄▄▄▄
   ▄█▀▀███▄▄              █▄
   ██    ██ ▄             ██
   ██    ██ ████▄▄█▀█▄ ▄████ ▄█▀█▄▀██ ██▀
   ██  ▄ ██ ██   ██▄█▀ ██ ██ ██▄█▀  ███
    ▀█████▄▄█▀  ▄▀█▄▄▄▄█▀███▄▀█▄▄▄▄██ ██▄
         ▀█

   Copyright (C) 2026 — 2026, Qredex, LTD. All Rights Reserved.

   DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.

   Licensed under the Apache License, Version 2.0. See LICENSE for the full license text.
   You may not use this file except in compliance with that License.
   Unless required by applicable law or agreed to in writing, software distributed under the
   License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
   either express or implied. See the License for the specific language governing permissions
   and limitations under the License.

   If you need additional information or have any questions, please email: copyright@qredex.com
-->

## Purpose

This document defines how any AI agent (or engineer acting as an agent) must work inside a Qredex PHP SDK repository to avoid drift, regressions, and "helpful but wrong" changes.

This SDK is not a generic HTTP wrapper.
It is a Qredex product surface and must remain aligned with the platform contract.

## Engineering Standards

Engineer this like a serious public infrastructure SDK.

### Standards

- Optimize for developer trust, safety, and long-term maintainability.
- Prefer fewer, stronger primitives over wide surface area.
- Make the SDK easy to use correctly and hard to misuse.
- Hide raw HTTP/auth/plumbing where appropriate, but never hide important behavior.
- Keep the public API small, explicit, typed, and predictable.
- Favor clean names, immutable value objects, stable contracts, and framework-neutral design.
- Prevent footguns: safe defaults, explicit config, clear errors, deterministic behavior, idempotency support, timeout/retry discipline, and no leaky internals.
- Do more with less: remove anything redundant, speculative, or low-value.
- Treat DX as part of the product: short happy path, strong docs, clean examples, good package metadata, and professional release/testing standards.
- If a design choice improves durability, readability, and correct usage, prefer it over cleverness or abstraction for its own sake.
- Build this to feel like a platform-grade SDK a company would trust in production.

## Efficiency, Quality, and Infrastructure Discipline

Use the minimum context, tokens, tool calls, edits, and validation needed to complete the task correctly.

### Working rules

- **Read narrowly first.** Expand only when needed for correctness.
- **Edit narrowly, but completely.** Include every directly connected change required for the result to be correct.
- **Validate with the lightest check that gives real confidence.**
- **Do not scan the whole codebase** unless the task truly requires it.
- **Do not perform broad refactors, broad searches, speculative cleanup, or optional exploration** unless requested or clearly necessary.
- **Do not invent new flows, abstractions, endpoints, or patterns** if the existing architecture already supports the task.
- **Reuse existing code paths, commands, adapters, and conventions** wherever possible.
- **Keep responses short, direct, and action-focused.**

### Quality guardrails

- **Accuracy is mandatory.**
- **Completeness matters more than superficial minimalism.**
- **Minimal work does not mean shallow work.**
- **If a wider check is required for safety, correctness, or integration integrity, do it — but keep it tightly scoped.**
- **If a requested change likely affects adjacent logic, inspect the smallest necessary connected surface before editing.**
- **Make the narrowest correct change, not the fastest careless change.**

### Qredex guardrails

- **Preserve determinism, idempotency, zoning, tenant scoping, and store scoping.**
- **Keep changes layer-correct and aligned with the existing architecture.**
- **Prefer canonical flows over parallel implementations.**
- **Avoid duplicate logic, fragmented behavior, and unnecessary abstractions.**

### Infrastructure and platform judgment

- **Act like a senior infrastructure/platform engineer, not a code generator.**
- **Recommend the most durable, secure, operationally safe, and platform-aligned path** when it is materially better than the requested implementation.
- **Favor standardization, observability, deterministic behavior, contract clarity, and clean boundaries** over clever shortcuts.
- **Call out drift, weak boundaries, duplicated responsibility, leaky abstractions, and anything that undermines Qredex as a platform.**
- **Treat naming, packaging, SDK boundaries, auth surfaces, API shape, and execution flow as strategic platform decisions, not local implementation details.**
- **When several options are viable, recommend the one that best improves long-term reliability, maintainability, developer experience, and platform leverage.**
- **Proactively surface high-value improvements, risks, and next best steps without waiting to be asked.**

### Model usage limit discipline

- **Treat model and agent usage limits as a hard engineering constraint.**
- **Optimize for minimum usage without degrading correctness, safety, or architectural quality.**
- **Stop exploring once sufficient evidence exists.** Do not keep reading or probing after the safe implementation path is already clear.
- **Use the fewest files, shortest useful command output, and narrowest validation** that still provides real confidence.
- **Avoid speculative work.** Do not expand scope unless the task or integration risk requires it.
- **Keep communication compressed, direct, and high-signal.**
- **Escalate only when necessary.** If materially more usage would be required to increase certainty, state the trade-off briefly before doing broader exploration or heavier validation.

## System Overview

This repository builds the Qredex Go server SDK.

### Product purpose

The Go SDK exists to make Qredex machine-to-machine integrations easy to adopt and hard to misuse.

It should help merchant backends and partner platforms:

- authenticate cleanly
- create and read creators
- create and read links
- issue IIT where backend issuance is appropriate
- lock PIT where authenticated machine lock is appropriate
- record paid orders and refunds
- read order attribution details for operational workflows

### What this SDK is not

This SDK must not become:

- a Merchant API SDK
- an Internal API SDK
- a browser/runtime agent
- a Shopify embedded/session helper
- a full mirror of every Qredex endpoint

## Source of Truth

The authoritative source of truth for platform behavior is the supplied Qredex Core contract material for the task.

That typically includes:

- OpenAPI
- SDK strategy / blueprint material
- auth and zoning rules
- attribution and ingestion model documentation
- completed reference SDKs in other languages when explicitly provided

If local SDK code, local examples, or local docs conflict with supplied Qredex platform docs, the supplied platform docs win.

Do not redefine platform behavior from the SDK repo.
If a required platform rule is unclear or missing, stop and surface the gap instead of inventing it.

## Non-negotiables

- **Do not invent flows.** If it is not in the canonical spec, stop and ask.
- **The public SDK surface is curated.** Generated transport may exist, but it must not become the primary public API.
- **Integrations API only.** No Merchant or Internal API drift.
- **Canonical naming is mandatory.** Preserve IIT, PIT, Order Attribution, `token_integrity`, `integrity_reason`, `resolution_status`.
- **Idempotency is mandatory.** Do not encourage write retry behavior that violates the platform contract.
- **Don't leak secrets.** Never log raw tokens, client secrets, IITs, PITs, or Authorization headers by default.
- **Use constants, not magic strings.** Error codes, header names, event names, and fixed values should not be repeated ad hoc.
- **Test the full flow.** When fixing bugs, add tests that reproduce the exact failure scenario.
- **Never create unused code.** Do not add classes, methods, interfaces, configuration, imports, or helpers that are not actually used.
- **Clean up after yourself.** If you create dead code during refactors, remove it immediately.

## Review Decision Protocol (Mandatory)

- **Classify every claim with one verdict only:** `VALID`, `MISPLACED_LAYER`, `INVALID`, or `UNVERIFIED`.
- **No verdict without evidence:** always cite the relevant local code or supplied contract source.
- **Do not implement before verdict:** first prove the claim, then patch.
- **If `MISPLACED_LAYER`, name the correct layer** (`public-api`, `resource`, `transport`, `auth`, `error-model`, `examples`, `docs`, `tests`).
- **Every behavior-changing fix must include tests that reproduce the failure scenario.**
- **Do not close work without running the repository’s canonical validation commands and reporting result.**

## Strategic Thinking & Challenge Protocol (Mandatory)

- **The agent must not default to agreement.**
- **If a proposed design has a cleaner, safer, or more scalable alternative, it must be presented.**
- **If a decision increases long-term complexity, the agent must explicitly call it out.**
- **If multiple valid approaches exist, the agent must:**
  - **Present the top 2 options**
  - **State trade-offs clearly**
  - **Recommend one with justification.**
- **If the user's idea is optimal, the agent must explicitly explain why it is optimal and why alternatives are inferior.**
- **The agent must prioritize architectural integrity over short-term convenience.**

### Decision evaluation criteria

When evaluating any design choice, assess:

- **Platform contract integrity** - Does this stay true to Qredex API behavior?
- **Authentication boundary integrity** - Does this respect the Integrations-only scope?
- **Layer purity** - Does this belong in public API, resource layer, transport, auth, or examples?
- **Long-term maintainability** - Will this make future SDKs easier or harder?
- **Blast radius of future changes** - How many consumers break if this changes?
- **Operational clarity** - Can SDK users debug this safely and clearly?

If a proposal weakens any of the above, the risk must be stated before implementation.

## Plan Mode

Use plan mode whenever work is more than 3 steps or touches architecture.

- Write checklist tasks.
- Identify risks.
- Define acceptance criteria.
- If something goes sideways, stop and re-plan.

## Change Discipline

- Prefer minimal, reversible changes.
- One change set = one theme when possible.
- Do not widen scope silently.
- Public API additions and changes are product decisions, not casual local edits.


## Naming Rules

- Canonical terms:
  - **IIT** = Influence Intent Token
  - **PIT** = Purchase Intent Token
  - **OrderAttribution** = resolved outcome from paid/refund events
- Canonical field names from the API must be preserved where they are part of the contract.
- Do not rename Qredex domain language into softer generic language.
- Keep public product identity as **Qredex**.

## Public API Rules

- The main public entrypoint should be simple and intentional.
- Support both:
  - environment/bootstrap path for convenience
  - explicit initialization/configuration path
- Group operations by resource, not by transport mechanics.
- Do not expose raw HTTP client details as the primary usage pattern.
- Do not flatten away important Qredex decision fields.
- Prefer explicit request input over long positional argument lists.

### Bad direction

- giant static helper classes
- hidden global state
- public transport-first API
- convenience abstractions that obscure Qredex contract facts

## Security Rules

- Never accept Merchant or Internal auth surfaces into this SDK.
- This SDK is for Integrations auth only.
- Token handling must be automatic by default, but explicit in configuration and observable in behavior.
- Redact secrets in logs and debug output.
- Do not build convenience helpers that weaken auth boundaries.
- Retry rules must stay safe:
  - reads may be retried if configured carefully
  - writes must not be blindly retried
  - prefer idempotent external IDs and platform-safe replay patterns



### SDK design expectations

- Prefer clear types and value objects over generic interfaces.
- Avoid massive interface hierarchies.
- Avoid framework-specific coupling unless explicitly requested.
- Support normal Go backend usage first (net/http, context patterns).
- If accepting request input as structs, validate and normalize consistently.
- If introducing abstractions, use them because they improve correctness, not because they look formal.
- Use Go's error handling idioms (error returns, error wrapping with fmt.Errorf, errors.Is/As).
- Prefer receiver methods on types over package-level functions when the operation logically belongs to the type.

### HTTP / transport expectations

- Keep transport swappable only if there is a real benefit.
- Do not over-abstract the HTTP layer.
- Make timeouts, retries, base URL/environment, and auth observable/configurable.
- Preserve raw API error details in typed error responses.
- Use http.RoundTripper for test doubles to enable dependency injection.

## Testing Guidelines

### Required coverage

At minimum cover:

- token issuance (OAuth flow)
- create/get/list creator
- create/get/list link
- link stats where in scope
- issue IIT
- lock PIT
- paid order submission
- order list / details where in scope
- refund submission
- failure paths for auth, validation, conflict, and network behavior

### Test types

- **Unit tests** (`*_test.go`) for public API behavior, auth behavior, error parsing, and serialization
- **Contract tests** for transport/request mapping against fixture responses or mock servers
- **Live tests** only as explicit opt-in, never as a default dependency for local or CI success

### Test rules

- Fixes require regression tests.
- Do not depend on external environments unless explicitly configured via `t.Setenv()` or env var overrides.
- Run `go test ./...`, `go vet ./...`, and `golangci-lint run ./...` as part of validation.
- Report exactly which validation commands were run and whether they passed.

### Test patterns in this repo

- **`FakeTransport`** (`fake_transport_test.go`) is the canonical test double. It implements `http.RoundTripper`, accepts queued `TransportResponse` or error via `Push()`, and records sent `TransportRequest` objects on `Requests()`. Always use this instead of mocking the http package directly.
- **`qredex_test.go`** — core SDK behavior: `Bootstrap()`, `New()`, config validation, resource access, error mapping.
- **`*_test.go` files** — integration and contract tests using `FakeTransport` to verify request serialization and response deserialization.
- **Error type tests** — HTTP status → typed error mapping (`APIError`, `AuthenticationError`, `AuthorizationError`, `NotFoundError`, `ConflictError`, `RateLimitError`, `ValidationError`).
- **Model tests** — `UnmarshalJSON` deserialization and `MarshalJSON` serialization for all model types.
- **Request validation tests** — validation rules for all write operations (required fields, UUID format, amount bounds).
- **Example tests** (`example_test.go`) — demonstrate canonical usage patterns with `Example_*` functions.
- To use `FakeTransport` in tests: create a `Config`, set `HTTPClient` to one wrapping `FakeTransport`, then push OAuth token response first, then resource responses.

### Environment variables

The SDK reads these env vars via `Bootstrap()` or can be set explicitly in `Config`:

| Variable | Required | Description |
|---|---|---|
| `QREDEX_CLIENT_ID` | **Yes** | OAuth client ID for Integrations auth |
| `QREDEX_CLIENT_SECRET` | **Yes** | OAuth client secret |
| `QREDEX_SCOPE` | No | Space-separated OAuth scopes (e.g., `direct:api direct:creators:write`) |
| `QREDEX_ENVIRONMENT` | No | `production` (default), `staging`, or `development` |
| `QREDEX_BASE_URL` | No | Override base URL (bypasses environment resolution) |
| `QREDEX_TIMEOUT_MS` | No | HTTP timeout in milliseconds (default: `10000`) |

### Scope values

Scope is a string constant representing OAuth scopes. Canonical values:

| Constant | Value |
|---|---|
| `ScopeAPI` | `direct:api` |
| `ScopeLinksRead` | `direct:links:read` |
| `ScopeLinksWrite` | `direct:links:write` |
| `ScopeCreatorsRead` | `direct:creators:read` |
| `ScopeCreatorsWrite` | `direct:creators:write` |
| `ScopeOrdersRead` | `direct:orders:read` |
| `ScopeOrdersWrite` | `direct:orders:write` |
| `ScopeIntentsRead` | `direct:intents:read` |
| `ScopeIntentsWrite` | `direct:intents:write` |

### Environment values

Environment is a string constant representing the Qredex API environment:

| Constant | Base URL |
|---|---|
| `Production` | `https://api.qredex.com` |
| `Staging` | `https://staging-api.qredex.com` |
| `Development` | `http://localhost:8080` |

### CI/CD workflows

GitHub Actions workflows in `.github/workflows/`:

- **`ci.yml`** — runs on PR/push: `go build`, `go test -race`, `go vet`, `golangci-lint`, `go fmt`
- **`create-release-tag-on-version-change.yml`** — creates the semantic-version tag when `SDKVersion` changes on `main`
- **`release.yml`** — validates the tagged release and creates the GitHub Release from semantic-version tags (`v*`)

## Documentation Rules

The repo should include clear package docs, typically:

- `README.md`
- `INTEGRATION_GUIDE.md`
- `ERRORS.md`
- examples for the canonical flow

Docs must explain:

- what the SDK is for
- what it is not for
- canonical IIT -> PIT -> paid/refund flow
- auth behavior
- common mistakes to avoid
- environment/bootstrap vs explicit init usage

Examples are part of the SDK contract.
They must reflect real canonical usage, not shortcuts or fake flows.

### README contact section

The bottom of `README.md` must contain a `Qredex Contact` section with exactly:

- Website: `https://qredex.com`
- X: `https://x.com/qredex`
- Email: `os@qredex.com`

Do not omit any of these contact points. If the section is missing or incomplete, add or fix it.

## Release and Changelog Instructions

When preparing a release or making release-related documentation updates, you must follow these rules:

- Read `docs/RELEASING.md` first. Follow the documented process exactly.
- Do not invent a new release process.
- The release source of truth is `SDKVersion` constant in `config.go`. Update it manually for each release.
- Review git history (`git log`) before writing or updating changelog entries: `git log --oneline <lastTag>..HEAD`
- `CHANGELOG.md` must exist. If it does not, create it immediately.
- All changelog entries must be grounded in actual repository history and real implemented changes.
- Do not invent versions, dates, entries, or release notes.
- If the release guide conflicts with the repository state, stop and surface the conflict instead of guessing.
- Run `go build ./...`, `go test ./...`, `go vet ./...`, and `golangci-lint run ./...` as the pre-release verification step.
- Releases are semantic-version Git tags such as `v0.2.0`.
- Preferred release path: update `SDKVersion`, move the matching changelog entries out of `Unreleased`, and push to `main`.
- GitHub Actions creates the tag from the version bump, then validates the tagged release and creates a GitHub Release.
- Go consumers should install tagged releases: `go get github.com/Qredex/qredex-go@v0.2.0`

## Before Starting Work

1. Read the smallest relevant portion of the supplied platform contract.
2. Check for an existing pattern in the SDK repo first.
3. Decide whether the change belongs in public API, transport, auth, errors, examples, or docs.
4. Identify adjacent behavior that must be inspected for correctness.

## Implementation Guidelines

1. Follow existing patterns unless they are clearly wrong.
2. Prefer additive changes in V1.
3. Keep public exports minimal and intentional.
4. Update tests and docs with behavior changes.
5. Run the repository’s canonical validation commands before closing work.
6. Do not leave unused code, stale examples, or dead exports behind.

## Code Review Checklist

- [ ] Public API stays Integrations-only
- [ ] Canonical Qredex naming is preserved
- [ ] Auth boundaries are respected
- [ ] Errors preserve status + `error_code` + debug identifiers where available
- [ ] No secrets are logged
- [ ] Tests cover happy path and failure path
- [ ] Examples still reflect canonical usage
- [ ] Docs are updated if public behavior changed
- [ ] No unused code remains
- [ ] The solution is the narrowest correct one

## Common Pitfalls

### 1. SDK drifting into Merchant/Internal scope
- Problem: convenience pressure widens scope beyond the intended product boundary
- Result: auth confusion, support burden, platform drift
- Solution: keep the SDK strictly Integrations-only unless product scope explicitly changes

### 2. Hiding canonical fields behind convenience booleans
- Problem: important Qredex contract fields get flattened away
- Result: consumers lose explainability and decision detail
- Solution: preserve meaningful fields like `resolution_status`, `token_integrity`, and `integrity_reason`

### 3. Over-abstracting HTTP too early
- Problem: the package becomes harder to reason about than the API itself
- Result: maintainability drops and debugging gets worse
- Solution: keep transport simple and subordinate to the public SDK surface

### 4. Magic strings everywhere
- Problem: error codes, event names, and header names drift
- Result: hard-to-maintain code and inconsistent behavior
- Solution: centralize constants where they improve correctness

### 5. Silent breaking changes
- Problem: method names, signatures, or return shapes change without warning
- Result: downstream consumers break unexpectedly
- Solution: treat public exports as a contract and document migration explicitly

## Definition of Done (Mandatory)

Before marking any task complete, all of the following must be satisfied:

### Testing requirements
- [ ] Relevant tests pass locally
- [ ] Bug fix includes a regression test
- [ ] Behavior change includes happy path and failure path coverage
- [ ] Repository standard validation commands were run and reported

### Documentation requirements
- [ ] Public docs updated if public behavior changed
- [ ] Examples updated if usage changed
- [ ] Any breaking or additive public API change is clearly explained

### Code quality requirements
- [ ] No new lint/static-analysis issues in the changed surface
- [ ] No unused code remains
- [ ] Public API remains curated and coherent

### Security requirements
- [ ] No secret logging
- [ ] Auth scope is still Integrations-only
- [ ] Retry and token behavior remain platform-safe

### Breaking change policy
- [ ] Breaking changes are explicitly called out
- [ ] Backward-compatible approach considered first
- [ ] Migration path documented if breaking change is unavoidable

## No Silent Breaking Changes (Mandatory)

Any breaking change must explicitly state:

1. **What breaks**
2. **Who is affected**
3. **Migration path**
4. **Deprecation timeline** where applicable

Preferred approach:

- deprecate first
- document clearly
- remove only after a deliberate compatibility window

## Final Reminder

This PHP SDK is part of the Qredex platform distribution layer.

Optimize for:

- platform consistency
- integration correctness
- merchant adoption speed
- long-term maintainability
- clear operational behavior

Not for:

- exposing every endpoint
- clever abstractions
- framework lock-in
- convenience that weakens the canonical Qredex model
---

## ⚠️ CRITICAL: Copyright Notice

**ALL created files MUST include the official Qredex Apache-2.0 header used in this repository:**

```go
//    ▄▄▄▄
//  ▄█▀▀███▄▄              █▄
//  ██    ██ ▄             ██
//  ██    ██ ████▄▄█▀█▄ ▄████ ▄█▀█▄▀██ ██▀
//  ██  ▄ ██ ██   ██▄█▀ ██ ██ ██▄█▀  ███
//   ▀█████▄▄█▀  ▄▀█▄▄▄▄█▀███▄▀█▄▄▄▄██ ██▄
//        ▀█
//
// Copyright (C) 2026 — 2026, Qredex, LTD. All Rights Reserved.
//
// DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.
//
// Licensed under the Apache License, Version 2.0. See LICENSE for the full license text.
// You may not use this file except in compliance with that License.
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.
//
// If you need additional information or have any questions, please email: copyright@qredex.com
```

**This applies to:** `.go`, `.md`, `.yaml`, `.yml`, `.json` - ALL files.

**Note:** If you create a new file, add this header at the top. If you modify an existing file without the header, add it.

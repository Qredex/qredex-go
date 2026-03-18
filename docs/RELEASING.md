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

# Releasing

## Source of Truth

- `config.go` `SDKVersion`
- `CHANGELOG.md`
- Git tag `vX.Y.Z`

## Pre-Release Checklist

1. Review repository history since the last tag.
2. Update `CHANGELOG.md`.
3. Update `SDKVersion` in `config.go`.
4. Run:

```bash
go build ./...
go test ./...
go vet ./...
golangci-lint run ./...
```

5. Verify examples and docs still match the current API.

## Tagging

Preferred path:

1. Update `SDKVersion` in `config.go`
2. Move the matching entries from `Unreleased` to `## [X.Y.Z] - YYYY-MM-DD` in `CHANGELOG.md`
3. Push the version bump to `main`

The `create-release-tag-from-version-bump` workflow will create and push `vX.Y.Z` automatically when `SDKVersion` changes on `main`.

Manual fallback:

```bash
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

GitHub Actions publishes the tagged release.

## Rules

- Do not invent release notes.
- Do not skip validation.
- Do not publish a tag whose changelog does not match the repo state.
- Do not treat SDKVersion and the Git tag as separate version sources.

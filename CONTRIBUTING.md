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

# Contributing

## Scope Rules

- Keep the SDK integrations-only.
- Preserve canonical Qredex naming.
- Prefer fewer, stronger public primitives.
- Do not add browser, Merchant API, or Internal API helpers.

## Required Validation

Run all of these before opening a PR:

```bash
go build ./...
go test ./...
go vet ./...
golangci-lint run ./...
```

## Change Rules

- Bug fixes need regression tests.
- Public API changes need docs updates.
- Never log secrets, bearer tokens, IITs, or PITs.
- Keep request validation and error typing intact.

## Release Hygiene

- Update `CHANGELOG.md` when behavior changes.
- Keep `SDKVersion` and the release tag aligned.
- Preferred release path: bump `SDKVersion`, move `CHANGELOG.md` out of `Unreleased`, and push to `main` so the tag is created automatically.

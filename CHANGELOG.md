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

# Changelog

All notable changes to the Qredex Go SDK are documented here.

The format is based on Keep a Changelog and this project follows Semantic Versioning.

## [Unreleased]

## [0.2.1] - 2026-03-31

### Added

- Local request validation with `RequestValidationError`
- `QREDEX_TIMEOUT_MS` support in `Bootstrap()`
- API reference documentation in `docs/API_REFERENCE.md`

### Changed

- Sanitized observability output now uses secret-safe paths
- 401 responses now clear the token cache and retry once with a fresh token
- Public docs, examples, CI, and release metadata were aligned with the actual Go SDK
- `GetLatestUnlocked` is now documented as a deprecated recovery helper

### Fixed

- Test suite no longer imports an undeclared test dependency
- Examples no longer print raw IIT or PIT token values
- Release workflow now publishes the correct Go module path

## [0.2.0] - 2026-03-18

### Fixed

- GET requests now serialize filter parameters as URL query parameters instead of a JSON request body
- `Config.UserAgentSuffix` is now propagated to resource requests

### Changed

- Read operations now retry on HTTP 429 and 5xx responses when `Config.RetryMax > 0`
- `Content-Type: application/json` is no longer set on GET or HEAD requests

### Added

- Example programs for the canonical operations
- `CONTRIBUTING.md`
- `SECURITY.md`
- Additional test coverage for GET query encoding, retries, and user-agent propagation

## [0.1.0] - 2026-01-17

### Added

- Initial Go SDK release

# Changelog

All notable changes to the Qredex Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2026-03-18

### Fixed
- GET requests (list operations) now correctly serialize filter parameters as URL query parameters instead of a JSON request body.
- `Config.UserAgentSuffix` is now propagated to all outgoing resource requests.  Previously it was set only on token requests.

### Changed
- Read operations (GET/HEAD) now retry on HTTP 429 and 5xx responses in addition to network-level failures when `Config.RetryMax > 0`.  The `Retry-After` header is honoured on 429 responses.
- `Content-Type: application/json` is no longer set on GET/HEAD requests.

### Added
- `isRetryableStatus` internal helper (429 and 5xx).
- `structToQueryParams` internal helper — reflection-based struct-to-query-string encoder used by all list operations.
- Example programs for all canonical operations: `create-creator`, `create-link`, `issue-iit`, `lock-pit`, `record-order`, `list-orders`, `record-refund`.
- `CONTRIBUTING.md` with development, testing, naming, and PR guidelines.
- `SECURITY.md` with vulnerability reporting and responsible-disclosure policy.
- Additional test coverage: GET query parameter serialisation, 5xx retry, `UserAgentSuffix` propagation, network error handling, `structToQueryParams` corner cases.

## [0.1.0] - 2026-01-17

### Initial Release

First production-ready release of the Qredex Go SDK.

#### Features
- **Authentication**: OAuth 2.0 client credentials flow with automatic token caching
- **Configuration**: Environment-based (`Bootstrap()`) and explicit (`New()`) initialization
- **Resources**:
  - `Creators().Create()`, `Creators().Get()`, `Creators().List()`
  - `Links().Create()`, `Links().Get()`, `Links().List()`, `Links().GetStats()`
  - `Intents().IssueInfluenceIntentToken()`, `Intents().LockPurchaseIntent()`, `Intents().GetPurchaseIntent()`, `Intents().GetLatestUnlocked()`
  - `Orders().RecordPaidOrder()`, `Orders().List()`, `Orders().GetDetails()`
  - `Refunds().RecordRefund()`
- **Error Handling**: Typed errors (`AuthenticationError`, `AuthorizationError`, `ValidationError`, `NotFoundError`, `ConflictError`, `RateLimitError`, `NetworkError`)
- **Testing**: Comprehensive test suite with `FakeTransport` for mocking HTTP responses
- **Documentation**: README, Integration Guide, Error Handling Guide, GoDoc examples

#### Technical Details
- Minimum Go version: 1.21
- No external dependencies (standard library only)
- Thread-safe token caching
- Context-aware API methods
- Configurable HTTP client, timeout, retry behavior

#### Known Issues
- None

#### Migration Notes
- Initial release — no migration required

---

## Legend
- `[Unreleased]` — Changes not yet released
- `[Added]` — New features
- `[Changed]` — Changes in existing functionality
- `[Deprecated]` — Soon-to-be removed features
- `[Removed]` — Removed features
- `[Fixed]` — Bug fixes
- `[Security]` — Security improvements

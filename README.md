# `@qredex/go`

[![Go Reference](https://pkg.go.dev/badge/github.com/Qredex/qredex-go.svg)](https://pkg.go.dev/github.com/Qredex/qredex-go)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/Qredex/qredex-go)](https://goreportcard.com/report/github.com/Qredex/qredex-go)

The official Qredex Go server SDK for machine-to-machine integrations with the Qredex Integrations API.

## What This SDK Is

This SDK covers the **Qredex Integrations API** for merchant backends and partner platforms:

- **Creators** — Create and manage creator accounts
- **Links** — Create and track influence links
- **Intents** — Issue Influence Intent Tokens (IIT) and lock Purchase Intent Tokens (PIT)
- **Orders** — Record paid orders and read attribution results
- **Refunds** — Record order refunds

Authentication is handled automatically via OAuth 2.0 client credentials flow.

## What This SDK Is NOT

This SDK does **not** cover:

- Merchant API (`/api/v1/merchant/**`) — human dashboard operations
- Internal API (`/api/v1/internal/**`) — Qredex admin operations
- Browser/runtime agent logic — use `@qredex/agent` for client-side operations
- Shopify OAuth/session exchange — use Shopify-specific flows
- Webhook receiver frameworks

## Requirements

- Go 1.21 or later
- Machine-to-machine OAuth credentials from Qredex:
  - Client ID
  - Client Secret

## Installation

```bash
go get github.com/Qredex/qredex-go
```

## Quick Start

### 1. Initialize the SDK

**From environment variables (recommended):**

```go
package main

import (
    "log"
    "github.com/Qredex/qredex-go"
)

func main() {
    // Requires: QREDEX_CLIENT_ID, QREDEX_CLIENT_SECRET
    q, err := qredex.Bootstrap()
    if err != nil {
        log.Fatal(err)
    }

    // Use q.Creators(), q.Links(), etc.
}
```

**Explicit configuration:**

```go
q, err := qredex.New(qredex.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    Environment:  qredex.Production,
    Scopes: []qredex.Scope{
        qredex.ScopeCreatorsWrite,
        qredex.ScopeLinksWrite,
        qredex.ScopeIntentsWrite,
        qredex.ScopeOrdersWrite,
    },
})
```

### 2. Canonical Flow: IIT → PIT → Paid Order → Refund

```go
package main

import (
    "context"
    "log"
    "github.com/Qredex/qredex-go"
)

func main() {
    q, err := qredex.Bootstrap()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Step 1: Create a creator
    creator, err := q.Creators().Create(ctx, qredex.CreateCreatorRequest{
        Handle:      "alice",
        DisplayName: strPtr("Alice"),
        Email:       strPtr("alice@example.com"),
    })
    if err != nil {
        log.Fatal(err)
    }

    // Step 2: Create an influence link
    link, err := q.Links().Create(ctx, qredex.CreateLinkRequest{
        StoreID:         "store-123",
        CreatorID:       creator.ID,
        LinkName:        "spring-launch",
        DestinationPath: "/products/spring",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Step 3: Issue an Influence Intent Token (IIT)
    iit, err := q.Intents().IssueInfluenceIntentToken(ctx, qredex.IssueInfluenceIntentTokenRequest{
        LinkID: link.ID,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Step 4: Lock a Purchase Intent Token (PIT)
    pit, err := q.Intents().LockPurchaseIntent(ctx, qredex.LockPurchaseIntentRequest{
        Token:  iit.Token,
        Source: strPtr("backend-cart"),
    })
    if err != nil {
        log.Fatal(err)
    }

    // Step 5: Record a paid order
    order, err := q.Orders().RecordPaidOrder(ctx, qredex.RecordPaidOrderRequest{
        StoreID:             "store-123",
        ExternalOrderID:     "order-100045",
        Currency:            "USD",
        TotalPrice:          floatPtr(110.00),
        PurchaseIntentToken: strPtr(pit.Token),
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Order attribution: %s (status: %s)", order.ID, order.ResolutionStatus)

    // Step 6: Record a refund (if needed)
    refund, err := q.Refunds().RecordRefund(ctx, qredex.RecordRefundRequest{
        StoreID:          "store-123",
        ExternalOrderID:  "order-100045",
        ExternalRefundID: "refund-100045-1",
        RefundTotal:      floatPtr(25.00),
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Refund recorded: %s", refund.ID)
}

// Helper functions
func strPtr(s string) *string     { return &s }
func floatPtr(f float64) *float64 { return &f }
```

## Environment Variables

The SDK reads these environment variables when using `Bootstrap()`:

| Variable | Required | Description |
|----------|----------|-------------|
| `QREDEX_CLIENT_ID` | **Yes** | OAuth client ID |
| `QREDEX_CLIENT_SECRET` | **Yes** | OAuth client secret |
| `QREDEX_SCOPE` | No | Space-separated OAuth scopes |
| `QREDEX_ENVIRONMENT` | No | `production` (default), `staging`, or `development` |
| `QREDEX_BASE_URL` | No | Override base URL (bypasses environment) |

## Configuration

### OAuth Scopes

Available scopes for the Integrations API:

```go
qredex.ScopeAPI              // "direct:api"
qredex.ScopeLinksRead        // "direct:links:read"
qredex.ScopeLinksWrite       // "direct:links:write"
qredex.ScopeCreatorsRead     // "direct:creators:read"
qredex.ScopeCreatorsWrite    // "direct:creators:write"
qredex.ScopeOrdersRead       // "direct:orders:read"
qredex.ScopeOrdersWrite      // "direct:orders:write"
qredex.ScopeIntentsRead      // "direct:intents:read"
qredex.ScopeIntentsWrite     // "direct:intents:write"
```

### Environments

```go
qredex.Production    // https://api.qredex.com
qredex.Staging       // https://staging-api.qredex.com
qredex.Development   // http://localhost:8080
```

### Custom HTTP Client

You can provide a custom HTTP client:

```go
customClient := &http.Client{
    Timeout: 30 * time.Second,
    // ... custom transport, etc.
}

q, err := qredex.New(qredex.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    HTTPClient:   customClient,
})
```

## Error Handling

The SDK provides typed errors for different failure categories:

```go
import "github.com/Qredex/qredex-go"

_, err := q.Creators().Create(ctx, req)
if err != nil {
    if qredex.IsAuthenticationError(err) {
        // HTTP 401 — invalid/expired credentials
    } else if qredex.IsAuthorizationError(err) {
        // HTTP 403 — insufficient scope
    } else if qredex.IsValidationError(err) {
        // HTTP 400 — invalid request payload
    } else if qredex.IsNotFoundError(err) {
        // HTTP 404 — resource not found
    } else if qredex.IsConflictError(err) {
        // HTTP 409 — duplicate or conflict
    } else if qredex.IsRateLimitError(err) {
        // HTTP 429 — rate limited
        if rateLimitErr, ok := err.(*qredex.RateLimitError); ok {
            log.Printf("Retry after %d seconds", rateLimitErr.RetryAfterSeconds)
        }
    } else if qredex.IsNetworkError(err) {
        // Transport/network failure
    }
}
```

### Error Information

All API errors preserve platform contract details:

```go
if apiErr, ok := err.(*qredex.APIError); ok {
    log.Printf("Status: %d", apiErr.Status)
    log.Printf("Error Code: %s", apiErr.ErrorCode)
    log.Printf("Message: %s", apiErr.Message)
    log.Printf("Request ID: %s", apiErr.RequestID)
    log.Printf("Trace ID: %s", apiErr.TraceID)
}
```


## Operational Defaults

### Token Management

- Tokens are automatically issued and cached
- Cached tokens are reused until 30 seconds before expiry
- No manual token management required

### Retry Behavior

- **Read operations (GET)**: Optional retries with exponential backoff (configurable via `RetryMax`)
- **Write operations (POST/PUT)**: No automatic retries (must be explicit and safe)

### Security

- Client secrets are never logged
- Bearer tokens are redacted from debug output
- IIT/PIT tokens are not logged by default

## Observability

The SDK supports pluggable observability hooks for logging, tracing, and metrics. All hooks are optional and safe by default:

- **No secrets, tokens, or PII are ever logged, traced, or recorded as metrics.**
- Hooks are configured via the `Config` struct.

### Logger

Implement the `Logger` interface to receive log messages (e.g., for HTTP requests, retries):

```go
type Logger interface {
    Printf(format string, v ...interface{})
}
```

Example:
```go
q, err := qredex.New(qredex.Config{
    Logger: log.Default(), // stdlib logger
    // ...other config
})
```

### Tracer

Implement the `Tracer` interface to receive trace events (e.g., request lifecycle):

```go
type Tracer interface {
    Trace(event string, fields map[string]interface{})
}
```

Example:
```go
q, err := qredex.New(qredex.Config{
    Tracer: myTracer{}, // implements Trace
    // ...other config
})
```

### Metrics

Implement the `Metrics` interface to record metrics (e.g., request count, error count, latency):

```go
type Metrics interface {
    Record(metric string, value float64, labels map[string]string)
}
```

Example:
```go
q, err := qredex.New(qredex.Config{
    Metrics: myMetrics{}, // implements Record
    // ...other config
})
```

**Note:** All observability hooks are optional. If not set, the SDK operates silently and safely by default.

## API Reference

### Creators

```go
// Create a creator
creator, err := q.Creators().Create(ctx, qredex.CreateCreatorRequest{
    Handle:      "alice",
    DisplayName: strPtr("Alice"),
})

// Get a creator
creator, err := q.Creators().Get(ctx, creatorID)

// List creators
page, err := q.Creators().List(ctx, qredex.ListCreatorsRequest{
    Page:   intPtr(1),
    Size:   intPtr(10),
    Status: (*qredex.CreatorStatus)(strPtr("ACTIVE")),
})
```

### Links

```go
// Create a link
link, err := q.Links().Create(ctx, qredex.CreateLinkRequest{
    StoreID:         "store-123",
    CreatorID:       creatorID,
    LinkName:        "spring-launch",
    DestinationPath: "/products/spring",
})

// Get a link
link, err := q.Links().Get(ctx, linkID)

// List links
page, err := q.Links().List(ctx, qredex.ListLinksRequest{
    Page:   intPtr(1),
    Size:   intPtr(10),
})

// Get link stats
stats, err := q.Links().GetStats(ctx, linkID)
```

### Intents

```go
// Issue an Influence Intent Token (IIT)
iit, err := q.Intents().IssueInfluenceIntentToken(ctx, qredex.IssueInfluenceIntentTokenRequest{
    LinkID:      linkID,
    LandingPath: strPtr("/products/spring"),
})

// Lock a Purchase Intent Token (PIT)
pit, err := q.Intents().LockPurchaseIntent(ctx, qredex.LockPurchaseIntentRequest{
    Token:  iit.Token,
    Source: strPtr("backend-cart"),
})

// Get PIT details
pit, err := q.Intents().GetPurchaseIntent(ctx, pitToken)

// Get latest unlocked PIT
pit, err := q.Intents().GetLatestUnlocked(ctx, intPtr(24))
```

### Orders

```go
// Record a paid order
order, err := q.Orders().RecordPaidOrder(ctx, qredex.RecordPaidOrderRequest{
    StoreID:             "store-123",
    ExternalOrderID:     "order-100045",
    Currency:            "USD",
    TotalPrice:          floatPtr(110.00),
    PurchaseIntentToken: strPtr(pitToken),
})

// List orders
page, err := q.Orders().List(ctx, qredex.ListOrdersRequest{
    Page: intPtr(1),
    Size: intPtr(10),
})

// Get order details
details, err := q.Orders().GetDetails(ctx, orderAttributionID)
```

### Refunds

```go
// Record a refund
refund, err := q.Refunds().RecordRefund(ctx, qredex.RecordRefundRequest{
    StoreID:          "store-123",
    ExternalOrderID:  "order-100045",
    ExternalRefundID: "refund-100045-1",
    RefundTotal:      floatPtr(25.00),
})
```

## Testing

### Unit Tests

Run the test suite:

```bash
go test ./...
```

### Live Integration Tests

Live tests are excluded by default. To run them:

```bash
go test ./... -tags=live
```

Requires environment variables:
- `QREDEX_LIVE_ENABLED=1`
- `QREDEX_LIVE_CLIENT_ID`
- `QREDEX_LIVE_CLIENT_SECRET`
- `QREDEX_LIVE_STORE_ID`
- `QREDEX_LIVE_CREATOR_ID`

## Documentation

- [Integration Guide](docs/INTEGRATION_GUIDE.md) — Complete integration walkthrough
- [Error Handling](docs/ERRORS.md) — Detailed error model documentation
- [Go Reference](https://pkg.go.dev/github.com/Qredex/qredex-go) — API documentation

## Examples

See the [`examples/`](examples/) directory for complete examples:

- `canonical-flow` — Complete IIT → PIT → paid → refund flow
- `create-creator` — Create a creator
- `create-link` — Create an influence link
- `issue-iit` — Issue an Influence Intent Token
- `lock-pit` — Lock a Purchase Intent Token
- `record-order` — Record a paid order
- `list-orders` — List order attributions
- `record-refund` — Record a refund

## Releases

This SDK follows semantic versioning. Releases are tagged and published on GitHub and Packagist.

See [CHANGELOG.md](CHANGELOG.md) for version history and [docs/RELEASING.md](docs/RELEASING.md) for release procedures.

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.

## Security

If you discover a security vulnerability, please email security@qredex.com before opening a public issue.

## License

Apache License 2.0 — See [LICENSE](LICENSE) for details.

## Qredex Contact

- **Website**: https://qredex.com
- **X**: https://x.com/qredex
- **Email**: os@qredex.com

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

# Qredex Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/Qredex/qredex-go.svg)](https://pkg.go.dev/github.com/Qredex/qredex-go)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/Qredex/qredex-go)](https://goreportcard.com/report/github.com/Qredex/qredex-go)

The official Qredex Go server SDK for machine-to-machine integrations with the Qredex Integrations API.

## What This SDK Is

- A server-side SDK for merchant backends and partner platforms.
- A resource-oriented client for creators, links, intents, orders, and refunds.
- An integrations client that handles OAuth client-credentials authentication automatically.

## What This SDK Is Not

- Not a browser or storefront SDK.
- Not a Merchant API or Internal API client.
- Not a webhook framework.
- Not a transport-first wrapper that expects callers to manage raw HTTP plumbing.

## Installation

```bash
go get github.com/Qredex/qredex-go
```

Requirements:

- Go 1.21+
- Qredex Integrations client credentials

## Quick Start

```go
package main

import (
	"context"
	"log"

	"github.com/Qredex/qredex-go"
)

func main() {
	client, err := qredex.Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	creator, err := client.Creators().Create(ctx, qredex.CreateCreatorRequest{
		Handle: "alice",
	})
	if err != nil {
		log.Fatal(err)
	}

	link, err := client.Links().Create(ctx, qredex.CreateLinkRequest{
		StoreID:         "store_123",
		CreatorID:       creator.ID,
		LinkName:        "spring-launch",
		DestinationPath: "/products/spring",
	})
	if err != nil {
		log.Fatal(err)
	}

	iit, err := client.Intents().IssueInfluenceIntentToken(ctx, qredex.IssueInfluenceIntentTokenRequest{
		LinkID: link.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	pit, err := client.Intents().LockPurchaseIntent(ctx, qredex.LockPurchaseIntentRequest{
		Token:  iit.Token,
		Source: qredex.String("backend-cart"),
	})
	if err != nil {
		log.Fatal(err)
	}

	order, err := client.Orders().RecordPaidOrder(ctx, qredex.RecordPaidOrderRequest{
		StoreID:             "store_123",
		ExternalOrderID:     "order-100045",
		Currency:            "USD",
		TotalPrice:          qredex.Float64(110.00),
		PurchaseIntentToken: qredex.String(pit.Token),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("recorded order attribution: %s (%s)", order.ID, order.ResolutionStatus)
}

```

## Canonical Backend Flow

The intended Qredex backend flow is:

1. Create creator
2. Create link
3. Issue IIT
4. Lock PIT
5. Record paid order
6. Record refund when needed

The SDK is optimized for that sequence. Do not use it as a browser/session helper, and do not log raw IIT or PIT values.

## Initialization

`Bootstrap()` reads:

| Variable | Required | Description |
| --- | --- | --- |
| `QREDEX_CLIENT_ID` | Yes | OAuth client ID |
| `QREDEX_CLIENT_SECRET` | Yes | OAuth client secret |
| `QREDEX_SCOPE` | No | Space-separated OAuth scopes |
| `QREDEX_ENVIRONMENT` | No | `production`, `staging`, or `development` |
| `QREDEX_BASE_URL` | No | Explicit base URL override |
| `QREDEX_TIMEOUT_MS` | No | Per-request timeout in milliseconds |

Explicit configuration is also supported:

```go
client, err := qredex.New(qredex.Config{
	ClientID:     "client-id",
	ClientSecret: "client-secret",
	Environment:  qredex.Production,
	Timeout:      15 * time.Second,
	RetryMax:     2,
	Scopes: []qredex.Scope{
		qredex.ScopeCreatorsWrite,
		qredex.ScopeLinksWrite,
		qredex.ScopeIntentsWrite,
		qredex.ScopeOrdersWrite,
	},
})
```

## Operational Behavior

### Authentication

- OAuth client-credentials tokens are fetched automatically.
- Tokens are cached in memory and refreshed 30 seconds before expiry.
- A 401 clears the cache and triggers one immediate token refresh retry.

### Validation

- The SDK validates required request fields before making a network call.
- Invalid local input returns `RequestValidationError`.
- Server-side 400 responses still return `ValidationError`.

### Retries

- Safe read operations can retry on transport errors, 429, and 5xx when `RetryMax > 0`.
- Writes are not retried automatically.
- `Retry-After` is honored on 429 responses.

### Observability

- Optional logger, tracer, and metrics hooks are supported.
- Diagnostics use sanitized request paths.
- Secrets, bearer tokens, IITs, and PITs must not be logged downstream.

### Idempotency

- Use stable `ExternalOrderID` and `ExternalRefundID` values.
- If you configure `IdempotencyKeyProvider`, generate deterministic keys from your external identifiers.
- Do not implement blind write retries without a deliberate idempotency strategy.

## Errors

The SDK exposes typed errors for configuration, request validation, transport failures, and API failures.

Important types:

- `ConfigurationError`
- `RequestValidationError`
- `ResponseDecodingError`
- `NetworkError`
- `APIError`
- `AuthenticationError`
- `AuthorizationError`
- `ValidationError`
- `NotFoundError`
- `ConflictError`
- `RateLimitError`

API errors preserve `Status`, `ErrorCode`, `RequestID`, and `TraceID`.

See [docs/ERRORS.md](docs/ERRORS.md) for handling guidance.

## Documentation

- [docs/INTEGRATION_GUIDE.md](docs/INTEGRATION_GUIDE.md) for the canonical backend flow
- [docs/API_REFERENCE.md](docs/API_REFERENCE.md) for the public SDK surface
- [docs/ERRORS.md](docs/ERRORS.md) for the error model
- [docs/RELEASING.md](docs/RELEASING.md) for release discipline

## Examples

- [examples/canonical-flow/main.go](examples/canonical-flow/main.go)
- [examples/create-creator/main.go](examples/create-creator/main.go)
- [examples/create-link/main.go](examples/create-link/main.go)
- [examples/issue-iit/main.go](examples/issue-iit/main.go)
- [examples/lock-pit/main.go](examples/lock-pit/main.go)
- [examples/record-order/main.go](examples/record-order/main.go)
- [examples/list-orders/main.go](examples/list-orders/main.go)
- [examples/record-refund/main.go](examples/record-refund/main.go)

## Development

Canonical validation commands:

```bash
go build ./...
go test ./...
go vet ./...
golangci-lint run ./...
```

## Security

Report vulnerabilities privately to `security@qredex.com`.

## License

Apache License 2.0. See [LICENSE](LICENSE).

## Qredex Contact

- Website: https://qredex.com
- X: https://x.com/qredex
- Email: os@qredex.com

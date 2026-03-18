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

# Qredex Go SDK Integration Guide

This guide covers the intended backend integration path for the Qredex Go SDK.

## Before You Start

You need:

- Qredex Integrations client credentials
- A backend process that can keep secrets
- Your store identifier
- A place to persist the locked PIT token until paid-order submission

Do not use this SDK directly from a browser or embedded storefront runtime.

## Initialization

Environment-based bootstrap:

```go
client, err := qredex.Bootstrap()
```

Explicit configuration:

```go
client, err := qredex.New(qredex.Config{
	ClientID:     "client-id",
	ClientSecret: "client-secret",
	Environment:  qredex.Production,
	Timeout:      15 * time.Second,
	RetryMax:     2,
})
```

## Canonical Flow

### 1. Create Creator

```go
creator, err := client.Creators().Create(ctx, qredex.CreateCreatorRequest{
	Handle: "alice",
})
```

### 2. Create Link

```go
link, err := client.Links().Create(ctx, qredex.CreateLinkRequest{
	StoreID:         "store_123",
	CreatorID:       creator.ID,
	LinkName:        "spring-launch",
	DestinationPath: "/products/spring",
})
```

### 3. Issue IIT

Issue the IIT when the backend decides attribution should begin for a visit or click source.

```go
iit, err := client.Intents().IssueInfluenceIntentToken(ctx, qredex.IssueInfluenceIntentTokenRequest{
	LinkID: link.ID,
})
```

Do not log the raw IIT token. Persist or forward it only through your intended secure flow.

### 4. Lock PIT

Lock the PIT from the backend when purchase intent becomes real.

```go
pit, err := client.Intents().LockPurchaseIntent(ctx, qredex.LockPurchaseIntentRequest{
	Token:  iit.Token,
	Source: qredex.String("backend-cart"),
})
```

Persist the PIT token with the cart or order context. It is the canonical input for paid-order attribution.

### 5. Record Paid Order

```go
order, err := client.Orders().RecordPaidOrder(ctx, qredex.RecordPaidOrderRequest{
	StoreID:             "store_123",
	ExternalOrderID:     "order-100045",
	Currency:            "USD",
	TotalPrice:          qredex.Float64(110.00),
	PurchaseIntentToken: qredex.String(pit.Token),
})
```

Use a stable `ExternalOrderID`. That is your first idempotency boundary.

### 6. Record Refund

```go
updated, err := client.Refunds().RecordRefund(ctx, qredex.RecordRefundRequest{
	StoreID:          "store_123",
	ExternalOrderID:  "order-100045",
	ExternalRefundID: "refund-100045-1",
	RefundTotal:      qredex.Float64(25.00),
})
```

## Validation Behavior

The SDK validates obvious request mistakes before sending any HTTP request. Examples:

- missing required identifiers
- empty required strings
- invalid `destination_path`
- invalid currency codes
- negative monetary amounts

These failures return `RequestValidationError`.

## Retry and Idempotency Guidance

- Reads can retry when `RetryMax > 0`.
- Writes are not retried automatically.
- Use stable external identifiers.
- If you add an `IdempotencyKeyProvider`, derive the key from your own durable identifiers.

## Observability Guidance

- Logger, tracer, and metrics hooks are optional.
- The SDK emits sanitized paths, not raw secret-bearing URLs.
- Keep the same discipline in downstream sinks: never log bearer tokens, IITs, or PITs.

## Deprecated Helper

`Intents().GetLatestUnlocked()` remains available for advanced recovery cases, but it is nondeterministic and should not be part of the normal paid-order flow.

## Production Checklist

- Credentials stored in secret management, not source control
- Explicit timeout configured
- External order/refund identifiers are stable
- Raw IIT/PIT tokens are never logged
- Paid-order and refund flows tested in staging
- `go build ./...`, `go test ./...`, `go vet ./...`, and `golangci-lint run ./...` pass before release

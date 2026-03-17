# Qredex Go SDK Integration Guide

This guide walks you through integrating the Qredex Go SDK into your merchant backend or partner platform.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Authentication Setup](#authentication-setup)
4. [Canonical Integration Flow](#canonical-integration-flow)
5. [Resource Operations](#resource-operations)
6. [Error Handling Best Practices](#error-handling-best-practices)
7. [Production Checklist](#production-checklist)

---

## Prerequisites

Before integrating the Qredex Go SDK, ensure you have:

1. **Qredex Merchant Account** — Active merchant account with Qredex
2. **OAuth Credentials** — Client ID and Client Secret from Qredex dashboard
3. **Store ID** — Your Qredex store identifier
4. **Go 1.21+** — The SDK requires Go 1.21 or later

### Obtaining OAuth Credentials

1. Log in to your Qredex Merchant Dashboard
2. Navigate to **Settings → Integrations**
3. Click **Create New Integration**
4. Select **Server-to-Server (OAuth 2.0)**
5. Copy the Client ID and Client Secret

**Important:** Store the Client Secret securely. It cannot be recovered if lost.

---

## Installation

```bash
go get github.com/Qredex/qredex-go
```

Add to your `go.mod`:

```go
require github.com/Qredex/qredex-go v0.1.0
```

---

## Authentication Setup

### Environment-Based Configuration (Recommended)

Set environment variables in your deployment:

```bash
export QREDEX_CLIENT_ID="your-client-id"
export QREDEX_CLIENT_SECRET="your-client-secret"
export QREDEX_ENVIRONMENT="production"  # or "staging", "development"
export QREDEX_SCOPE="direct:creators:write direct:links:write direct:intents:write direct:orders:write"
```

Initialize in your application:

```go
package main

import (
    "log"
    "github.com/Qredex/qredex-go"
)

func initQredex() (*qredex.Qredex, error) {
    return qredex.Bootstrap()
}
```

### Explicit Configuration

For more control:

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
    Timeout: 15 * time.Second,
})
```

### Custom HTTP Client

For advanced scenarios (custom TLS, proxies, etc.):

```go
customClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

q, err := qredex.New(qredex.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    HTTPClient:   customClient,
})
```

---

## Canonical Integration Flow

The canonical flow for Qredex attribution is:

```
1. Create Creator → 2. Create Link → 3. Issue IIT → 4. Lock PIT → 5. Record Paid Order → 6. Record Refund (optional)
```

### Step 1: Create a Creator

```go
creator, err := q.Creators().Create(ctx, qredex.CreateCreatorRequest{
    Handle:      "alice",
    DisplayName: strPtr("Alice"),
    Email:       strPtr("alice@example.com"),
    Socials:     map[string]string{"twitter": "@alice"},
})
if err != nil {
    log.Fatal(err)
}

log.Printf("Created creator: %s (%s)", creator.Handle, creator.ID)
```

### Step 2: Create an Influence Link

```go
link, err := q.Links().Create(ctx, qredex.CreateLinkRequest{
    StoreID:               "store-123",
    CreatorID:             creator.ID,
    LinkName:              "spring-launch",
    DestinationPath:       "/products/spring",
    AttributionWindowDays: intPtr(30),
})
if err != nil {
    log.Fatal(err)
}

log.Printf("Created link: %s (%s)", link.LinkName, link.ID)
log.Printf("Public URL: %s", link.PublicLinkURL)
```

### Step 3: Issue an Influence Intent Token (IIT)

The IIT is issued when a user clicks a link or when you want to attribute traffic from a known source.

```go
iit, err := q.Intents().IssueInfluenceIntentToken(ctx, qredex.IssueInfluenceIntentTokenRequest{
    LinkID:      link.ID,
    LandingPath: strPtr("/products/spring"),
    IPHash:      strPtr(hashIP(userIP)),
})
if err != nil {
    log.Fatal(err)
}

log.Printf("Issued IIT: %s", iit.Token)
```

**Important:** The IIT token is a JWT. Pass it to your frontend or store it for later PIT locking.

### Step 4: Lock a Purchase Intent Token (PIT)

Lock the PIT when the user shows purchase intent (e.g., adds to cart, initiates checkout).

```go
pit, err := q.Intents().LockPurchaseIntent(ctx, qredex.LockPurchaseIntentRequest{
    Token:  iit.Token,  // The IIT from Step 3
    Source: strPtr("backend-cart"),
})
if err != nil {
    log.Fatal(err)
}

log.Printf("Locked PIT: %s", pit.Token)
log.Printf("Eligible: %v", *pit.Eligible)
```

**Important:** Store the PIT token with the user's session or cart. It will be used when recording the paid order.

### Step 5: Record a Paid Order

After the order is paid, submit it for attribution.

```go
order, err := q.Orders().RecordPaidOrder(ctx, qredex.RecordPaidOrderRequest{
    StoreID:             "store-123",
    ExternalOrderID:     "order-100045",  // Your internal order ID
    OrderNumber:         strPtr("100045"),
    Currency:            "USD",
    TotalPrice:          floatPtr(110.00),
    PaidAt:              &paidAt,
    PurchaseIntentToken: strPtr(pit.Token),  // The PIT from Step 4
})
if err != nil {
    log.Fatal(err)
}

log.Printf("Order attribution: %s", order.ID)
log.Printf("Resolution Status: %s", order.ResolutionStatus)
log.Printf("Token Integrity: %v", order.TokenIntegrity)
log.Printf("Attributed to Creator: %s", order.CreatorHandle)
```

### Step 6: Record a Refund (Optional)

If the order is refunded, submit the refund for attribution adjustment.

```go
refund, err := q.Refunds().RecordRefund(ctx, qredex.RecordRefundRequest{
    StoreID:          "store-123",
    ExternalOrderID:  "order-100045",
    ExternalRefundID: "refund-100045-1",  // Your internal refund ID
    RefundTotal:      floatPtr(25.00),
    RefundedAt:       &refundedAt,
})
if err != nil {
    log.Fatal(err)
}

log.Printf("Refund recorded: %s", refund.ID)
```

---

## Resource Operations

### Creators

```go
// Get a creator
creator, err := q.Creators().Get(ctx, creatorID)

// List creators (paginated)
page, err := q.Creators().List(ctx, qredex.ListCreatorsRequest{
    Page:   intPtr(1),
    Size:   intPtr(20),
    Status: (*qredex.CreatorStatus)(strPtr("ACTIVE")),
})

for _, creator := range page.Items {
    log.Printf("%s: %s", creator.Handle, creator.Status)
}
```

### Links

```go
// Get a link
link, err := q.Links().Get(ctx, linkID)

// List links (paginated)
page, err := q.Links().List(ctx, qredex.ListLinksRequest{
    Page:   intPtr(1),
    Size:   intPtr(20),
    Status: (*qredex.LinkStatus)(strPtr("ACTIVE")),
})

// Get link stats
stats, err := q.Links().GetStats(ctx, linkID)
log.Printf("Clicks: %d, Orders: %d, Revenue: %.2f", 
    stats.ClicksCount, stats.OrdersCount, stats.RevenueTotal)
```

### Intents

```go
// Get PIT by token
pit, err := q.Intents().GetPurchaseIntent(ctx, pitToken)

// Get latest unlocked PIT (e.g., for last 24 hours)
pit, err := q.Intents().GetLatestUnlocked(ctx, intPtr(24))
```

### Orders

```go
// List orders (paginated)
page, err := q.Orders().List(ctx, qredex.ListOrdersRequest{
    Page: intPtr(1),
    Size: intPtr(20),
})

// Get order details (includes score breakdown and timeline)
details, err := q.Orders().GetDetails(ctx, orderAttributionID)
log.Printf("Integrity Score: %d (%s)", details.IntegrityScore, details.IntegrityBand)
log.Printf("Review Required: %v", details.ReviewRequired)
```

---

## Error Handling Best Practices

### Categorize Errors

```go
func handleAttributionError(err error) {
    if qredex.IsAuthenticationError(err) {
        // Check credentials, token expiry
        log.Printf("Auth error: %v", err)
    } else if qredex.IsAuthorizationError(err) {
        // Check OAuth scopes
        log.Printf("Authz error: %v", err)
    } else if qredex.IsValidationError(err) {
        // Fix request payload
        log.Printf("Validation error: %v", err)
    } else if qredex.IsConflictError(err) {
        // Handle duplicate order
        log.Printf("Conflict: %v", err)
    } else if qredex.IsRateLimitError(err) {
        // Wait and retry
        if rateLimitErr, ok := err.(*qredex.RateLimitError); ok {
            time.Sleep(time.Duration(rateLimitErr.RetryAfterSeconds) * time.Second)
        }
    } else if qredex.IsNetworkError(err) {
        // Retry with backoff
        log.Printf("Network error: %v", err)
    } else {
        // Unexpected error
        log.Printf("Unexpected error: %v", err)
    }
}
```

### Preserve Error Context

```go
if apiErr, ok := err.(*qredex.APIError); ok {
    log.Printf("Status: %d", apiErr.Status)
    log.Printf("Error Code: %s", apiErr.ErrorCode)
    log.Printf("Request ID: %s", apiErr.RequestID)  // For support tickets
    log.Printf("Trace ID: %s", apiErr.TraceID)      // For debugging
}
```

### Idempotency for Paid Orders

Paid order submission should be idempotent. Use the same `ExternalOrderID` for retries:

```go
func recordPaidOrderWithRetry(ctx context.Context, orderID string, req qredex.RecordPaidOrderRequest) (*qredex.OrderAttribution, error) {
    for attempt := 0; attempt < 3; attempt++ {
        order, err := q.Orders().RecordPaidOrder(ctx, req)
        if err == nil {
            return order, nil
        }
        
        if qredex.IsConflictError(err) {
            // Order already recorded - fetch existing attribution
            log.Printf("Order %s already recorded", orderID)
            break
        }
        
        if !qredex.IsNetworkError(err) && !qredex.IsRateLimitError(err) {
            return nil, err
        }
        
        // Retry with backoff
        time.Sleep(time.Duration(attempt+1) * time.Second)
    }
    
    return nil, fmt.Errorf("failed to record order after retries")
}
```

---

## Production Checklist

Before deploying to production:

### Configuration

- [ ] OAuth credentials stored securely (not in code)
- [ ] Using production environment (`QREDEX_ENVIRONMENT=production`)
- [ ] Appropriate OAuth scopes configured
- [ ] HTTP timeout configured (15-30 seconds recommended)

### Error Handling

- [ ] All API errors are handled
- [ ] Network errors trigger retries with backoff
- [ ] Rate limit errors respect `Retry-After` headers
- [ ] Conflict errors are handled for idempotent order submission

### Observability

- [ ] Request IDs are logged for debugging
- [ ] Error codes are captured for monitoring
- [ ] Attribution resolution status is tracked

### Testing

- [ ] Tested in staging environment first
- [ ] Canonical flow tested end-to-end
- [ ] Error scenarios tested (invalid credentials, validation failures, etc.)
- [ ] Load tested for expected traffic volume

### Security

- [ ] Client secret not logged or exposed
- [ ] Bearer tokens not logged
- [ ] IIT/PIT tokens handled securely
- [ ] HTTPS enforced in production

---

## Next Steps

- Review the [API Reference](https://pkg.go.dev/github.com/Qredex/qredex-go) for detailed method documentation
- See [ERRORS.md](ERRORS.md) for comprehensive error handling guidance
- Check out the [examples/](examples/) directory for complete code samples
- Contact os@qredex.com for support or questions

# Qredex Go SDK Error Handling

This document describes the error model for the Qredex Go SDK.

## Error Hierarchy

The SDK defines a typed error hierarchy that separates failure categories:

```
error
├── ConfigurationError    — SDK misconfiguration (before any request)
├── NetworkError          — Transport/network failures
└── APIError              — HTTP 4xx/5xx responses from Qredex API
    ├── AuthenticationError  — HTTP 401
    ├── AuthorizationError   — HTTP 403
    ├── ValidationError      — HTTP 400
    ├── NotFoundError        — HTTP 404
    ├── ConflictError        — HTTP 409
    └── RateLimitError       — HTTP 429
```

## Error Detection

Use the provided helper functions to detect error types:

```go
import "github.com/Qredex/qredex-go"

_, err := q.Creators().Create(ctx, req)
if err != nil {
    if qredex.IsAuthenticationError(err) {
        // Handle 401
    } else if qredex.IsAuthorizationError(err) {
        // Handle 403
    } else if qredex.IsValidationError(err) {
        // Handle 400
    } else if qredex.IsNotFoundError(err) {
        // Handle 404
    } else if qredex.IsConflictError(err) {
        // Handle 409
    } else if qredex.IsRateLimitError(err) {
        // Handle 429
    } else if qredex.IsNetworkError(err) {
        // Handle network failure
    } else if qredex.IsAPIError(err) {
        // Handle any API error
    }
}
```

## Error Types

### ConfigurationError

Returned when the SDK is misconfigured before any request is made.

**Causes:**
- Missing `ClientID`
- Missing `ClientSecret`
- Invalid environment configuration

**Example:**

```go
q, err := qredex.New(qredex.Config{
    ClientID: "",  // Empty!
    ClientSecret: "secret",
})
// err: *ConfigurationError: "Qredex Config requires ClientID"
```

**Handling:**
- Fix configuration before retrying
- This is a developer error, not a runtime error

### NetworkError

Wraps transport-level failures that occur before a valid HTTP response is received.

**Causes:**
- DNS resolution failures
- Connection timeouts
- TLS handshake failures
- Context cancellation
- Network unreachable

**Fields:**
```go
type NetworkError struct {
    Message string  // Description of the failure
    Cause   error   // Underlying error (if available)
}
```

**Handling:**
- Safe to retry with exponential backoff
- Check for context cancellation
- Verify network connectivity

**Example:**

```go
_, err := q.Creators().Create(ctx, req)
if qredex.IsNetworkError(err) {
    log.Printf("Network failure: %v", err)
    // Retry with backoff
}
```

### APIError (Base)

Base type for all HTTP 4xx/5xx responses from the Qredex API.

**Fields:**
```go
type APIError struct {
    Status    int     // HTTP status code
    ErrorCode string  // Qredex error code (machine-readable)
    Message   string  // Human-readable error message
    RequestID string  // X-Request-Id correlation header
    TraceID   string  // X-Trace-Id correlation header
}
```

**Handling:**
- Inspect `ErrorCode` for programmatic handling
- Log `RequestID` and `TraceID` for support tickets
- Use `Status` for general categorization

### AuthenticationError (401)

Returned on HTTP 401 responses.

**Causes:**
- Invalid client credentials
- Expired access token
- Revoked token
- Malformed Authorization header

**Example Response:**
```json
{
  "error_code": "invalid_token",
  "message": "The access token is invalid or expired"
}
```

**Handling:**
- Clear token cache: `qredex.Auth().ClearTokenCache()`
- Re-authenticate with fresh credentials
- Verify `ClientID` and `ClientSecret` are correct

**Example:**

```go
_, err := q.Creators().Create(ctx, req)
if qredex.IsAuthenticationError(err) {
    log.Printf("Auth failed: %v", err)
    // Clear cache and retry
    qredex.Auth().ClearTokenCache()
}
```

### AuthorizationError (403)

Returned on HTTP 403 responses.

**Causes:**
- Missing required OAuth scope
- Insufficient permissions
- Forbidden resource access

**Example Response:**
```json
{
  "error_code": "insufficient_scope",
  "message": "Missing required scope: direct:creators:write"
}
```

**Handling:**
- Verify OAuth scopes in configuration
- Request additional scopes from Qredex
- Check resource permissions

**Example:**

```go
_, err := q.Creators().Create(ctx, req)
if qredex.IsAuthorizationError(err) {
    log.Printf("Missing scope: %v", err)
    // Update config to include direct:creators:write
}
```

### ValidationError (400)

Returned on HTTP 400 responses.

**Causes:**
- Missing required fields
- Invalid field format (e.g., malformed UUID)
- Invalid field values (e.g., negative amount)
- Schema validation failures

**Example Response:**
```json
{
  "error_code": "validation_error",
  "message": "Invalid request: handle is required"
}
```

**Handling:**
- Fix request payload
- Validate input before sending
- Check field requirements in API docs

**Example:**

```go
_, err := q.Creators().Create(ctx, qredex.CreateCreatorRequest{
    Handle: "",  // Empty!
})
if qredex.IsValidationError(err) {
    log.Printf("Validation failed: %v", err)
    // Fix the request
}
```

### NotFoundError (404)

Returned on HTTP 404 responses.

**Causes:**
- Resource does not exist
- Invalid resource ID
- Resource deleted or archived

**Example Response:**
```json
{
  "error_code": "not_found",
  "message": "Creator not found"
}
```

**Handling:**
- Verify resource ID is correct
- Check if resource was deleted
- Handle gracefully in application logic

**Example:**

```go
creator, err := q.Creators().Get(ctx, "invalid-id")
if qredex.IsNotFoundError(err) {
    log.Printf("Creator not found: %v", err)
    // Handle missing resource
}
```

### ConflictError (409)

Returned on HTTP 409 responses.

**Causes:**
- Duplicate order submission
- Cross-source duplicate detection
- Idempotency conflict

**Example Response:**
```json
{
  "error_code": "REJECTED_CROSS_SOURCE_DUPLICATE",
  "message": "Order already exists from different source"
}
```

**Handling:**
- Check if operation already completed
- Use idempotent request patterns
- Fetch existing resource instead of creating

**Example:**

```go
order, err := q.Orders().RecordPaidOrder(ctx, req)
if qredex.IsConflictError(err) {
    log.Printf("Order already recorded: %v", err)
    // Fetch existing order attribution
}
```

### RateLimitError (429)

Returned on HTTP 429 responses.

**Causes:**
- Too many requests per second
- Exceeded daily quota
- Burst limit exceeded

**Fields:**
```go
type RateLimitError struct {
    APIError
    RetryAfterSeconds int  // Seconds to wait before retrying
}
```

**Example Response:**
```json
{
  "error_code": "rate_limited",
  "message": "Rate limit exceeded"
}
```

**Handling:**
- Wait for `RetryAfterSeconds` before retrying
- Implement request throttling
- Use exponential backoff

**Example:**

```go
for attempt := 0; attempt < 3; attempt++ {
    _, err := q.Creators().List(ctx, qredex.ListCreatorsRequest{})
    if err == nil {
        break
    }
    
    if rateLimitErr, ok := err.(*qredex.RateLimitError); ok {
        log.Printf("Rate limited, waiting %d seconds", rateLimitErr.RetryAfterSeconds)
        time.Sleep(time.Duration(rateLimitErr.RetryAfterSeconds) * time.Second)
        continue
    }
    
    return err
}
```

## Error Codes

The SDK preserves Qredex error codes from API responses. Common error codes:

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `invalid_client` | 401 | Invalid client credentials |
| `invalid_token` | 401 | Token expired or revoked |
| `insufficient_scope` | 403 | Missing required OAuth scope |
| `validation_error` | 400 | Request validation failed |
| `not_found` | 404 | Resource not found |
| `REJECTED_CROSS_SOURCE_DUPLICATE` | 409 | Duplicate from different source |
| `rate_limited` | 429 | Rate limit exceeded |
| `INTERNAL_ERROR` | 500 | Internal server error |

## Best Practices

### 1. Always Check Error Types

```go
// BAD: Generic error handling
if err != nil {
    log.Fatal(err)
}

// GOOD: Typed error handling
if err != nil {
    if qredex.IsValidationError(err) {
        // Fix request
    } else if qredex.IsConflictError(err) {
        // Handle duplicate
    } else {
        log.Fatal(err)
    }
}
```

### 2. Log Correlation Identifiers

```go
if apiErr, ok := err.(*qredex.APIError); ok {
    log.Printf("Error: %s [%s] - RequestID: %s, TraceID: %s",
        apiErr.Message, apiErr.ErrorCode, apiErr.RequestID, apiErr.TraceID)
}
```

### 3. Implement Retry Logic

```go
func withRetry(ctx context.Context, fn func() error) error {
    for attempt := 0; attempt < 3; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        // Don't retry client errors
        if qredex.IsValidationError(err) || 
           qredex.IsAuthenticationError(err) ||
           qredex.IsAuthorizationError(err) {
            return err
        }
        
        // Retry network errors and rate limits
        if qredex.IsRateLimitError(err) {
            if rateLimitErr, ok := err.(*qredex.RateLimitError); ok {
                time.Sleep(time.Duration(rateLimitErr.RetryAfterSeconds) * time.Second)
                continue
            }
        }
        
        if qredex.IsNetworkError(err) {
            time.Sleep(time.Duration(attempt+1) * time.Second)
            continue
        }
        
        return err
    }
    return fmt.Errorf("failed after retries")
}
```

### 4. Handle Idempotency

```go
func recordPaidOrder(ctx context.Context, orderID string, req qredex.RecordPaidOrderRequest) (*qredex.OrderAttribution, error) {
    order, err := q.Orders().RecordPaidOrder(ctx, req)
    if err == nil {
        return order, nil
    }
    
    // Handle duplicate gracefully
    if qredex.IsConflictError(err) {
        log.Printf("Order %s already recorded", orderID)
        // Fetch existing attribution
        return fetchExistingOrder(ctx, orderID)
    }
    
    return nil, err
}
```

### 5. Never Log Secrets

```go
// BAD: Logs sensitive data
log.Printf("Request: %+v", req)  // May include tokens!

// GOOD: Log only safe fields
log.Printf("Recording order: store=%s, external_id=%s", 
    req.StoreID, req.ExternalOrderID)
```

## Debugging

### Enable Debug Logging

The SDK does not log by default. To debug:

```go
// Wrap HTTP client with logging transport
type loggingTransport struct {
    base http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    log.Printf("Request: %s %s", req.Method, req.URL)
    resp, err := t.base.RoundTrip(req)
    if err != nil {
        log.Printf("Error: %v", err)
    } else {
        log.Printf("Response: %d", resp.StatusCode)
    }
    return resp, err
}

q, err := qredex.New(qredex.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    HTTPClient: &http.Client{
        Transport: &loggingTransport{base: http.DefaultTransport},
    },
})
```

### Contact Support

When contacting Qredex support, include:

1. `RequestID` from the error
2. `TraceID` from the error
3. Timestamp of the failure
4. Operation being performed
5. Relevant code snippet

Example:

```
RequestID: req-abc123
TraceID: trace-xyz789
Timestamp: 2026-01-15T10:30:00Z
Operation: RecordPaidOrder
Error: ConflictError [REJECTED_CROSS_SOURCE_DUPLICATE]
```

## Next Steps

- Review the [Integration Guide](INTEGRATION_GUIDE.md) for complete integration walkthrough
- See the [README](../README.md) for quick start
- Check [API Reference](https://pkg.go.dev/github.com/Qredex/qredex-go) for method documentation

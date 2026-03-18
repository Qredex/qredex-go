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

# Qredex Go SDK Errors

## Error Categories

The SDK separates failures into four buckets:

- `ConfigurationError`: client misconfiguration before any request
- `RequestValidationError`: invalid caller input rejected locally
- `ResponseDecodingError`: successful HTTP response that could not be decoded into the expected SDK model
- `NetworkError`: transport failure before a valid API response
- `APIError`: non-2xx API response from Qredex

`APIError` has typed subtypes:

- `AuthenticationError`
- `AuthorizationError`
- `ValidationError`
- `NotFoundError`
- `ConflictError`
- `RateLimitError`

## Local Validation vs API Validation

`RequestValidationError` means the SDK rejected bad input before the request was sent.

`ValidationError` means the request reached Qredex and the API rejected it with HTTP 400.

That distinction is intentional. It helps you separate caller bugs from platform responses.

## Preserved API Details

`APIError` preserves:

- `Status`
- `ErrorCode`
- `Message`
- `RequestID`
- `TraceID`

Use `RequestID` and `TraceID` in support escalations.

## Typical Handling Pattern

```go
_, err := client.Orders().RecordPaidOrder(ctx, req)
if err != nil {
	switch {
	case qredex.IsRequestValidationError(err):
		// Fix local caller input.
	case qredex.IsAuthenticationError(err):
		// Credentials or token flow failed.
	case qredex.IsAuthorizationError(err):
		// Missing scope or forbidden operation.
	case qredex.IsConflictError(err):
		// Idempotency boundary was hit.
	case qredex.IsRateLimitError(err):
		// Honor Retry-After.
	case qredex.IsResponseDecodingError(err):
		// API responded, but the body did not match the expected contract.
	case qredex.IsNetworkError(err):
		// Transport failure. Reads may be retried carefully.
	case qredex.IsAPIError(err):
		// Other platform-side API failure.
	default:
		// Unexpected error.
	}
}
```

## Authentication Behavior

The SDK automatically refreshes once on a 401 by clearing its in-memory token cache and fetching a fresh token. If the retry also fails, the caller receives `AuthenticationError`.

## Logging Guidance

- Safe: `status`, `error_code`, `requestId`, `traceId`
- Unsafe: raw bearer tokens, IITs, PITs, Authorization headers, full secret-bearing URLs

## Rate Limiting

`RateLimitError` preserves `RetryAfterSeconds` when the header is present. Reads with retries enabled already respect `Retry-After`. Write retry policy remains your responsibility.

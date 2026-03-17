// Copyright (C) 2026 — 2026, Qredex, LTD. All Rights Reserved.
//
// DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.
//
// Licensed under the Apache License, Version 2.0. See LICENSE for the full license text.
// You may not use this file except in compliance with that License.
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.
//
// If you need additional information or have any questions, please email: copyright@qredex.com

package qredex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// ConfigurationError is returned when the SDK is misconfigured before any request is made.
type ConfigurationError struct {
	Message string
}

func (e *ConfigurationError) Error() string {
	return "qredex: configuration error: " + e.Message
}

// APIError is the base error type for non-network API responses from Qredex.
// It carries the HTTP status, error_code, message, and correlation identifiers.
type APIError struct {
	// Status is the HTTP response status code.
	Status int
	// ErrorCode is the machine-readable Qredex error code from the response body.
	ErrorCode string
	// Message is the human-readable error message.
	Message string
	// RequestID is the X-Request-Id correlation header, when present.
	RequestID string
	// TraceID is the X-Trace-Id correlation header, when present.
	TraceID string
}

func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("qredex: API error %d [%s]: %s", e.Status, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("qredex: API error %d: %s", e.Status, e.Message)
}

// AuthenticationError is returned on HTTP 401 responses.
// This typically indicates invalid credentials or an expired/revoked token.
type AuthenticationError struct{ APIError }

// AuthorizationError is returned on HTTP 403 responses.
// This typically indicates a missing or insufficient OAuth scope.
type AuthorizationError struct{ APIError }

// ValidationError is returned on HTTP 400 responses.
// This indicates a malformed or invalid request payload.
type ValidationError struct{ APIError }

// NotFoundError is returned on HTTP 404 responses.
type NotFoundError struct{ APIError }

// ConflictError is returned on HTTP 409 responses.
// This indicates a duplicate order submission or cross-source conflict.
type ConflictError struct{ APIError }

// RateLimitError is returned on HTTP 429 responses.
// RetryAfterSeconds is populated from the Retry-After response header when present.
type RateLimitError struct {
	APIError
	// RetryAfterSeconds is the number of seconds to wait before retrying, if provided.
	RetryAfterSeconds int
}

// NetworkError wraps a transport-level failure that occurred before a valid HTTP response
// was received. This includes DNS failures, connection timeouts, and context cancellations.
type NetworkError struct {
	// Message describes the network failure.
	Message string
	// Cause is the underlying error, if available.
	Cause error
}

func (e *NetworkError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("qredex: network error: %s: %v", e.Message, e.Cause)
	}
	return "qredex: network error: " + e.Message
}

// Unwrap returns the underlying cause of the network error.
func (e *NetworkError) Unwrap() error { return e.Cause }

// IsAuthenticationError reports whether err is an AuthenticationError (HTTP 401).
func IsAuthenticationError(err error) bool {
	_, ok := err.(*AuthenticationError)
	return ok
}

// IsAuthorizationError reports whether err is an AuthorizationError (HTTP 403).
func IsAuthorizationError(err error) bool {
	_, ok := err.(*AuthorizationError)
	return ok
}

// IsValidationError reports whether err is a ValidationError (HTTP 400).
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsNotFoundError reports whether err is a NotFoundError (HTTP 404).
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

// IsConflictError reports whether err is a ConflictError (HTTP 409).
func IsConflictError(err error) bool {
	_, ok := err.(*ConflictError)
	return ok
}

// IsRateLimitError reports whether err is a RateLimitError (HTTP 429).
func IsRateLimitError(err error) bool {
	_, ok := err.(*RateLimitError)
	return ok
}

// IsNetworkError reports whether err is a NetworkError (transport-level failure).
func IsNetworkError(err error) bool {
	_, ok := err.(*NetworkError)
	return ok
}

// IsAPIError reports whether err is any API-originated Qredex error (non-network).
func IsAPIError(err error) bool {
	switch err.(type) {
	case *APIError, *AuthenticationError, *AuthorizationError, *ValidationError,
		*NotFoundError, *ConflictError, *RateLimitError:
		return true
	}
	return false
}

// apiErrorBody is the JSON shape of a Qredex API error response body.
type apiErrorBody struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// parseAPIError reads an HTTP error response and returns the appropriate typed error.
// The response body is consumed and closed by the caller.
func parseAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var parsed apiErrorBody
	_ = json.Unmarshal(body, &parsed)

	if parsed.Message == "" {
		parsed.Message = http.StatusText(resp.StatusCode)
	}

	requestID := resp.Header.Get("X-Request-Id")
	traceID := resp.Header.Get("X-Trace-Id")

	base := APIError{
		Status:    resp.StatusCode,
		ErrorCode: parsed.ErrorCode,
		Message:   parsed.Message,
		RequestID: requestID,
		TraceID:   traceID,
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &AuthenticationError{base}
	case http.StatusForbidden:
		return &AuthorizationError{base}
	case http.StatusBadRequest:
		return &ValidationError{base}
	case http.StatusNotFound:
		return &NotFoundError{base}
	case http.StatusConflict:
		return &ConflictError{base}
	case http.StatusTooManyRequests:
		retryAfter := 0
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			retryAfter, _ = strconv.Atoi(ra)
		}
		return &RateLimitError{APIError: base, RetryAfterSeconds: retryAfter}
	default:
		return &base
	}
}

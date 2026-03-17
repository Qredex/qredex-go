//	▄▄▄▄
//	▄█▀▀███▄▄              █▄
//	██    ██ ▄             ██
//	██    ██ ████▄▄█▀█▄ ▄████ ▄█▀█▄▀██ ██▀
//	██  ▄ ██ ██   ██▄█▀ ██ ██ ██▄█▀  ███
//	 ▀█████▄▄█▀  ▄▀█▄▄▄▄█▀███▄▀█▄▄▄▄██ ██▄
//	     ▀█
//
//	Copyright (C) 2026 — 2026, Qredex, LTD. All Rights Reserved.
//
//	DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.
//
//	Licensed under the Apache License, Version 2.0. See LICENSE for the full license text.
//	You may not use this file except in compliance with that License.
//	Unless required by applicable law or agreed to in writing, software distributed under the
//	License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//	either express or implied. See the License for the specific language governing permissions
//	and limitations under the License.
//
//	If you need additional information or have any questions, please email: copyright@qredex.com

package qredex

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

// TransportRequest represents an outgoing HTTP request for testing.
type TransportRequest struct {
	Method    string
	URL       string
	Headers   http.Header
	Body      []byte
	Timestamp time.Time
}

// TransportResponse represents a mocked HTTP response for testing.
type TransportResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// FakeTransport is a test double that implements http.RoundTripper.
// It records all requests and returns queued responses.
type FakeTransport struct {
	mu        sync.Mutex
	responses []TransportResponse
	requests  []TransportRequest
	errs      []error
	errIndex  int
	respIndex int
}

// NewFakeTransport creates a new FakeTransport instance.
func NewFakeTransport() *FakeTransport {
	return &FakeTransport{
		responses: make([]TransportResponse, 0),
		requests:  make([]TransportRequest, 0),
		errs:      make([]error, 0),
	}
}

// RoundTrip implements http.RoundTripper.
func (ft *FakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	// Record the request
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
	}
	ft.requests = append(ft.requests, TransportRequest{
		Method:    req.Method,
		URL:       req.URL.String(),
		Headers:   req.Header.Clone(),
		Body:      body,
		Timestamp: time.Now(),
	})

	// Check if we have a queued error
	if ft.errIndex < len(ft.errs) {
		err := ft.errs[ft.errIndex]
		ft.errIndex++
		return nil, err
	}

	// Return queued response
	if ft.respIndex < len(ft.responses) {
		resp := ft.responses[ft.respIndex]
		ft.respIndex++

		return &http.Response{
			StatusCode: resp.StatusCode,
			Header:     resp.Headers,
			Body:       io.NopCloser(bytes.NewReader(resp.Body)),
			Request:    req,
		}, nil
	}

	// Default response if nothing queued
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
		Request:    req,
	}, nil
}

// PushResponse queues a response to be returned.
func (ft *FakeTransport) PushResponse(statusCode int, body interface{}) {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	ft.responses = append(ft.responses, TransportResponse{
		StatusCode: statusCode,
		Body:       bodyBytes,
	})
}

// PushError queues an error to be returned.
func (ft *FakeTransport) PushError(err error) {
	ft.mu.Lock()
	defer ft.mu.Unlock()
	ft.errs = append(ft.errs, err)
}

// Requests returns all recorded requests.
func (ft *FakeTransport) Requests() []TransportRequest {
	ft.mu.Lock()
	defer ft.mu.Unlock()
	return ft.requests
}

// Reset clears all queued responses and recorded requests.
func (ft *FakeTransport) Reset() {
	ft.mu.Lock()
	defer ft.mu.Unlock()
	ft.responses = make([]TransportResponse, 0)
	ft.requests = make([]TransportRequest, 0)
	ft.errs = make([]error, 0)
	ft.errIndex = 0
	ft.respIndex = 0
}

// createTestQredex creates a Qredex instance with a FakeTransport for testing.
func createTestQredex(transport *FakeTransport) (*Qredex, error) {
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	cfg := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		BaseURL:      "https://api.qredex.com",
		HTTPClient:   httpClient,
		Timeout:      10 * time.Second,
	}

	return New(cfg)
}

// TestCanonicalFlow tests the complete IIT -> PIT -> paid order -> refund flow.
func TestCanonicalFlow(t *testing.T) {
	transport := NewFakeTransport()

	// Queue OAuth token response
	transport.PushResponse(http.StatusOK, map[string]interface{}{
		"access_token": "test-access-token",
		"token_type":   "Bearer",
		"expires_in":   3600,
	})

	// Queue creator response
	transport.PushResponse(http.StatusOK, map[string]interface{}{
		"id":           "creator-123",
		"handle":       "alice",
		"status":       "ACTIVE",
		"display_name": "Alice",
		"created_at":   "2026-01-01T00:00:00Z",
		"updated_at":   "2026-01-01T00:00:00Z",
	})

	// Queue link response
	transport.PushResponse(http.StatusOK, map[string]interface{}{
		"id":                      "link-123",
		"merchant_id":             "merchant-123",
		"store_id":                "store-123",
		"creator_id":              "creator-123",
		"link_name":               "spring-launch",
		"link_code":               "abc123",
		"public_link_url":         "https://qredex.com/l/abc123",
		"destination_path":        "/products/spring",
		"status":                  "ACTIVE",
		"attribution_window_days": 30,
		"created_at":              "2026-01-01T00:00:00Z",
		"updated_at":              "2026-01-01T00:00:00Z",
	})

	// Queue IIT response
	transport.PushResponse(http.StatusOK, map[string]interface{}{
		"id":                "iit-123",
		"merchant_id":       "merchant-123",
		"link_id":           "link-123",
		"token":             "eyJhbGc...iit-token",
		"token_id":          "token-123",
		"issued_at":         "2026-01-01T00:00:00Z",
		"expires_at":        "2026-01-02T00:00:00Z",
		"status":            "ACTIVE",
		"integrity_version": 1,
	})

	// Queue PIT response
	transport.PushResponse(http.StatusOK, map[string]interface{}{
		"id":                      "pit-123",
		"merchant_id":             "merchant-123",
		"store_id":                "store-123",
		"link_id":                 "link-123",
		"token":                   "eyJhbGc...pit-token",
		"token_id":                "token-456",
		"source":                  "backend-cart",
		"origin_match_status":     "MATCH",
		"window_status":           "WITHIN",
		"attribution_window_days": 30,
		"store_domain_snapshot":   "example.com",
		"issued_at":               "2026-01-01T00:00:00Z",
		"expires_at":              "2026-01-02T00:00:00Z",
		"locked_at":               "2026-01-01T00:01:00Z",
		"integrity_version":       1,
		"eligible":                true,
		"created_at":              "2026-01-01T00:00:00Z",
		"updated_at":              "2026-01-01T00:01:00Z",
	})

	// Queue paid order response
	transport.PushResponse(http.StatusOK, map[string]interface{}{
		"id":                    "order-attr-123",
		"merchant_id":           "merchant-123",
		"order_source":          "DIRECT_API",
		"external_order_id":     "order-100045",
		"order_number":          "100045",
		"paid_at":               "2026-01-01T00:02:00Z",
		"currency":              "USD",
		"total_price":           110.00,
		"purchase_intent_token": "eyJhbGc...pit-token",
		"link_id":               "link-123",
		"link_name":             "spring-launch",
		"link_code":             "abc123",
		"creator_id":            "creator-123",
		"creator_handle":        "alice",
		"duplicate_suspect":     false,
		"integrity_score":       95,
		"integrity_band":        "HIGH",
		"review_required":       false,
		"resolution_status":     "ATTRIBUTED",
		"token_integrity":       "VALID",
		"window_status":         "WITHIN",
		"origin_match_status":   "MATCH",
		"created_at":            "2026-01-01T00:02:00Z",
		"updated_at":            "2026-01-01T00:02:00Z",
	})

	// Queue refund response
	transport.PushResponse(http.StatusOK, map[string]interface{}{
		"id":                "order-attr-123",
		"merchant_id":       "merchant-123",
		"order_source":      "DIRECT_API",
		"external_order_id": "order-100045",
		"order_number":      "100045",
		"paid_at":           "2026-01-01T00:02:00Z",
		"currency":          "USD",
		"total_price":       110.00,
		"duplicate_suspect": false,
		"integrity_score":   95,
		"integrity_band":    "HIGH",
		"review_required":   false,
		"resolution_status": "ATTRIBUTED",
		"created_at":        "2026-01-01T00:02:00Z",
		"updated_at":        "2026-01-01T00:03:00Z",
	})

	qredex, err := createTestQredex(transport)
	if err != nil {
		t.Fatalf("Failed to create Qredex instance: %v", err)
	}

	ctx := context.Background()

	// Step 1: Create creator
	creator, err := qredex.Creators().Create(ctx, CreateCreatorRequest{
		Handle:      "alice",
		DisplayName: strPtr("Alice"),
	})
	if err != nil {
		t.Fatalf("Create creator failed: %v", err)
	}
	if creator.ID != "creator-123" {
		t.Errorf("Expected creator ID 'creator-123', got %q", creator.ID)
	}

	// Step 2: Create link
	link, err := qredex.Links().Create(ctx, CreateLinkRequest{
		StoreID:         "store-123",
		CreatorID:       creator.ID,
		LinkName:        "spring-launch",
		DestinationPath: "/products/spring",
	})
	if err != nil {
		t.Fatalf("Create link failed: %v", err)
	}
	if link.ID != "link-123" {
		t.Errorf("Expected link ID 'link-123', got %q", link.ID)
	}

	// Step 3: Issue IIT
	iit, err := qredex.Intents().IssueInfluenceIntentToken(ctx, IssueInfluenceIntentTokenRequest{
		LinkID: link.ID,
	})
	if err != nil {
		t.Fatalf("Issue IIT failed: %v", err)
	}
	if iit.Token != "eyJhbGc...iit-token" {
		t.Errorf("Expected IIT token, got %q", iit.Token)
	}

	// Step 4: Lock PIT
	pit, err := qredex.Intents().LockPurchaseIntent(ctx, LockPurchaseIntentRequest{
		Token:  iit.Token,
		Source: strPtr("backend-cart"),
	})
	if err != nil {
		t.Fatalf("Lock PIT failed: %v", err)
	}
	if pit.Token != "eyJhbGc...pit-token" {
		t.Errorf("Expected PIT token, got %q", pit.Token)
	}

	// Step 5: Record paid order
	order, err := qredex.Orders().RecordPaidOrder(ctx, RecordPaidOrderRequest{
		StoreID:             "store-123",
		ExternalOrderID:     "order-100045",
		Currency:            "USD",
		TotalPrice:          floatPtr(110.00),
		PurchaseIntentToken: strPtr(pit.Token),
	})
	if err != nil {
		t.Fatalf("Record paid order failed: %v", err)
	}
	if order.ResolutionStatus != ResolutionStatusAttributed {
		t.Errorf("Expected resolution status 'ATTRIBUTED', got %q", order.ResolutionStatus)
	}
	if order.TokenIntegrity == nil || *order.TokenIntegrity != TokenIntegrityValid {
		t.Errorf("Expected token integrity 'VALID'")
	}

	// Step 6: Record refund
	refund, err := qredex.Refunds().RecordRefund(ctx, RecordRefundRequest{
		StoreID:          "store-123",
		ExternalOrderID:  "order-100045",
		ExternalRefundID: "refund-100045-1",
		RefundTotal:      floatPtr(25.00),
	})
	if err != nil {
		t.Fatalf("Record refund failed: %v", err)
	}
	if refund.ID != "order-attr-123" {
		t.Errorf("Expected order attribution ID 'order-attr-123', got %q", refund.ID)
	}

	// Verify all requests were made (6 resource requests + 1 token request = 7)
	requests := transport.Requests()
	if len(requests) != 7 {
		t.Errorf("Expected 7 requests (6 resource + 1 token), got %d", len(requests))
	}
}

// TestErrorHandling tests various error scenarios.
func TestErrorHandling(t *testing.T) {
	t.Run("AuthenticationError", func(t *testing.T) {
		transport := NewFakeTransport()
		// Token request succeeds
		transport.PushResponse(http.StatusOK, map[string]interface{}{
			"access_token": "test-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
		// Actual API call fails with 401
		transport.PushResponse(http.StatusUnauthorized, map[string]interface{}{
			"error_code": "invalid_token",
			"message":    "The access token is invalid or expired",
		})

		qredex, _ := createTestQredex(transport)
		ctx := context.Background()

		_, err := qredex.Creators().Create(ctx, CreateCreatorRequest{
			Handle: "test",
		})

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !IsAuthenticationError(err) {
			t.Errorf("Expected AuthenticationError, got %T", err)
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		transport := NewFakeTransport()
		// First call is token request
		transport.PushResponse(http.StatusOK, map[string]interface{}{
			"access_token": "test-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
		// Second call is the actual request
		transport.PushResponse(http.StatusBadRequest, map[string]interface{}{
			"error_code": "validation_error",
			"message":    "Invalid request: handle is required",
		})

		qredex, _ := createTestQredex(transport)
		ctx := context.Background()

		_, err := qredex.Creators().Create(ctx, CreateCreatorRequest{})

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !IsValidationError(err) {
			t.Errorf("Expected ValidationError, got %T", err)
		}
	})

	t.Run("ConflictError", func(t *testing.T) {
		transport := NewFakeTransport()
		transport.PushResponse(http.StatusOK, map[string]interface{}{
			"access_token": "test-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
		transport.PushResponse(http.StatusConflict, map[string]interface{}{
			"error_code": "REJECTED_CROSS_SOURCE_DUPLICATE",
			"message":    "Order already exists from different source",
		})

		qredex, _ := createTestQredex(transport)
		ctx := context.Background()

		_, err := qredex.Orders().RecordPaidOrder(ctx, RecordPaidOrderRequest{
			StoreID:         "store-123",
			ExternalOrderID: "order-123",
			Currency:        "USD",
			TotalPrice:      floatPtr(100.00),
		})

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !IsConflictError(err) {
			t.Errorf("Expected ConflictError, got %T", err)
		}
	})

	t.Run("NotFoundError", func(t *testing.T) {
		transport := NewFakeTransport()
		transport.PushResponse(http.StatusOK, map[string]interface{}{
			"access_token": "test-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
		transport.PushResponse(http.StatusNotFound, map[string]interface{}{
			"error_code": "not_found",
			"message":    "Creator not found",
		})

		qredex, _ := createTestQredex(transport)
		ctx := context.Background()

		_, err := qredex.Creators().Get(ctx, "non-existent-id")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !IsNotFoundError(err) {
			t.Errorf("Expected NotFoundError, got %T", err)
		}
	})

	t.Run("RateLimitError", func(t *testing.T) {
		transport := NewFakeTransport()
		transport.PushResponse(http.StatusOK, map[string]interface{}{
			"access_token": "test-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})

		headers := make(http.Header)
		headers.Set("Retry-After", "60")
		transport.responses = append(transport.responses, TransportResponse{
			StatusCode: http.StatusTooManyRequests,
			Headers:    headers,
			Body:       []byte(`{"error_code":"rate_limited","message":"Rate limit exceeded"}`),
		})

		qredex, _ := createTestQredex(transport)
		ctx := context.Background()

		_, err := qredex.Creators().List(ctx, ListCreatorsRequest{})

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !IsRateLimitError(err) {
			t.Errorf("Expected RateLimitError, got %T", err)
		}

		rateLimitErr, ok := err.(*RateLimitError)
		if ok && rateLimitErr.RetryAfterSeconds != 60 {
			t.Errorf("Expected RetryAfterSeconds=60, got %d", rateLimitErr.RetryAfterSeconds)
		}
	})
}

// TestTokenCaching tests that tokens are cached and reused.
func TestTokenCaching(t *testing.T) {
	transport := NewFakeTransport()

	// Queue two token responses (should only use first)
	transport.PushResponse(http.StatusOK, map[string]interface{}{
		"access_token": "cached-token",
		"token_type":   "Bearer",
		"expires_in":   3600,
	})

	// Queue two successful creator responses
	transport.PushResponse(http.StatusOK, map[string]interface{}{
		"id":         "creator-1",
		"handle":     "alice",
		"status":     "ACTIVE",
		"created_at": "2026-01-01T00:00:00Z",
		"updated_at": "2026-01-01T00:00:00Z",
	})
	transport.PushResponse(http.StatusOK, map[string]interface{}{
		"id":         "creator-2",
		"handle":     "bob",
		"status":     "ACTIVE",
		"created_at": "2026-01-01T00:00:00Z",
		"updated_at": "2026-01-01T00:00:00Z",
	})

	qredex, err := createTestQredex(transport)
	if err != nil {
		t.Fatalf("Failed to create Qredex: %v", err)
	}

	ctx := context.Background()

	// Make two requests
	_, err = qredex.Creators().Get(ctx, "creator-1")
	if err != nil {
		t.Fatalf("First Get failed: %v", err)
	}

	_, err = qredex.Creators().Get(ctx, "creator-2")
	if err != nil {
		t.Fatalf("Second Get failed: %v", err)
	}

	// Verify only one token request was made
	requests := transport.Requests()
	tokenRequests := 0
	for _, req := range requests {
		if req.URL == "https://api.qredex.com/api/v1/auth/token" {
			tokenRequests++
		}
	}

	if tokenRequests != 1 {
		t.Errorf("Expected 1 token request, got %d", tokenRequests)
	}
}

// TestModelSerialization tests JSON serialization/deserialization.
func TestModelSerialization(t *testing.T) {
	t.Run("Creator", func(t *testing.T) {
		jsonData := `{
			"id": "creator-123",
			"handle": "alice",
			"status": "ACTIVE",
			"display_name": "Alice",
			"email": "alice@example.com",
			"socials": {"twitter": "@alice"},
			"created_at": "2026-01-01T00:00:00Z",
			"updated_at": "2026-01-01T00:00:00Z"
		}`

		var creator Creator
		err := json.Unmarshal([]byte(jsonData), &creator)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if creator.ID != "creator-123" {
			t.Errorf("Expected ID 'creator-123', got %q", creator.ID)
		}
		if creator.Handle != "alice" {
			t.Errorf("Expected handle 'alice', got %q", creator.Handle)
		}
		if creator.Status != CreatorStatusActive {
			t.Errorf("Expected status ACTIVE, got %q", creator.Status)
		}
		if creator.DisplayName == nil || *creator.DisplayName != "Alice" {
			t.Errorf("Expected display name 'Alice'")
		}
	})

	t.Run("OrderAttribution", func(t *testing.T) {
		jsonData := `{
			"id": "order-attr-123",
			"merchant_id": "merchant-123",
			"order_source": "DIRECT_API",
			"external_order_id": "order-100045",
			"order_number": "100045",
			"paid_at": "2026-01-01T00:00:00Z",
			"currency": "USD",
			"total_price": 110.00,
			"link_id": "link-123",
			"link_name": "spring-launch",
			"link_code": "abc123",
			"creator_id": "creator-123",
			"creator_handle": "alice",
			"duplicate_suspect": false,
			"integrity_score": 95,
			"integrity_band": "HIGH",
			"review_required": false,
			"resolution_status": "ATTRIBUTED",
			"token_integrity": "VALID",
			"window_status": "WITHIN",
			"origin_match_status": "MATCH",
			"created_at": "2026-01-01T00:00:00Z",
			"updated_at": "2026-01-01T00:00:00Z"
		}`

		var order OrderAttribution
		err := json.Unmarshal([]byte(jsonData), &order)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if order.ResolutionStatus != ResolutionStatusAttributed {
			t.Errorf("Expected resolution status ATTRIBUTED, got %q", order.ResolutionStatus)
		}
		if order.TokenIntegrity == nil || *order.TokenIntegrity != TokenIntegrityValid {
			t.Errorf("Expected token integrity VALID")
		}
		if order.IntegrityScore != 95 {
			t.Errorf("Expected integrity score 95, got %d", order.IntegrityScore)
		}
	})

	t.Run("LinkStats", func(t *testing.T) {
		jsonData := `{
			"link_id": "link-123",
			"clicks_count": 150,
			"sessions_count": 120,
			"orders_count": 25,
			"revenue_total": 2750.50,
			"token_invalid_count": 5,
			"token_missing_count": 10,
			"last_click_at": "2026-01-01T00:00:00Z",
			"last_order_at": "2026-01-01T00:00:00Z"
		}`

		var stats LinkStats
		err := json.Unmarshal([]byte(jsonData), &stats)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if stats.ClicksCount != 150 {
			t.Errorf("Expected 150 clicks, got %d", stats.ClicksCount)
		}
		if stats.RevenueTotal != 2750.50 {
			t.Errorf("Expected revenue 2750.50, got %f", stats.RevenueTotal)
		}
	})
}

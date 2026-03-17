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
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBootstrap_Valid(t *testing.T) {
	t.Setenv("QREDEX_CLIENT_ID", "test-id")
	t.Setenv("QREDEX_CLIENT_SECRET", "test-secret")

	q, err := Bootstrap()
	if err != nil {
		t.Fatalf("Bootstrap failed: %v", err)
	}
	if q == nil {
		t.Fatal("Expected non-nil Qredex instance")
	}
}

func TestBootstrap_Missing_ClientID(t *testing.T) {
	t.Setenv("QREDEX_CLIENT_ID", "")
	t.Setenv("QREDEX_CLIENT_SECRET", "test-secret")

	_, err := Bootstrap()
	if err == nil {
		t.Fatal("Expected error for missing CLIENT_ID")
	}
	if _, ok := err.(*ConfigurationError); !ok {
		t.Fatalf("Expected ConfigurationError, got %T", err)
	}
}

func TestNew_Valid(t *testing.T) {
	cfg := Config{
		ClientID:     "test-id",
		ClientSecret: "test-secret",
	}

	q, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if q == nil {
		t.Fatal("Expected non-nil Qredex instance")
	}
}

func TestNew_Missing_ClientID(t *testing.T) {
	cfg := Config{
		ClientSecret: "test-secret",
	}

	_, err := New(cfg)
	if err == nil {
		t.Fatal("Expected error for missing ClientID")
	}
	if _, ok := err.(*ConfigurationError); !ok {
		t.Fatalf("Expected ConfigurationError, got %T", err)
	}
}

func TestConfig_ResolvedBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantURL string
	}{
		{
			name:    "production",
			cfg:     Config{Environment: Production},
			wantURL: "https://api.qredex.com",
		},
		{
			name:    "staging",
			cfg:     Config{Environment: Staging},
			wantURL: "https://staging-api.qredex.com",
		},
		{
			name:    "development",
			cfg:     Config{Environment: Development},
			wantURL: "http://localhost:8080",
		},
		{
			name:    "custom baseurl",
			cfg:     Config{BaseURL: "https://custom.example.com/"},
			wantURL: "https://custom.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.resolvedBaseURL(); got != tt.wantURL {
				t.Errorf("resolvedBaseURL() = %s, want %s", got, tt.wantURL)
			}
		})
	}
}

func TestAPIError_Messages(t *testing.T) {
	err := &APIError{
		Status:    400,
		ErrorCode: "invalid_request",
		Message:   "invalid link_id",
	}

	errStr := err.Error()
	if errStr != "qredex: API error 400 [invalid_request]: invalid link_id" {
		t.Errorf("APIError.Error() = %q, unexpected", errStr)
	}
}

func TestIsAuthenticationError(t *testing.T) {
	err := &AuthenticationError{
		APIError: APIError{
			Status:  401,
			Message: "Unauthorized",
		},
	}

	if !IsAuthenticationError(err) {
		t.Fatal("IsAuthenticationError failed")
	}
	if IsAuthorizationError(err) {
		t.Fatal("IsAuthorizationError should be false")
	}
}

func TestTokenProvider_IssueToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "test-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`))
	}))
	defer server.Close()

	cfg := Config{
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		BaseURL:      server.URL,
	}

	tp := newTokenProvider(&cfg, server.Client())
	token, err := tp.issueToken(context.Background())
	if err != nil {
		t.Fatalf("issueToken failed: %v", err)
	}
	if token.AccessToken != "test-token" {
		t.Errorf("AccessToken = %s, want test-token", token.AccessToken)
	}
	if token.ExpiresIn != 3600 {
		t.Errorf("ExpiresIn = %d, want 3600", token.ExpiresIn)
	}
}

func TestTokenCache(t *testing.T) {
	cache := &tokenCache{}

	// Cache should be nil initially
	if cached := cache.get(); cached != nil {
		t.Fatal("cache should be nil initially")
	}

	// Set a token
	cache.set("test-token", 3600)
	cached := cache.get()
	if cached == nil {
		t.Fatal("cache.get() should not be nil after set")
	}
	if cached.accessToken != "test-token" {
		t.Errorf("accessToken = %s, want test-token", cached.accessToken)
	}

	// Clear the cache
	cache.clear()
	if cached := cache.get(); cached != nil {
		t.Fatal("cache should be nil after clear")
	}
}

func TestBackoffDelay(t *testing.T) {
	base := 100 * time.Millisecond
	max := 10 * time.Second

	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{0, 100 * time.Millisecond},
		{1, 200 * time.Millisecond},
		{2, 400 * time.Millisecond},
		{3, 800 * time.Millisecond},
		{5, 3200 * time.Millisecond},
		{10, 10 * time.Second}, // capped at max
	}

	for _, tt := range tests {
		t.Run("attempt_"+string(rune(tt.attempt+'0')), func(t *testing.T) {
			if got := backoffDelay(tt.attempt, base, max); got != tt.want {
				t.Errorf("backoffDelay(%d) = %v, want %v", tt.attempt, got, tt.want)
			}
		})
	}
}

func TestHTTPClient_Request_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("missing Authorization header")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "test-id", "handle": "alice"}`))
	}))
	defer server.Close()

	cfg := Config{
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		BaseURL:      server.URL,
	}

	q, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	var result Creator
	err = q.hc.request(context.Background(), "GET", "/test", nil, &result)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

// TestGETQueryParams verifies that List* request structs are serialised as
// URL query parameters rather than a JSON request body for GET operations.
func TestGETQueryParams(t *testing.T) {
	var capturedURL string
	var capturedBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"t","token_type":"Bearer","expires_in":3600}`))
			return
		}
		capturedURL = r.URL.RequestURI()
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"items":[],"page":1,"size":10,"total_elements":0,"total_pages":0}`))
	}))
	defer server.Close()

	cfg := Config{ClientID: "id", ClientSecret: "sec", BaseURL: server.URL}
	q, _ := New(cfg)

	page := intPtr(2)
	size := intPtr(5)
	status := CreatorStatusActive
	_, err := q.Creators().List(context.Background(), ListCreatorsRequest{
		Page:   page,
		Size:   size,
		Status: &status,
	})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(capturedBody) != 0 {
		t.Errorf("GET request should have no body, got %q", string(capturedBody))
	}
	if !strings.Contains(capturedURL, "page=2") {
		t.Errorf("expected query param page=2 in %q", capturedURL)
	}
	if !strings.Contains(capturedURL, "size=5") {
		t.Errorf("expected query param size=5 in %q", capturedURL)
	}
	if !strings.Contains(capturedURL, "status=ACTIVE") {
		t.Errorf("expected query param status=ACTIVE in %q", capturedURL)
	}
}

// TestUserAgentSuffix verifies that Config.UserAgentSuffix is appended to
// the User-Agent header on outgoing resource requests.
func TestUserAgentSuffix(t *testing.T) {
	var capturedUA string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"t","token_type":"Bearer","expires_in":3600}`))
			return
		}
		capturedUA = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"items":[],"page":1,"size":10,"total_elements":0,"total_pages":0}`))
	}))
	defer server.Close()

	cfg := Config{
		ClientID:        "id",
		ClientSecret:    "sec",
		BaseURL:         server.URL,
		UserAgentSuffix: "my-platform/1.0",
	}
	q, _ := New(cfg)

	q.Creators().List(context.Background(), ListCreatorsRequest{})

	if !strings.HasPrefix(capturedUA, "qredex-go/") {
		t.Errorf("User-Agent should start with qredex-go/, got %q", capturedUA)
	}
	if !strings.Contains(capturedUA, "my-platform/1.0") {
		t.Errorf("User-Agent should contain suffix 'my-platform/1.0', got %q", capturedUA)
	}
}

// TestRetryOn5xx verifies that GET requests are retried on 5xx responses
// when RetryMax > 0.
func TestRetryOn5xx(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"t","token_type":"Bearer","expires_in":3600}`))
			return
		}
		callCount++
		if callCount < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error_code":"service_unavailable","message":"try again"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"items":[],"page":1,"size":10,"total_elements":0,"total_pages":0}`))
	}))
	defer server.Close()

	cfg := Config{
		ClientID:       "id",
		ClientSecret:   "sec",
		BaseURL:        server.URL,
		RetryMax:       3,
		RetryBaseDelay: 1 * time.Millisecond,
		RetryMaxDelay:  10 * time.Millisecond,
	}
	q, _ := New(cfg)

	_, err := q.Creators().List(context.Background(), ListCreatorsRequest{})
	if err != nil {
		t.Fatalf("List should succeed after retries, got: %v", err)
	}
	if callCount != 3 {
		t.Errorf("expected 3 resource calls (2 failures + 1 success), got %d", callCount)
	}
}

// TestNetworkError verifies that a connection failure produces a NetworkError.
func TestNetworkError(t *testing.T) {
	transport := NewFakeTransport()
	// Token succeeds
	transport.PushResponse(http.StatusOK, map[string]interface{}{
		"access_token": "t",
		"token_type":   "Bearer",
		"expires_in":   3600,
	})
	// Resource call fails at transport level
	transport.PushError(errors.New("dial tcp: connection refused"))

	qredex, _ := createTestQredex(transport)

	_, err := qredex.Creators().Get(context.Background(), "creator-xyz")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNetworkError(err) {
		t.Errorf("expected NetworkError, got %T: %v", err, err)
	}
}

// TestStructToQueryParams verifies the reflection-based query param encoder.
func TestStructToQueryParams(t *testing.T) {
	t.Run("nil body", func(t *testing.T) {
		params := structToQueryParams(nil)
		if len(params) != 0 {
			t.Errorf("expected empty params, got %v", params)
		}
	})

	t.Run("list request with all fields set", func(t *testing.T) {
		page := 2
		size := 25
		params := structToQueryParams(ListOrdersRequest{Page: &page, Size: &size})
		if params.Get("page") != "2" {
			t.Errorf("page = %q, want 2", params.Get("page"))
		}
		if params.Get("size") != "25" {
			t.Errorf("size = %q, want 25", params.Get("size"))
		}
	})

	t.Run("nil pointer fields are omitted", func(t *testing.T) {
		params := structToQueryParams(ListCreatorsRequest{})
		if len(params) != 0 {
			t.Errorf("expected empty params for zero-value request, got %v", params)
		}
	})

	t.Run("string enum field", func(t *testing.T) {
		status := CreatorStatusActive
		params := structToQueryParams(ListCreatorsRequest{Status: &status})
		if params.Get("status") != "ACTIVE" {
			t.Errorf("status = %q, want ACTIVE", params.Get("status"))
		}
	})
}

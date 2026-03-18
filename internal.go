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
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// tokenCache holds a cached OAuth access token with thread-safe access.
type tokenCache struct {
	mu    sync.RWMutex
	token *cachedToken
}

type cachedToken struct {
	accessToken string
	expiresAt   time.Time
}

func (tc *tokenCache) get() *cachedToken {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if tc.token == nil || time.Now().Add(30*time.Second).After(tc.token.expiresAt) {
		return nil
	}
	return tc.token
}

func (tc *tokenCache) set(accessToken string, expiresIn int) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.token = &cachedToken{
		accessToken: accessToken,
		expiresAt:   time.Now().Add(time.Duration(expiresIn) * time.Second),
	}
}

func (tc *tokenCache) clear() {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.token = nil
}

// tokenProvider manages OAuth 2.0 client-credentials token acquisition and caching.
type tokenProvider struct {
	config *Config
	cache  *tokenCache
	client *http.Client
}

func newTokenProvider(config *Config, httpClient *http.Client) *tokenProvider {
	return &tokenProvider{
		config: config,
		cache:  &tokenCache{},
		client: httpClient,
	}
}

// issueToken fetches a fresh OAuth token from the /api/v1/auth/token endpoint.
func (tp *tokenProvider) issueToken(ctx context.Context) (*OAuthTokenResponse, error) {
	baseURL := tp.config.resolvedBaseURL()
	tokenURL := baseURL + "/api/v1/auth/token"

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	if scope := tp.config.scopeString(); scope != "" {
		form.Set("scope", scope)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, &NetworkError{Message: "failed to create token request", Cause: err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent(tp.config.UserAgentSuffix))
	req.SetBasicAuth(tp.config.ClientID, tp.config.ClientSecret)

	ctx, cancel := context.WithTimeout(ctx, tp.config.resolvedTimeout())
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := tp.client.Do(req)
	if err != nil {
		return nil, &NetworkError{Message: "token request failed", Cause: err}
	}
	defer func() { _ = resp.Body.Close() }() // explicitly ignore error per linter

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
	}

	var result OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, &NetworkError{Message: "failed to parse token response", Cause: err}
	}

	return &result, nil
}

// getToken returns a valid cached token or fetches a new one.
func (tp *tokenProvider) getToken(ctx context.Context) (string, error) {
	cached := tp.cache.get()
	if cached != nil {
		return cached.accessToken, nil
	}

	token, err := tp.issueToken(ctx)
	if err != nil {
		return "", err
	}

	tp.cache.set(token.AccessToken, token.ExpiresIn)
	return token.AccessToken, nil
}

// httpClient manages HTTP requests with auth, retries, and error handling.
// httpClient manages HTTP requests with auth, retries, and error handling.
type httpClient struct {
	baseURL                string
	httpClient             *http.Client
	tokenProvider          *tokenProvider
	userAgentSuffix        string
	timeout                time.Duration
	retryMax               int
	retryBaseDelay         time.Duration
	retryMaxDelay          time.Duration
	logger                 Logger
	tracer                 Tracer
	metrics                Metrics
	idempotencyKeyProvider IdempotencyKeyProvider
}

func newHTTPClient(config *Config, httpCli *http.Client, tp *tokenProvider) *httpClient {
	return &httpClient{
		baseURL:                config.resolvedBaseURL(),
		httpClient:             httpCli,
		tokenProvider:          tp,
		userAgentSuffix:        config.UserAgentSuffix,
		timeout:                config.resolvedTimeout(),
		retryMax:               config.resolvedRetryMax(),
		retryBaseDelay:         config.resolvedRetryBaseDelay(),
		retryMaxDelay:          config.resolvedRetryMaxDelay(),
		logger:                 config.Logger,
		tracer:                 config.Tracer,
		metrics:                config.Metrics,
		idempotencyKeyProvider: config.IdempotencyKeyProvider,
	}
}

// request executes an HTTP request with automatic auth, retries on safe operations, and error handling.
func (hc *httpClient) request(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	return hc.doRequest(ctx, method, path, body, result, 0)
}

func (hc *httpClient) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}, attempt int) error {
	token, err := hc.tokenProvider.getToken(ctx)
	if err != nil {
		return err
	}

	rawURL := hc.baseURL + path
	startTime := time.Now()
	var reqBody io.Reader
	isReadMethod := method == "GET" || method == "HEAD"

	if body != nil && isReadMethod {
		// For safe/idempotent methods, encode the body struct as URL query parameters.
		params := structToQueryParams(body)
		if len(params) > 0 {
			rawURL += "?" + params.Encode()
		}
	} else if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return &NetworkError{Message: "failed to marshal request body", Cause: err}
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, reqBody)
	if err != nil {
		return &NetworkError{Message: "failed to create request", Cause: err}
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", userAgent(hc.userAgentSuffix))
	if !isReadMethod {
		req.Header.Set("Content-Type", "application/json")
		if hc.idempotencyKeyProvider != nil {
			key := hc.idempotencyKeyProvider.GetIdempotencyKey(ctx, method, path, body)
			if key != "" {
				req.Header.Set("Idempotency-Key", key)
			}
		}
	}
	if hc.logger != nil {
		hc.logger.Printf("[qredex] %s %s", method, rawURL)
	}
	if hc.tracer != nil {
		hc.tracer.Trace("qredex.request", map[string]interface{}{
			"method":  method,
			"url":     rawURL,
			"attempt": attempt,
		})
	}
	if hc.metrics != nil {
		hc.metrics.Record("qredex.request.count", 1, map[string]string{
			"method": method,
			"path":   path,
		})
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, hc.timeout)
	defer cancel()
	req = req.WithContext(ctxWithTimeout)

	resp, err := hc.httpClient.Do(req)
	if err != nil {
		if attempt < hc.retryMax && isReadMethod {
			delay := backoffDelay(attempt, hc.retryBaseDelay, hc.retryMaxDelay)
			select {
			case <-time.After(delay):
				return hc.doRequest(ctx, method, path, body, result, attempt+1)
			case <-ctx.Done():
				return &NetworkError{Message: "context canceled", Cause: ctx.Err()}
			}
		}
		return &NetworkError{Message: "HTTP request failed", Cause: err}
	}
	defer func() { _ = resp.Body.Close() }() // explicitly ignore error per linter

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &NetworkError{Message: "failed to read response body", Cause: err}
	}

	if resp.StatusCode >= 400 {
		if hc.metrics != nil {
			hc.metrics.Record("qredex.request.error", 1, map[string]string{
				"method": method,
				"path":   path,
				"status": strconv.Itoa(resp.StatusCode),
			})
		}
		// Retry safe methods on 429 or 5xx responses.
		if attempt < hc.retryMax && isReadMethod && isRetryableStatus(resp.StatusCode) {
			delay := backoffDelay(attempt, hc.retryBaseDelay, hc.retryMaxDelay)
			if resp.StatusCode == http.StatusTooManyRequests {
				if ra := resp.Header.Get("Retry-After"); ra != "" {
					if secs, parseErr := strconv.Atoi(ra); parseErr == nil && secs > 0 {
						d := time.Duration(secs) * time.Second
						if d < hc.retryMaxDelay {
							delay = d
						} else {
							delay = hc.retryMaxDelay
						}
					}
				}
			}
			select {
			case <-time.After(delay):
				return hc.doRequest(ctx, method, path, body, result, attempt+1)
			case <-ctx.Done():
				return &NetworkError{Message: "context canceled", Cause: ctx.Err()}
			}
		}
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
		return parseAPIError(resp)
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return &NetworkError{Message: "failed to parse response", Cause: err}
		}
	}

	if hc.metrics != nil {
		latency := float64(time.Since(startTime).Milliseconds())
		hc.metrics.Record("qredex.request.latency_ms", latency, map[string]string{
			"method": method,
			"path":   path,
		})
	}
	return nil
}

func userAgent(suffix string) string {
	ua := "qredex-go/" + SDKVersion
	if suffix != "" {
		ua += " " + suffix
	}
	return ua
}

// isRetryableStatus reports whether an HTTP status code warrants a retry.
// Retries are safe on 429 (Too Many Requests) and any 5xx server error.
func isRetryableStatus(status int) bool {
	return status == http.StatusTooManyRequests || (status >= 500 && status <= 599)
}

// structToQueryParams converts a struct to url.Values using its json struct tags.
// Only exported fields with a non-empty json tag are included.
// Nil pointer fields and zero-string fields tagged omitempty are skipped.
func structToQueryParams(v interface{}) url.Values {
	params := url.Values{}
	if v == nil {
		return params
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return params
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return params
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		parts := strings.SplitN(tag, ",", 2)
		name := parts[0]
		if name == "" || name == "-" {
			continue
		}
		omitempty := len(parts) > 1 && strings.Contains(parts[1], "omitempty")

		// Dereference pointers; skip nil.
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				continue
			}
			fieldVal = fieldVal.Elem()
		}

		var s string
		switch fieldVal.Kind() {
		case reflect.String:
			s = fieldVal.String()
			if omitempty && s == "" {
				continue
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v := fieldVal.Int()
			if omitempty && v == 0 {
				continue
			}
			s = strconv.FormatInt(v, 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v := fieldVal.Uint()
			if omitempty && v == 0 {
				continue
			}
			s = strconv.FormatUint(v, 10)
		case reflect.Float32, reflect.Float64:
			s = strconv.FormatFloat(fieldVal.Float(), 'f', -1, 64)
		case reflect.Bool:
			b := fieldVal.Bool()
			if omitempty && !b {
				continue
			}
			s = strconv.FormatBool(b)
		default:
			continue
		}

		params.Set(name, s)
	}

	return params
}

func backoffDelay(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	delay := baseDelay
	for i := 0; i < attempt; i++ {
		delay *= 2
	}
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}

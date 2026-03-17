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

/*
Package qredexsdk is the official Qredex Go server SDK for machine-to-machine integrations.

It covers the Qredex Integrations API: creators, links, intents (IIT/PIT), orders, and refunds.
Auth is handled automatically via the OAuth 2.0 client credentials flow.

Quick start:

	qredex, err := qredexsdk.Bootstrap()
	if err != nil {
	    log.Fatal(err)
	}

	creator, err := qredex.Creators().Create(ctx, qredexsdk.CreateCreatorRequest{
	    Handle:      "alice",
	    DisplayName: ptr("Alice"),
	})

See README.md and INTEGRATION_GUIDE.md for full documentation.
*/
package qredex

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// SDKVersion is the current release version of the Qredex Go SDK.
const SDKVersion = "0.1.0"

// Environment identifies the Qredex API environment.
type Environment string

const (
	// Production is the live Qredex API environment.
	Production Environment = "production"
	// Staging is the Qredex staging environment.
	Staging Environment = "staging"
	// Development is a local development environment.
	Development Environment = "development"
)

// baseURL returns the base URL for the environment.
func (e Environment) baseURL() string {
	switch e {
	case Staging:
		return "https://staging-api.qredex.com"
	case Development:
		return "http://localhost:8080"
	default:
		return "https://api.qredex.com"
	}
}

// Scope is a Qredex OAuth scope string for the Integrations API.
type Scope string

const (
	// ScopeAPI grants access to the full direct API surface.
	ScopeAPI Scope = "direct:api"
	// ScopeLinksRead grants read access to links.
	ScopeLinksRead Scope = "direct:links:read"
	// ScopeLinksWrite grants write access to links.
	ScopeLinksWrite Scope = "direct:links:write"
	// ScopeCreatorsRead grants read access to creators.
	ScopeCreatorsRead Scope = "direct:creators:read"
	// ScopeCreatorsWrite grants write access to creators.
	ScopeCreatorsWrite Scope = "direct:creators:write"
	// ScopeOrdersRead grants read access to orders.
	ScopeOrdersRead Scope = "direct:orders:read"
	// ScopeOrdersWrite grants write access to orders.
	ScopeOrdersWrite Scope = "direct:orders:write"
	// ScopeIntentsRead grants read access to intents.
	ScopeIntentsRead Scope = "direct:intents:read"
	// ScopeIntentsWrite grants write access to intents.
	ScopeIntentsWrite Scope = "direct:intents:write"
)

// Logger is a minimal logger interface for SDK observability.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Tracer is a minimal tracer interface for SDK observability.
type Tracer interface {
	Trace(event string, fields map[string]interface{})
}

// Config holds all configuration for the Qredex SDK.
// Build one explicitly with New, or load from environment with Bootstrap.
type Config struct {
	// ClientID is the OAuth client ID issued by Qredex. Required.
	ClientID string
	// ClientSecret is the OAuth client secret. Required. Never logged.
	ClientSecret string
	// Scopes is the optional list of OAuth scopes to request.
	// Defaults to the full direct:api scope when empty.
	Scopes []Scope
	// Environment selects the Qredex API environment. Defaults to Production.
	Environment Environment
	// BaseURL overrides the resolved environment base URL.
	// Use for custom proxy or local mocking.
	BaseURL string
	// HTTPClient is an optional custom HTTP client.
	// If nil, a default client with sensible timeouts is used.
	HTTPClient *http.Client
	// Timeout is the per-request HTTP timeout. Defaults to 10 seconds.
	Timeout time.Duration
	// UserAgentSuffix is appended to the SDK user-agent string.
	UserAgentSuffix string
	// RetryMax is the number of retry attempts for retryable GET requests.
	// Defaults to 0 (no retries). Max effective value is 5.
	RetryMax int
	// RetryBaseDelay is the base delay for exponential backoff. Defaults to 500ms.
	RetryBaseDelay time.Duration
	// RetryMaxDelay caps the retry delay. Defaults to 30s.
	RetryMaxDelay time.Duration

	// Logger is a minimal logger interface for SDK observability.
	Logger Logger
	// Tracer is a minimal tracer interface for SDK observability.
	Tracer Tracer

	// IdempotencyKeyProvider allows injection of idempotency keys for write requests.
	IdempotencyKeyProvider IdempotencyKeyProvider
}

func (c *Config) validate() error {
	if strings.TrimSpace(c.ClientID) == "" {
		return &ConfigurationError{Message: "Qredex Config requires ClientID"}
	}
	if strings.TrimSpace(c.ClientSecret) == "" {
		return &ConfigurationError{Message: "Qredex Config requires ClientSecret"}
	}
	return nil
}

func (c *Config) resolvedBaseURL() string {
	if c.BaseURL != "" {
		return strings.TrimRight(c.BaseURL, "/")
	}
	if c.Environment == "" {
		return Production.baseURL()
	}
	return c.Environment.baseURL()
}

func (c *Config) resolvedTimeout() time.Duration {
	if c.Timeout > 0 {
		return c.Timeout
	}
	return 10 * time.Second
}

func (c *Config) resolvedRetryMax() int {
	if c.RetryMax <= 0 {
		return 0
	}
	if c.RetryMax > 5 {
		return 5
	}
	return c.RetryMax
}

func (c *Config) resolvedRetryBaseDelay() time.Duration {
	if c.RetryBaseDelay > 0 {
		return c.RetryBaseDelay
	}
	return 500 * time.Millisecond
}

func (c *Config) resolvedRetryMaxDelay() time.Duration {
	if c.RetryMaxDelay > 0 {
		return c.RetryMaxDelay
	}
	return 30 * time.Second
}

func (c *Config) scopeString() string {
	if len(c.Scopes) == 0 {
		return ""
	}
	parts := make([]string, len(c.Scopes))
	for i, s := range c.Scopes {
		parts[i] = string(s)
	}
	return strings.Join(parts, " ")
}

// ConfigBuilder provides a fluent builder for Config.
type ConfigBuilder struct {
	cfg Config
}

// NewConfigBuilder returns a new ConfigBuilder.
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{cfg: Config{}}
}

// WithClientID sets the ClientID.
func (b *ConfigBuilder) WithClientID(id string) *ConfigBuilder {
	b.cfg.ClientID = id
	return b
}

// WithClientSecret sets the ClientSecret.
func (b *ConfigBuilder) WithClientSecret(secret string) *ConfigBuilder {
	b.cfg.ClientSecret = secret
	return b
}

// WithScopes sets the OAuth scopes.
func (b *ConfigBuilder) WithScopes(scopes ...Scope) *ConfigBuilder {
	b.cfg.Scopes = scopes
	return b
}

// WithEnvironment sets the environment.
func (b *ConfigBuilder) WithEnvironment(env Environment) *ConfigBuilder {
	b.cfg.Environment = env
	return b
}

// WithBaseURL sets the base URL.
func (b *ConfigBuilder) WithBaseURL(url string) *ConfigBuilder {
	b.cfg.BaseURL = url
	return b
}

// WithHTTPClient sets a custom HTTP client.
func (b *ConfigBuilder) WithHTTPClient(client *http.Client) *ConfigBuilder {
	b.cfg.HTTPClient = client
	return b
}

// WithTimeout sets the request timeout.
func (b *ConfigBuilder) WithTimeout(timeout time.Duration) *ConfigBuilder {
	b.cfg.Timeout = timeout
	return b
}

// WithUserAgentSuffix sets the user agent suffix.
func (b *ConfigBuilder) WithUserAgentSuffix(suffix string) *ConfigBuilder {
	b.cfg.UserAgentSuffix = suffix
	return b
}

// WithRetry configures retry parameters.
func (b *ConfigBuilder) WithRetry(max int, baseDelay, maxDelay time.Duration) *ConfigBuilder {
	b.cfg.RetryMax = max
	b.cfg.RetryBaseDelay = baseDelay
	b.cfg.RetryMaxDelay = maxDelay
	return b
}

// Build returns the constructed Config.
func (b *ConfigBuilder) Build() Config {
	return b.cfg
}

// Bootstrap creates a Qredex instance from environment variables.
//
// Required environment variables:
//   - QREDEX_CLIENT_ID
//   - QREDEX_CLIENT_SECRET
//
// Optional environment variables:
//   - QREDEX_SCOPE        (space-separated scopes)
//   - QREDEX_ENVIRONMENT  (production | staging | development)
//   - QREDEX_BASE_URL     (overrides environment-resolved URL)
func Bootstrap() (*Qredex, error) {
	clientID := strings.TrimSpace(os.Getenv("QREDEX_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("QREDEX_CLIENT_SECRET"))
	rawScope := strings.TrimSpace(os.Getenv("QREDEX_SCOPE"))
	rawEnv := strings.TrimSpace(os.Getenv("QREDEX_ENVIRONMENT"))
	rawBaseURL := strings.TrimSpace(os.Getenv("QREDEX_BASE_URL"))

	if clientID == "" {
		return nil, &ConfigurationError{Message: "Bootstrap requires QREDEX_CLIENT_ID environment variable"}
	}
	if clientSecret == "" {
		return nil, &ConfigurationError{Message: "Bootstrap requires QREDEX_CLIENT_SECRET environment variable"}
	}

	env := Production
	if rawEnv != "" {
		switch rawEnv {
		case "production":
			env = Production
		case "staging":
			env = Staging
		case "development":
			env = Development
		default:
			return nil, &ConfigurationError{
				Message: fmt.Sprintf("QREDEX_ENVIRONMENT must be 'production', 'staging', or 'development'; got %q", rawEnv),
			}
		}
	}

	var scopes []Scope
	if rawScope != "" {
		for _, s := range strings.Fields(rawScope) {
			scopes = append(scopes, Scope(s))
		}
	}

	cfg := Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Environment:  env,
		BaseURL:      rawBaseURL,
		Scopes:       scopes,
	}

	return New(cfg)
}

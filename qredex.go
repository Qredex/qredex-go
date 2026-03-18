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
	"fmt"
	"net/http"
)

// Qredex is the main entrypoint for the Qredex Integrations API SDK.
type Qredex struct {
	creators *CreatorsResource
	links    *LinksResource
	intents  *IntentsResource
	orders   *OrdersResource
	refunds  *RefundsResource
	hc       *httpClient
	config   *Config
}

// New creates a Qredex client with explicit configuration.
// The configuration is validated before returning. If validation fails, an error is returned.
func New(cfg Config) (*Qredex, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	tp := newTokenProvider(&cfg, httpClient)
	hc := newHTTPClient(&cfg, httpClient, tp)

	return &Qredex{
		creators: newCreatorsResource(hc),
		links:    newLinksResource(hc),
		intents:  newIntentsResource(hc),
		orders:   newOrdersResource(hc),
		refunds:  newRefundsResource(hc),
		hc:       hc,
		config:   &cfg,
	}, nil
}

// Creators returns the Creators resource client.
func (q *Qredex) Creators() *CreatorsResource {
	return q.creators
}

// Links returns the Links resource client.
func (q *Qredex) Links() *LinksResource {
	return q.links
}

// Intents returns the Intents resource client.
func (q *Qredex) Intents() *IntentsResource {
	return q.intents
}

// Orders returns the Orders resource client.
func (q *Qredex) Orders() *OrdersResource {
	return q.orders
}

// Refunds returns the Refunds resource client.
func (q *Qredex) Refunds() *RefundsResource {
	return q.refunds
}

// CreatorsResource is the creators API resource group.
type CreatorsResource struct {
	hc *httpClient
}

func newCreatorsResource(hc *httpClient) *CreatorsResource {
	return &CreatorsResource{hc: hc}
}

// Create creates a new creator.
func (cr *CreatorsResource) Create(ctx context.Context, req CreateCreatorRequest) (*Creator, error) {
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("qredex: creators.create failed: %w", err)
	}
	var result Creator
	err := cr.hc.request(ctx, "POST", "/api/v1/integrations/creators", req, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: creators.create failed: %w", err)
	}
	return &result, nil
}

// Get retrieves a creator by ID.
func (cr *CreatorsResource) Get(ctx context.Context, creatorID string) (*Creator, error) {
	if err := validateIdentifier("creator_id", creatorID); err != nil {
		return nil, fmt.Errorf("qredex: creators.get failed: %w", err)
	}
	var result Creator
	err := cr.hc.request(ctx, "GET", "/api/v1/integrations/creators/"+creatorID, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: creators.get failed: %w", err)
	}
	return &result, nil
}

// List retrieves a paginated list of creators.
func (cr *CreatorsResource) List(ctx context.Context, req ListCreatorsRequest) (*CreatorPage, error) {
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("qredex: creators.list failed: %w", err)
	}
	var result CreatorPage
	err := cr.hc.request(ctx, "GET", "/api/v1/integrations/creators", req, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: creators.list failed: %w", err)
	}
	return &result, nil
}

// LinksResource is the links API resource group.
type LinksResource struct {
	hc *httpClient
}

func newLinksResource(hc *httpClient) *LinksResource {
	return &LinksResource{hc: hc}
}

// Create creates a new influence link.
func (lr *LinksResource) Create(ctx context.Context, req CreateLinkRequest) (*Link, error) {
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("qredex: links.create failed: %w", err)
	}
	var result Link
	err := lr.hc.request(ctx, "POST", "/api/v1/integrations/links", req, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: links.create failed: %w", err)
	}
	return &result, nil
}

// Get retrieves a link by ID.
func (lr *LinksResource) Get(ctx context.Context, linkID string) (*Link, error) {
	if err := validateIdentifier("link_id", linkID); err != nil {
		return nil, fmt.Errorf("qredex: links.get failed: %w", err)
	}
	var result Link
	err := lr.hc.request(ctx, "GET", "/api/v1/integrations/links/"+linkID, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: links.get failed: %w", err)
	}
	return &result, nil
}

// List retrieves a paginated list of links.
func (lr *LinksResource) List(ctx context.Context, req ListLinksRequest) (*LinkPage, error) {
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("qredex: links.list failed: %w", err)
	}
	var result LinkPage
	err := lr.hc.request(ctx, "GET", "/api/v1/integrations/links", req, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: links.list failed: %w", err)
	}
	return &result, nil
}

// GetStats retrieves statistics for a link.
func (lr *LinksResource) GetStats(ctx context.Context, linkID string) (*LinkStats, error) {
	if err := validateIdentifier("link_id", linkID); err != nil {
		return nil, fmt.Errorf("qredex: links.getStats failed: %w", err)
	}
	var result LinkStats
	err := lr.hc.request(ctx, "GET", "/api/v1/integrations/links/"+linkID+"/stats", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: links.getStats failed: %w", err)
	}
	return &result, nil
}

// IntentsResource is the intents (IIT/PIT) API resource group.
type IntentsResource struct {
	hc *httpClient
}

func newIntentsResource(hc *httpClient) *IntentsResource {
	return &IntentsResource{hc: hc}
}

// IssueInfluenceIntentToken issues a new Influence Intent Token (IIT).
func (ir *IntentsResource) IssueInfluenceIntentToken(ctx context.Context, req IssueInfluenceIntentTokenRequest) (*InfluenceIntent, error) {
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("qredex: intents.issueInfluenceIntentToken failed: %w", err)
	}
	var result InfluenceIntent
	err := ir.hc.request(ctx, "POST", "/api/v1/integrations/intents/token", req, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: intents.issueInfluenceIntentToken failed: %w", err)
	}
	return &result, nil
}

// LockPurchaseIntent locks a Purchase Intent Token (PIT).
func (ir *IntentsResource) LockPurchaseIntent(ctx context.Context, req LockPurchaseIntentRequest) (*PurchaseIntent, error) {
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("qredex: intents.lockPurchaseIntent failed: %w", err)
	}
	var result PurchaseIntent
	err := ir.hc.request(ctx, "POST", "/api/v1/integrations/intents/lock", req, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: intents.lockPurchaseIntent failed: %w", err)
	}
	return &result, nil
}

// GetPurchaseIntent retrieves a Purchase Intent by its PIT token.
func (ir *IntentsResource) GetPurchaseIntent(ctx context.Context, pit string) (*PurchaseIntent, error) {
	if err := validateIdentifier("purchase_intent_token", pit); err != nil {
		return nil, fmt.Errorf("qredex: intents.getPurchaseIntent failed: %w", err)
	}
	var result PurchaseIntent
	err := ir.hc.request(ctx, "GET", "/api/v1/integrations/intents/"+pit, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: intents.getPurchaseIntent failed: %w", err)
	}
	return &result, nil
}

// Deprecated: GetLatestUnlocked is a nondeterministic recovery helper.
// Prefer explicit PIT handling in canonical integrations.
func (ir *IntentsResource) GetLatestUnlocked(ctx context.Context, hours *int) (*PurchaseIntent, error) {
	if err := optionalPositiveInt("hours", hours); err != nil {
		return nil, fmt.Errorf("qredex: intents.getLatestUnlocked failed: %w", err)
	}
	var result PurchaseIntent
	query := struct {
		Hours *int `json:"hours,omitempty"`
	}{Hours: hours}
	err := ir.hc.request(ctx, "GET", "/api/v1/integrations/intents/latest-unlocked", query, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: intents.getLatestUnlocked failed: %w", err)
	}
	return &result, nil
}

// OrdersResource is the orders API resource group.
type OrdersResource struct {
	hc *httpClient
}

func newOrdersResource(hc *httpClient) *OrdersResource {
	return &OrdersResource{hc: hc}
}

// RecordPaidOrder records a paid order for attribution.
func (or *OrdersResource) RecordPaidOrder(ctx context.Context, req RecordPaidOrderRequest) (*OrderAttribution, error) {
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("qredex: orders.recordPaidOrder failed: %w", err)
	}
	var result OrderAttribution
	err := or.hc.request(ctx, "POST", "/api/v1/integrations/orders/paid", req, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: orders.recordPaidOrder failed: %w", err)
	}
	return &result, nil
}

// List retrieves a paginated list of order attributions.
func (or *OrdersResource) List(ctx context.Context, req ListOrdersRequest) (*OrderAttributionPage, error) {
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("qredex: orders.list failed: %w", err)
	}
	var result OrderAttributionPage
	err := or.hc.request(ctx, "GET", "/api/v1/integrations/orders", req, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: orders.list failed: %w", err)
	}
	return &result, nil
}

// GetDetails retrieves full details for an order attribution by ID.
func (or *OrdersResource) GetDetails(ctx context.Context, orderAttributionID string) (*OrderAttributionDetails, error) {
	if err := validateIdentifier("order_attribution_id", orderAttributionID); err != nil {
		return nil, fmt.Errorf("qredex: orders.getDetails failed: %w", err)
	}
	var result OrderAttributionDetails
	err := or.hc.request(ctx, "GET", "/api/v1/integrations/orders/"+orderAttributionID+"/details", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: orders.getDetails failed: %w", err)
	}
	return &result, nil
}

// RefundsResource is the refunds API resource group.
type RefundsResource struct {
	hc *httpClient
}

func newRefundsResource(hc *httpClient) *RefundsResource {
	return &RefundsResource{hc: hc}
}

// RecordRefund records an order refund.
func (rr *RefundsResource) RecordRefund(ctx context.Context, req RecordRefundRequest) (*OrderAttribution, error) {
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("qredex: refunds.recordRefund failed: %w", err)
	}
	var result OrderAttribution
	err := rr.hc.request(ctx, "POST", "/api/v1/integrations/orders/refund", req, &result)
	if err != nil {
		return nil, fmt.Errorf("qredex: refunds.recordRefund failed: %w", err)
	}
	return &result, nil
}

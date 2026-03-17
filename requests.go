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

import "time"

// CreateCreatorRequest is a request to create a new creator.
type CreateCreatorRequest struct {
	Handle      string            `json:"handle"`
	DisplayName *string           `json:"display_name,omitempty"`
	Email       *string           `json:"email,omitempty"`
	Socials     map[string]string `json:"socials,omitempty"`
}

// ListCreatorsRequest is a request to list creators with optional pagination and filtering.
type ListCreatorsRequest struct {
	Page   *int           `json:"page,omitempty"`
	Size   *int           `json:"size,omitempty"`
	Status *CreatorStatus `json:"status,omitempty"`
}

// CreateLinkRequest is a request to create a new influence link.
type CreateLinkRequest struct {
	StoreID               string      `json:"store_id"`
	CreatorID             string      `json:"creator_id"`
	LinkName              string      `json:"link_name"`
	DestinationPath       string      `json:"destination_path"`
	Note                  *string     `json:"note,omitempty"`
	AttributionWindowDays *int        `json:"attribution_window_days,omitempty"`
	LinkExpiryAt          *time.Time  `json:"link_expiry_at,omitempty"`
	DiscountCode          *string     `json:"discount_code,omitempty"`
	Status                *LinkStatus `json:"status,omitempty"`
}

// ListLinksRequest is a request to list links with optional pagination and filtering.
type ListLinksRequest struct {
	Page        *int        `json:"page,omitempty"`
	Size        *int        `json:"size,omitempty"`
	Status      *LinkStatus `json:"status,omitempty"`
	Destination *string     `json:"destination,omitempty"`
	Expired     *bool       `json:"expired,omitempty"`
}

// IssueInfluenceIntentTokenRequest is a request to issue an Influence Intent Token.
type IssueInfluenceIntentTokenRequest struct {
	LinkID           string     `json:"link_id"`
	IPHash           *string    `json:"ip_hash,omitempty"`
	UserAgentHash    *string    `json:"user_agent_hash,omitempty"`
	Referrer         *string    `json:"referrer,omitempty"`
	LandingPath      *string    `json:"landing_path,omitempty"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	IntegrityVersion *int       `json:"integrity_version,omitempty"`
}

// LockPurchaseIntentRequest is a request to lock a Purchase Intent Token.
type LockPurchaseIntentRequest struct {
	Token            string  `json:"token"`
	Source           *string `json:"source,omitempty"`
	IntegrityVersion *int    `json:"integrity_version,omitempty"`
}

// RecordPaidOrderRequest is a request to record a paid order.
type RecordPaidOrderRequest struct {
	StoreID             string     `json:"store_id"`
	ExternalOrderID     string     `json:"external_order_id"`
	OrderNumber         *string    `json:"order_number,omitempty"`
	PaidAt              *time.Time `json:"paid_at,omitempty"`
	Currency            string     `json:"currency"`
	SubtotalPrice       *float64   `json:"subtotal_price,omitempty"`
	DiscountTotal       *float64   `json:"discount_total,omitempty"`
	TotalPrice          *float64   `json:"total_price,omitempty"`
	CustomerEmailHash   *string    `json:"customer_email_hash,omitempty"`
	CheckoutToken       *string    `json:"checkout_token,omitempty"`
	PurchaseIntentToken *string    `json:"purchase_intent_token,omitempty"`
}

// ListOrdersRequest is a request to list order attributions with optional pagination.
type ListOrdersRequest struct {
	Page *int `json:"page,omitempty"`
	Size *int `json:"size,omitempty"`
}

// RecordRefundRequest is a request to record an order refund.
type RecordRefundRequest struct {
	StoreID          string     `json:"store_id"`
	ExternalOrderID  string     `json:"external_order_id"`
	ExternalRefundID string     `json:"external_refund_id"`
	RefundTotal      *float64   `json:"refund_total,omitempty"`
	RefundedAt       *time.Time `json:"refunded_at,omitempty"`
}

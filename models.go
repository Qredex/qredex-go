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

// CreatorStatus is the status of a creator account.
type CreatorStatus string

const (
	CreatorStatusActive   CreatorStatus = "ACTIVE"
	CreatorStatusDisabled CreatorStatus = "DISABLED"
)

// LinkStatus is the status of an influence link.
type LinkStatus string

const (
	LinkStatusActive   LinkStatus = "ACTIVE"
	LinkStatusDisabled LinkStatus = "DISABLED"
)

// OrderSource indicates where an order originated.
type OrderSource string

const (
	OrderSourceShopify   OrderSource = "SHOPIFY"
	OrderSourceDirectAPI OrderSource = "DIRECT_API"
)

// DuplicateConfidence indicates the likelihood that an order is a duplicate.
type DuplicateConfidence string

const (
	DuplicateConfidenceLow    DuplicateConfidence = "LOW"
	DuplicateConfidenceMedium DuplicateConfidence = "MEDIUM"
	DuplicateConfidenceHigh   DuplicateConfidence = "HIGH"
)

// IntegrityReason explains why a token integrity check failed.
type IntegrityReason string

const (
	IntegrityReasonMissing         IntegrityReason = "MISSING"
	IntegrityReasonTampered        IntegrityReason = "TAMPERED"
	IntegrityReasonExpired         IntegrityReason = "EXPIRED"
	IntegrityReasonMismatched      IntegrityReason = "MISMATCHED"
	IntegrityReasonReplaced        IntegrityReason = "REPLACED"
	IntegrityReasonLinkInactive    IntegrityReason = "LINK_INACTIVE"
	IntegrityReasonCreatorInactive IntegrityReason = "CREATOR_INACTIVE"
)

// TokenIntegrity is the result of a token integrity check.
type TokenIntegrity string

const (
	TokenIntegrityValid   TokenIntegrity = "VALID"
	TokenIntegrityInvalid TokenIntegrity = "INVALID"
)

// ResolutionStatus is the final attribution resolution state of an order.
type ResolutionStatus string

const (
	ResolutionStatusAttributed   ResolutionStatus = "ATTRIBUTED"
	ResolutionStatusUnattributed ResolutionStatus = "UNATTRIBUTED"
	ResolutionStatusRejected     ResolutionStatus = "REJECTED"
)

// OriginMatchStatus indicates whether the request origin matched stored intent origins.
type OriginMatchStatus string

const (
	OriginMatchStatusMatch    OriginMatchStatus = "MATCH"
	OriginMatchStatusMismatch OriginMatchStatus = "MISMATCH"
	OriginMatchStatusAbsent   OriginMatchStatus = "ABSENT"
	OriginMatchStatusUnknown  OriginMatchStatus = "UNKNOWN"
)

// IntegrityBand is a banded integrity score for order review.
type IntegrityBand string

const (
	IntegrityBandHigh     IntegrityBand = "HIGH"
	IntegrityBandMedium   IntegrityBand = "MEDIUM"
	IntegrityBandLow      IntegrityBand = "LOW"
	IntegrityBandCritical IntegrityBand = "CRITICAL"
)

// WindowStatus indicates whether an order falls within the attribution window.
type WindowStatus string

const (
	WindowStatusWithin  WindowStatus = "WITHIN"
	WindowStatusOutside WindowStatus = "OUTSIDE"
	WindowStatusUnknown WindowStatus = "UNKNOWN"
)

// Creator represents a creator account in the Qredex system.
type Creator struct {
	ID          string            `json:"id"`
	Handle      string            `json:"handle"`
	Status      CreatorStatus     `json:"status"`
	DisplayName *string           `json:"display_name"`
	Email       *string           `json:"email"`
	Socials     map[string]string `json:"socials"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// CreatorListItem includes creator fields plus link/order/revenue statistics.
type CreatorListItem struct {
	Creator
	LinksCount   int64   `json:"links_count"`
	OrdersCount  int64   `json:"orders_count"`
	RevenueTotal float64 `json:"revenue_total"`
}

// CreatorPage is a paginated list of creators.
type CreatorPage struct {
	Items         []CreatorListItem `json:"items"`
	Page          int               `json:"page"`
	Size          int               `json:"size"`
	TotalElements int64             `json:"total_elements"`
	TotalPages    int               `json:"total_pages"`
}

// Link represents an influence link.
type Link struct {
	ID                    string     `json:"id"`
	MerchantID            string     `json:"merchant_id"`
	StoreID               string     `json:"store_id"`
	CreatorID             string     `json:"creator_id"`
	LinkName              string     `json:"link_name"`
	LinkCode              string     `json:"link_code"`
	PublicLinkURL         string     `json:"public_link_url"`
	DestinationPath       string     `json:"destination_path"`
	Note                  *string    `json:"note"`
	Status                LinkStatus `json:"status"`
	AttributionWindowDays int        `json:"attribution_window_days"`
	LinkExpiryAt          *time.Time `json:"link_expiry_at"`
	DisabledAt            *time.Time `json:"disabled_at"`
	DiscountCode          *string    `json:"discount_code"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// LinkListItem includes link fields plus creator and statistics.
type LinkListItem struct {
	Link
	CreatorHandle      string  `json:"creator_handle"`
	CreatorDisplayName *string `json:"creator_display_name"`
	ClicksCount        int64   `json:"clicks_count"`
	OrdersCount        int64   `json:"orders_count"`
	RevenueTotal       float64 `json:"revenue_total"`
}

// LinkPage is a paginated list of links.
type LinkPage struct {
	Items         []LinkListItem `json:"items"`
	Page          int            `json:"page"`
	Size          int            `json:"size"`
	TotalElements int64          `json:"total_elements"`
	TotalPages    int            `json:"total_pages"`
}

// LinkStats contains click, order, and revenue statistics for a link.
type LinkStats struct {
	LinkID            string     `json:"link_id"`
	ClicksCount       int64      `json:"clicks_count"`
	SessionsCount     int64      `json:"sessions_count"`
	OrdersCount       int64      `json:"orders_count"`
	RevenueTotal      float64    `json:"revenue_total"`
	TokenInvalidCount int64      `json:"token_invalid_count"`
	TokenMissingCount int64      `json:"token_missing_count"`
	LastClickAt       *time.Time `json:"last_click_at"`
	LastOrderAt       *time.Time `json:"last_order_at"`
}

// InfluenceIntent is an Influence Intent Token (IIT).
type InfluenceIntent struct {
	ID               string    `json:"id"`
	MerchantID       string    `json:"merchant_id"`
	LinkID           string    `json:"link_id"`
	Token            string    `json:"token"`
	TokenID          string    `json:"token_id"`
	IssuedAt         time.Time `json:"issued_at"`
	ExpiresAt        time.Time `json:"expires_at"`
	Status           string    `json:"status"`
	IntegrityVersion int       `json:"integrity_version"`
	IPHash           *string   `json:"ip_hash"`
	UserAgentHash    *string   `json:"user_agent_hash"`
	Referrer         *string   `json:"referrer"`
	LandingPath      *string   `json:"landing_path"`
}

// PurchaseIntent is a Purchase Intent Token (PIT).
type PurchaseIntent struct {
	ID                            string             `json:"id"`
	MerchantID                    string             `json:"merchant_id"`
	StoreID                       string             `json:"store_id"`
	LinkID                        string             `json:"link_id"`
	InfluenceIntentID             *string            `json:"influence_intent_id"`
	Token                         string             `json:"token"`
	TokenID                       string             `json:"token_id"`
	Source                        *string            `json:"source"`
	OriginMatchStatus             *OriginMatchStatus `json:"origin_match_status"`
	WindowStatus                  *WindowStatus      `json:"window_status"`
	AttributionWindowDays         *int               `json:"attribution_window_days"`
	AttributionWindowDaysSnapshot *int               `json:"attribution_window_days_snapshot"`
	StoreDomainSnapshot           string             `json:"store_domain_snapshot"`
	LinkExpiryAtSnapshot          *time.Time         `json:"link_expiry_at_snapshot"`
	DiscountCodeSnapshot          *string            `json:"discount_code_snapshot"`
	IssuedAt                      time.Time          `json:"issued_at"`
	ExpiresAt                     time.Time          `json:"expires_at"`
	LockedAt                      *time.Time         `json:"locked_at"`
	IntegrityVersion              int                `json:"integrity_version"`
	Eligible                      *bool              `json:"eligible"`
	CreatedAt                     *time.Time         `json:"created_at"`
	UpdatedAt                     *time.Time         `json:"updated_at"`
}

// OrderAttribution represents a recorded order and its attribution to a creator/link.
type OrderAttribution struct {
	ID                            string               `json:"id"`
	MerchantID                    string               `json:"merchant_id"`
	OrderSource                   OrderSource          `json:"order_source"`
	ExternalOrderID               string               `json:"external_order_id"`
	OrderNumber                   *string              `json:"order_number"`
	PaidAt                        *time.Time           `json:"paid_at"`
	Currency                      string               `json:"currency"`
	SubtotalPrice                 *float64             `json:"subtotal_price"`
	DiscountTotal                 *float64             `json:"discount_total"`
	TotalPrice                    *float64             `json:"total_price"`
	PurchaseIntentToken           *string              `json:"purchase_intent_token"`
	LinkID                        *string              `json:"link_id"`
	LinkName                      *string              `json:"link_name"`
	LinkCode                      *string              `json:"link_code"`
	CreatorID                     *string              `json:"creator_id"`
	CreatorHandle                 *string              `json:"creator_handle"`
	CreatorDisplayName            *string              `json:"creator_display_name"`
	DuplicateSuspect              bool                 `json:"duplicate_suspect"`
	DuplicateConfidence           *DuplicateConfidence `json:"duplicate_confidence"`
	DuplicateReason               *string              `json:"duplicate_reason"`
	DuplicateOfOrderAttributionID *string              `json:"duplicate_of_order_attribution_id"`
	WindowStatus                  *WindowStatus        `json:"window_status"`
	TokenIntegrity                *TokenIntegrity      `json:"token_integrity"`
	IntegrityReason               *IntegrityReason     `json:"integrity_reason"`
	OriginMatchStatus             *OriginMatchStatus   `json:"origin_match_status"`
	IntegrityScore                int                  `json:"integrity_score"`
	IntegrityBand                 IntegrityBand        `json:"integrity_band"`
	ReviewRequired                bool                 `json:"review_required"`
	ResolutionStatus              ResolutionStatus     `json:"resolution_status"`
	CreatedAt                     time.Time            `json:"created_at"`
	UpdatedAt                     time.Time            `json:"updated_at"`
}

// OrderAttributionPage is a paginated list of order attributions.
type OrderAttributionPage struct {
	Items         []OrderAttribution `json:"items"`
	Page          int                `json:"page"`
	Size          int                `json:"size"`
	TotalElements int64              `json:"total_elements"`
	TotalPages    int                `json:"total_pages"`
}

// OrderAttributionDetails includes full score breakdown and timeline for an order.
type OrderAttributionDetails struct {
	OrderAttribution
	AttributionLockedAt   *time.Time `json:"attribution_locked_at"`
	AttributionWindowDays *int       `json:"attribution_window_days"`
}

// OAuthTokenResponse is the response from the OAuth 2.0 token endpoint.
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

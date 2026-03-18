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
	"regexp"
	"strings"
)

var currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)

func validateIdentifier(name, value string) error {
	return requireNonEmptyString(name, value)
}

func (r CreateCreatorRequest) validate() error {
	if err := requireNonEmptyString("handle", r.Handle); err != nil {
		return err
	}
	if err := optionalNonEmptyString("display_name", r.DisplayName); err != nil {
		return err
	}
	if err := optionalNonEmptyString("email", r.Email); err != nil {
		return err
	}
	return nil
}

func (r ListCreatorsRequest) validate() error {
	if err := optionalPositiveInt("page", r.Page); err != nil {
		return err
	}
	if err := optionalPositiveInt("size", r.Size); err != nil {
		return err
	}
	return nil
}

func (r CreateLinkRequest) validate() error {
	if err := requireNonEmptyString("store_id", r.StoreID); err != nil {
		return err
	}
	if err := requireNonEmptyString("creator_id", r.CreatorID); err != nil {
		return err
	}
	if err := requireNonEmptyString("link_name", r.LinkName); err != nil {
		return err
	}
	if err := requireNonEmptyString("destination_path", r.DestinationPath); err != nil {
		return err
	}
	if !strings.HasPrefix(r.DestinationPath, "/") {
		return &RequestValidationError{Message: "destination_path must start with '/'"}
	}
	if err := optionalNonEmptyString("note", r.Note); err != nil {
		return err
	}
	if err := optionalNonEmptyString("discount_code", r.DiscountCode); err != nil {
		return err
	}
	if err := optionalPositiveInt("attribution_window_days", r.AttributionWindowDays); err != nil {
		return err
	}
	return nil
}

func (r ListLinksRequest) validate() error {
	if err := optionalPositiveInt("page", r.Page); err != nil {
		return err
	}
	if err := optionalPositiveInt("size", r.Size); err != nil {
		return err
	}
	if err := optionalNonEmptyString("destination", r.Destination); err != nil {
		return err
	}
	return nil
}

func (r IssueInfluenceIntentTokenRequest) validate() error {
	if err := requireNonEmptyString("link_id", r.LinkID); err != nil {
		return err
	}
	if err := optionalNonEmptyString("ip_hash", r.IPHash); err != nil {
		return err
	}
	if err := optionalNonEmptyString("user_agent_hash", r.UserAgentHash); err != nil {
		return err
	}
	if err := optionalNonEmptyString("referrer", r.Referrer); err != nil {
		return err
	}
	if err := optionalNonEmptyString("landing_path", r.LandingPath); err != nil {
		return err
	}
	if err := optionalPositiveInt("integrity_version", r.IntegrityVersion); err != nil {
		return err
	}
	return nil
}

func (r LockPurchaseIntentRequest) validate() error {
	if err := requireNonEmptyString("token", r.Token); err != nil {
		return err
	}
	if err := optionalNonEmptyString("source", r.Source); err != nil {
		return err
	}
	if err := optionalPositiveInt("integrity_version", r.IntegrityVersion); err != nil {
		return err
	}
	return nil
}

func (r RecordPaidOrderRequest) validate() error {
	if err := requireNonEmptyString("store_id", r.StoreID); err != nil {
		return err
	}
	if err := requireNonEmptyString("external_order_id", r.ExternalOrderID); err != nil {
		return err
	}
	if err := requireCurrency("currency", r.Currency); err != nil {
		return err
	}
	if err := optionalNonEmptyString("order_number", r.OrderNumber); err != nil {
		return err
	}
	if err := optionalNonEmptyString("customer_email_hash", r.CustomerEmailHash); err != nil {
		return err
	}
	if err := optionalNonEmptyString("checkout_token", r.CheckoutToken); err != nil {
		return err
	}
	if err := optionalNonEmptyString("purchase_intent_token", r.PurchaseIntentToken); err != nil {
		return err
	}
	if err := optionalNonNegativeFloat("subtotal_price", r.SubtotalPrice); err != nil {
		return err
	}
	if err := optionalNonNegativeFloat("discount_total", r.DiscountTotal); err != nil {
		return err
	}
	if err := optionalNonNegativeFloat("total_price", r.TotalPrice); err != nil {
		return err
	}
	return nil
}

func (r ListOrdersRequest) validate() error {
	if err := optionalPositiveInt("page", r.Page); err != nil {
		return err
	}
	if err := optionalPositiveInt("size", r.Size); err != nil {
		return err
	}
	return nil
}

func (r RecordRefundRequest) validate() error {
	if err := requireNonEmptyString("store_id", r.StoreID); err != nil {
		return err
	}
	if err := requireNonEmptyString("external_order_id", r.ExternalOrderID); err != nil {
		return err
	}
	if err := requireNonEmptyString("external_refund_id", r.ExternalRefundID); err != nil {
		return err
	}
	if err := optionalNonNegativeFloat("refund_total", r.RefundTotal); err != nil {
		return err
	}
	return nil
}

func requireNonEmptyString(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return &RequestValidationError{Message: name + " must be a non-empty string"}
	}
	return nil
}

func optionalNonEmptyString(name string, value *string) error {
	if value == nil {
		return nil
	}
	if strings.TrimSpace(*value) == "" {
		return &RequestValidationError{Message: name + " must be a non-empty string when provided"}
	}
	return nil
}

func optionalPositiveInt(name string, value *int) error {
	if value == nil {
		return nil
	}
	if *value <= 0 {
		return &RequestValidationError{Message: name + " must be greater than zero when provided"}
	}
	return nil
}

func optionalNonNegativeFloat(name string, value *float64) error {
	if value == nil {
		return nil
	}
	if *value < 0 {
		return &RequestValidationError{Message: name + " must be non-negative when provided"}
	}
	return nil
}

func requireCurrency(name, value string) error {
	if err := requireNonEmptyString(name, value); err != nil {
		return err
	}
	if !currencyPattern.MatchString(value) {
		return &RequestValidationError{Message: name + " must be a three-letter uppercase ISO 4217 code"}
	}
	return nil
}

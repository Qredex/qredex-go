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

// Package main demonstrates locking a Purchase Intent Token (PIT) via the
// Qredex Integrations API.
//
// The PIT lock is typically called from the merchant backend at the checkout
// initiation step, converting the visitor's IIT cookie into a locked PIT that
// can then be submitted alongside the paid order for attribution.
//
// Required environment variables:
//
//	QREDEX_CLIENT_ID
//	QREDEX_CLIENT_SECRET
//	QREDEX_IIT_TOKEN   — the raw IIT token string from the visitor's cookie
package main

import (
	"context"
	"log"
	"os"

	"github.com/qredex/sdk-go"
)

func main() {
	q, err := qredex.Bootstrap()
	if err != nil {
		log.Fatalf("Failed to initialise Qredex SDK: %v", err)
	}

	iitToken := os.Getenv("QREDEX_IIT_TOKEN")
	if iitToken == "" {
		log.Fatal("QREDEX_IIT_TOKEN environment variable is required")
	}

	ctx := context.Background()

	pit, err := q.Intents().LockPurchaseIntent(ctx, qredex.LockPurchaseIntentRequest{
		Token:  iitToken,
		Source: strPtr("backend-checkout"),
	})
	if err != nil {
		log.Fatalf("Failed to lock PIT: %v", err)
	}

	log.Printf("Locked PIT: token_id=%s", pit.TokenID)
	if pit.Eligible != nil {
		log.Printf("  Eligible for attribution: %v", *pit.Eligible)
	}
	if pit.OriginMatchStatus != nil {
		log.Printf("  Origin match: %s", *pit.OriginMatchStatus)
	}
	if pit.WindowStatus != nil {
		log.Printf("  Window: %s", *pit.WindowStatus)
	}
	log.Printf("  PIT token (first 20 chars): %.20s...", pit.Token)
	// Store pit.Token alongside the order to submit for attribution.
}

func strPtr(s string) *string { return &s }

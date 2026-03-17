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

// Package main demonstrates recording a paid order for attribution via the
// Qredex Integrations API.
//
// This call should be made from your order-confirmation webhook or
// post-payment handler.  Pass the locked PIT token alongside the order
// to attribute revenue to the correct creator.
//
// Required environment variables:
//
//	QREDEX_CLIENT_ID
//	QREDEX_CLIENT_SECRET
//	QREDEX_STORE_ID
package main

import (
	"context"
	"log"
	"os"

	"github.com/Qredex/qredex-go"
)

func main() {
	q, err := qredex.Bootstrap()
	if err != nil {
		log.Fatalf("Failed to initialise Qredex SDK: %v", err)
	}

	storeID := os.Getenv("QREDEX_STORE_ID")
	if storeID == "" {
		log.Fatal("QREDEX_STORE_ID environment variable is required")
	}

	// In a real integration this comes from your order payload.
	pitToken := os.Getenv("QREDEX_PIT_TOKEN") // optional — attribution still recorded if absent

	ctx := context.Background()

	req := qredex.RecordPaidOrderRequest{
		StoreID:         storeID,
		ExternalOrderID: "order-100045",
		OrderNumber:     strPtr("100045"),
		Currency:        "USD",
		TotalPrice:      floatPtr(110.00),
		SubtotalPrice:   floatPtr(100.00),
		DiscountTotal:   floatPtr(10.00),
	}
	if pitToken != "" {
		req.PurchaseIntentToken = &pitToken
	}

	order, err := q.Orders().RecordPaidOrder(ctx, req)
	if err != nil {
		if qredex.IsConflictError(err) {
			// Order already recorded — safe to ignore in idempotent retry scenarios.
			log.Printf("Order already recorded (conflict): %v", err)
			return
		}
		log.Fatalf("Failed to record paid order: %v", err)
	}

	log.Printf("Recorded order: id=%s external_order_id=%s", order.ID, order.ExternalOrderID)
	log.Printf("  Resolution: %s", order.ResolutionStatus)
	if order.TokenIntegrity != nil {
		log.Printf("  Token integrity: %s", *order.TokenIntegrity)
	}
	if order.CreatorHandle != nil {
		log.Printf("  Attributed to: %s", *order.CreatorHandle)
	}
}

func strPtr(s string) *string     { return &s }
func floatPtr(f float64) *float64 { return &f }

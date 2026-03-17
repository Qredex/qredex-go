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

// Package main demonstrates recording an order refund via the Qredex
// Integrations API.
//
// Call this from your refund webhook or post-refund handler so that
// attributed revenue is adjusted correctly.
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

	"github.com/qredex/sdk-go"
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

	ctx := context.Background()

	updated, err := q.Refunds().RecordRefund(ctx, qredex.RecordRefundRequest{
		StoreID:          storeID,
		ExternalOrderID:  "order-100045",
		ExternalRefundID: "refund-100045-1",
		RefundTotal:      floatPtr(25.00),
	})
	if err != nil {
		log.Fatalf("Failed to record refund: %v", err)
	}

	log.Printf("Refund recorded: order_id=%s external_order_id=%s",
		updated.ID, updated.ExternalOrderID)
	log.Printf("  Resolution: %s", updated.ResolutionStatus)
}

func floatPtr(f float64) *float64 { return &f }

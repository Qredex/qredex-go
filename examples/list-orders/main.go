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

// Package main demonstrates listing order attributions via the Qredex
// Integrations API with simple pagination.
//
// Required environment variables:
//
//	QREDEX_CLIENT_ID
//	QREDEX_CLIENT_SECRET
package main

import (
	"context"
	"log"

	"github.com/qredex/sdk-go"
)

func main() {
	q, err := qredex.Bootstrap()
	if err != nil {
		log.Fatalf("Failed to initialise Qredex SDK: %v", err)
	}

	ctx := context.Background()

	page, err := q.Orders().List(ctx, qredex.ListOrdersRequest{
		Page: intPtr(1),
		Size: intPtr(20),
	})
	if err != nil {
		log.Fatalf("Failed to list orders: %v", err)
	}

	log.Printf("Orders: page=%d/%d total=%d",
		page.Page, page.TotalPages, page.TotalElements)

	for _, order := range page.Items {
		creator := "(unattributed)"
		if order.CreatorHandle != nil {
			creator = *order.CreatorHandle
		}
		log.Printf("  %s | %-12s | creator=%-20s | integrity=%d",
			order.ExternalOrderID,
			order.ResolutionStatus,
			creator,
			order.IntegrityScore,
		)
	}
}

func intPtr(i int) *int { return &i }

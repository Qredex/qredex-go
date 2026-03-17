// Copyright (C) 2026 — 2026, Qredex, LTD. All Rights Reserved.
//
// DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.
//
// Licensed under the Apache License, Version 2.0. See LICENSE for the full license text.
// You may not use this file except in compliance with that License.
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// limitations under the License.
//
// If you need additional information or have any questions, please email: copyright@qredex.com

// Package main demonstrates the complete canonical Qredex integration flow:
// Create Creator → Create Link → Issue IIT → Lock PIT → Record Paid Order → Record Refund
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/qredex/sdk-go"
)

func main() {
	// Parse flags
	dryRun := flag.Bool("dry-run", false, "Run without making API calls")
	flag.Parse()

	if *dryRun {
		log.Println("Running in dry-run mode - no API calls will be made")
	}

	// Initialize Qredex SDK from environment
	q, err := qredex.Bootstrap()
	if err != nil {
		log.Fatalf("Failed to initialize Qredex: %v", err)
	}

	ctx := context.Background()

	// Step 1: Create a creator
	log.Println("Step 1: Creating creator...")
	creator, err := q.Creators().Create(ctx, qredex.CreateCreatorRequest{
		Handle:      fmt.Sprintf("demo-%d", time.Now().Unix()),
		DisplayName: strPtr("Demo Creator"),
		Email:       strPtr("demo@example.com"),
	})
	if err != nil {
		log.Fatalf("Failed to create creator: %v", err)
	}
	log.Printf("✓ Created creator: %s (%s)", creator.Handle, creator.ID)

	// Step 2: Create an influence link
	log.Println("Step 2: Creating influence link...")
	storeID := os.Getenv("QREDEX_STORE_ID")
	if storeID == "" {
		log.Fatal("QREDEX_STORE_ID environment variable is required")
	}

	link, err := q.Links().Create(ctx, qredex.CreateLinkRequest{
		StoreID:               storeID,
		CreatorID:             creator.ID,
		LinkName:              "demo-spring-launch",
		DestinationPath:       "/products/spring",
		AttributionWindowDays: intPtr(30),
	})
	if err != nil {
		log.Fatalf("Failed to create link: %v", err)
	}
	log.Printf("✓ Created link: %s (%s)", link.LinkName, link.ID)
	log.Printf("  Public URL: %s", link.PublicLinkURL)

	// Step 3: Issue an Influence Intent Token (IIT)
	log.Println("Step 3: Issuing Influence Intent Token (IIT)...")
	iit, err := q.Intents().IssueInfluenceIntentToken(ctx, qredex.IssueInfluenceIntentTokenRequest{
		LinkID:      link.ID,
		LandingPath: strPtr("/products/spring"),
	})
	if err != nil {
		log.Fatalf("Failed to issue IIT: %v", err)
	}
	log.Printf("✓ Issued IIT: %s", iit.TokenID)
	log.Printf("  Token: %s...", iit.Token[:20])

	// Step 4: Lock a Purchase Intent Token (PIT)
	log.Println("Step 4: Locking Purchase Intent Token (PIT)...")
	pit, err := q.Intents().LockPurchaseIntent(ctx, qredex.LockPurchaseIntentRequest{
		Token:  iit.Token,
		Source: strPtr("demo-backend"),
	})
	if err != nil {
		log.Fatalf("Failed to lock PIT: %v", err)
	}
	log.Printf("✓ Locked PIT: %s", pit.TokenID)
	log.Printf("  Eligible: %v", *pit.Eligible)

	// Step 5: Record a paid order
	log.Println("Step 5: Recording paid order...")
	order, err := q.Orders().RecordPaidOrder(ctx, qredex.RecordPaidOrderRequest{
		StoreID:             storeID,
		ExternalOrderID:     fmt.Sprintf("demo-order-%d", time.Now().Unix()),
		OrderNumber:         strPtr("DEMO-001"),
		Currency:            "USD",
		TotalPrice:          floatPtr(99.99),
		PaidAt:              &[]time.Time{time.Now()}[0],
		PurchaseIntentToken: strPtr(pit.Token),
	})
	if err != nil {
		log.Fatalf("Failed to record paid order: %v", err)
	}
	log.Printf("✓ Recorded order: %s", order.ID)
	log.Printf("  Resolution Status: %s", order.ResolutionStatus)
	log.Printf("  Token Integrity: %v", order.TokenIntegrity)
	log.Printf("  Integrity Score: %d (%s)", order.IntegrityScore, order.IntegrityBand)
	if order.CreatorHandle != nil {
		log.Printf("  Attributed to: %s", *order.CreatorHandle)
	}

	// Step 6: Record a refund
	log.Println("Step 6: Recording refund...")
	refund, err := q.Refunds().RecordRefund(ctx, qredex.RecordRefundRequest{
		StoreID:          storeID,
		ExternalOrderID:  order.ExternalOrderID,
		ExternalRefundID: fmt.Sprintf("demo-refund-%d", time.Now().Unix()),
		RefundTotal:      floatPtr(25.00),
		RefundedAt:       &[]time.Time{time.Now()}[0],
	})
	if err != nil {
		log.Fatalf("Failed to record refund: %v", err)
	}
	log.Printf("✓ Recorded refund: %s", refund.ID)

	// Summary
	log.Println("\n=== Canonical Flow Complete ===")
	log.Printf("Creator: %s (%s)", creator.Handle, creator.ID)
	log.Printf("Link: %s (%s)", link.LinkName, link.ID)
	log.Printf("IIT: %s", iit.TokenID)
	log.Printf("PIT: %s", pit.TokenID)
	log.Printf("Order: %s (%s)", order.ID, order.ResolutionStatus)
	log.Printf("Refund: %s", refund.ID)
}

// Helper functions
func strPtr(s string) *string     { return &s }
func intPtr(i int) *int           { return &i }
func floatPtr(f float64) *float64 { return &f }

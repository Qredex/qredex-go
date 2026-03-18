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
	"log"
	"os"
)

// Example_bootstrap demonstrates creating a Qredex client from environment variables.
func Example_bootstrap() {
	// Requires:
	//   QREDEX_CLIENT_ID
	//   QREDEX_CLIENT_SECRET
	// Optional:
	//   QREDEX_SCOPE (space-separated OAuth scopes)
	//   QREDEX_ENVIRONMENT (production, staging, development)

	qredex, err := Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	// Use qredex.Creators(), qredex.Links(), qredex.Orders(), etc.
	_ = qredex
}

// Example_explicit shows creating a Qredex client with explicit configuration.
func Example_explicit() {
	qredex, err := New(Config{
		ClientID:     "my-client-id",
		ClientSecret: "my-client-secret",
		Environment:  Production,
		Scopes: []Scope{
			ScopeCreatorsWrite,
			ScopeLinksWrite,
			ScopeOrdersWrite,
			ScopeIntentsWrite,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	_ = qredex
}

// Example_createCreator demonstrates creating a creator.
func Example_createCreator() {
	qredex, err := Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	creator, err := qredex.Creators().Create(ctx, CreateCreatorRequest{
		Handle:      "alice",
		DisplayName: String("Alice"),
		Email:       String("alice@example.com"),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created creator: %s (%s)", creator.Handle, creator.ID)
}

// Example_createLink demonstrates creating an influence link.
func Example_createLink() {
	qredex, err := Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	storeID := os.Getenv("STORE_ID")     // Replace with your store ID
	creatorID := os.Getenv("CREATOR_ID") // Replace with your creator ID

	link, err := qredex.Links().Create(ctx, CreateLinkRequest{
		StoreID:               storeID,
		CreatorID:             creatorID,
		LinkName:              "spring-launch",
		DestinationPath:       "/collections/spring",
		AttributionWindowDays: Int(30),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created link: %s (%s)", link.LinkName, link.ID)
}

// Example_issueIIT demonstrates issuing an Influence Intent Token.
func Example_issueIIT() {
	qredex, err := Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	linkID := os.Getenv("LINK_ID") // Replace with your link ID

	iit, err := qredex.Intents().IssueInfluenceIntentToken(ctx, IssueInfluenceIntentTokenRequest{
		LinkID: linkID,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Issued IIT: %s", iit.TokenID)
}

// Example_lockPIT demonstrates locking a Purchase Intent Token.
func Example_lockPIT() {
	qredex, err := Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	iitToken := os.Getenv("IIT_TOKEN") // Replace with your IIT token

	pit, err := qredex.Intents().LockPurchaseIntent(ctx, LockPurchaseIntentRequest{
		Token:  iitToken,
		Source: String("backend-cart"),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Locked PIT: %s", pit.TokenID)
}

// Example_recordOrder demonstrates recording a paid order.
func Example_recordOrder() {
	qredex, err := Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	storeID := os.Getenv("STORE_ID") // Replace with your store ID

	order, err := qredex.Orders().RecordPaidOrder(ctx, RecordPaidOrderRequest{
		StoreID:             storeID,
		ExternalOrderID:     "order-12345",
		Currency:            "USD",
		TotalPrice:          Float64(99.99),
		PurchaseIntentToken: String("eyJhbGc..."), // Replace with actual PIT
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Recorded order: %s (status: %s)", order.ExternalOrderID, order.ResolutionStatus)
}

// Example_listOrders demonstrates listing order attributions.
func Example_listOrders() {
	qredex, err := Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	page, err := qredex.Orders().List(ctx, ListOrdersRequest{
		Page: Int(1),
		Size: Int(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d orders", len(page.Items))
	for _, order := range page.Items {
		log.Printf("  - %s: %v", order.ExternalOrderID, order.ResolutionStatus)
	}
}

// Example_recordRefund demonstrates recording a refund.
func Example_recordRefund() {
	qredex, err := Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	storeID := os.Getenv("STORE_ID") // Replace with your store ID

	updated, err := qredex.Refunds().RecordRefund(ctx, RecordRefundRequest{
		StoreID:          storeID,
		ExternalOrderID:  "order-12345",
		ExternalRefundID: "refund-12345-1",
		RefundTotal:      Float64(25.50),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Updated order: %s", updated.ID)
}

// Example_scopes demonstrates all available OAuth scope constants.
func Example_scopes() {
	_ = ScopeAPI
	_ = ScopeLinksRead
	_ = ScopeLinksWrite
	_ = ScopeCreatorsRead
	_ = ScopeCreatorsWrite
	_ = ScopeOrdersRead
	_ = ScopeOrdersWrite
	_ = ScopeIntentsRead
	_ = ScopeIntentsWrite
	// Output:
}

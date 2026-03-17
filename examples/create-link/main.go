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

// Package main demonstrates creating an influence link via the Qredex Integrations API.
//
// Required environment variables:
//
//	QREDEX_CLIENT_ID
//	QREDEX_CLIENT_SECRET
//	QREDEX_STORE_ID
//	QREDEX_CREATOR_ID
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
	creatorID := os.Getenv("QREDEX_CREATOR_ID")
	if creatorID == "" {
		log.Fatal("QREDEX_CREATOR_ID environment variable is required")
	}

	ctx := context.Background()

	link, err := q.Links().Create(ctx, qredex.CreateLinkRequest{
		StoreID:               storeID,
		CreatorID:             creatorID,
		LinkName:              "spring-launch",
		DestinationPath:       "/collections/spring",
		AttributionWindowDays: intPtr(30),
		DiscountCode:          strPtr("ALICE10"),
	})
	if err != nil {
		log.Fatalf("Failed to create link: %v", err)
	}

	log.Printf("Created link: id=%s name=%s", link.ID, link.LinkName)
	log.Printf("  Public URL: %s", link.PublicLinkURL)
	log.Printf("  Code: %s", link.LinkCode)
}

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }

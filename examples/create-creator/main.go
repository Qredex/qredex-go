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

// Package main demonstrates creating a creator via the Qredex Integrations API.
//
// Required environment variables:
//
//	QREDEX_CLIENT_ID
//	QREDEX_CLIENT_SECRET
package main

import (
	"context"
	"log"

	"github.com/Qredex/qredex-go"
)

func main() {
	q, err := qredex.Bootstrap()
	if err != nil {
		log.Fatalf("Failed to initialise Qredex SDK: %v", err)
	}

	ctx := context.Background()

	creator, err := q.Creators().Create(ctx, qredex.CreateCreatorRequest{
		Handle:      "alice",
		DisplayName: strPtr("Alice"),
		Email:       strPtr("alice@example.com"),
		Socials: map[string]string{
			"instagram": "@alice",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create creator: %v", err)
	}

	log.Printf("Created creator: id=%s handle=%s status=%s", creator.ID, creator.Handle, creator.Status)
}

func strPtr(s string) *string { return &s }

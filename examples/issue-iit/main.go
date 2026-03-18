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

// Package main demonstrates issuing an Influence Intent Token (IIT) via the
// Qredex Integrations API.
//
// The IIT is issued by the merchant backend when a visitor arrives via an
// influence link.  The token is returned to the browser/storefront and stored
// as a first-party cookie so it can later be used to lock a PIT at checkout.
//
// Required environment variables:
//
//	QREDEX_CLIENT_ID
//	QREDEX_CLIENT_SECRET
//	QREDEX_LINK_ID
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

	linkID := os.Getenv("QREDEX_LINK_ID")
	if linkID == "" {
		log.Fatal("QREDEX_LINK_ID environment variable is required")
	}

	ctx := context.Background()

	iit, err := q.Intents().IssueInfluenceIntentToken(ctx, qredex.IssueInfluenceIntentTokenRequest{
		LinkID:      linkID,
		LandingPath: qredex.String("/collections/spring"),
		// Optionally capture visitor signals for integrity scoring.
		// IPHash:        qredex.String(hashIP(r.RemoteAddr)),
		// UserAgentHash: qredex.String(hashUA(r.Header.Get("User-Agent"))),
		// Referrer:      qredex.String(r.Referer()),
	})
	if err != nil {
		log.Fatalf("Failed to issue IIT: %v", err)
	}

	log.Printf("Issued IIT: token_id=%s", iit.TokenID)
	log.Printf("  Expires at: %s", iit.ExpiresAt)
}

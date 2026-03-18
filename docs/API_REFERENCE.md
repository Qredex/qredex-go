<!--
     ▄▄▄▄
   ▄█▀▀███▄▄              █▄
   ██    ██ ▄             ██
   ██    ██ ████▄▄█▀█▄ ▄████ ▄█▀█▄▀██ ██▀
   ██  ▄ ██ ██   ██▄█▀ ██ ██ ██▄█▀  ███
    ▀█████▄▄█▀  ▄▀█▄▄▄▄█▀███▄▀█▄▄▄▄██ ██▄
         ▀█

   Copyright (C) 2026 — 2026, Qredex, LTD. All Rights Reserved.

   DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.

   Licensed under the Apache License, Version 2.0. See LICENSE for the full license text.
   You may not use this file except in compliance with that License.
   Unless required by applicable law or agreed to in writing, software distributed under the
   License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
   either express or implied. See the License for the specific language governing permissions
   and limitations under the License.

   If you need additional information or have any questions, please email: copyright@qredex.com
-->

# API Reference

## Construction

- `New(Config) (*Qredex, error)`
- `Bootstrap() (*Qredex, error)`

## Creators

- `Creators().Create(ctx, CreateCreatorRequest) (*Creator, error)`
- `Creators().Get(ctx, creatorID string) (*Creator, error)`
- `Creators().List(ctx, ListCreatorsRequest) (*CreatorPage, error)`

## Links

- `Links().Create(ctx, CreateLinkRequest) (*Link, error)`
- `Links().Get(ctx, linkID string) (*Link, error)`
- `Links().List(ctx, ListLinksRequest) (*LinkPage, error)`
- `Links().GetStats(ctx, linkID string) (*LinkStats, error)`

## Intents

- `Intents().IssueInfluenceIntentToken(ctx, IssueInfluenceIntentTokenRequest) (*InfluenceIntent, error)`
- `Intents().LockPurchaseIntent(ctx, LockPurchaseIntentRequest) (*PurchaseIntent, error)`
- `Intents().GetPurchaseIntent(ctx, pit string) (*PurchaseIntent, error)`
- `Intents().GetLatestUnlocked(ctx, hours *int) (*PurchaseIntent, error)`

`GetLatestUnlocked` is deprecated for normal integrations. Prefer explicit PIT handling.

## Orders

- `Orders().RecordPaidOrder(ctx, RecordPaidOrderRequest) (*OrderAttribution, error)`
- `Orders().List(ctx, ListOrdersRequest) (*OrderAttributionPage, error)`
- `Orders().GetDetails(ctx, orderAttributionID string) (*OrderAttributionDetails, error)`

## Refunds

- `Refunds().RecordRefund(ctx, RecordRefundRequest) (*OrderAttribution, error)`

## Errors

- `ConfigurationError`
- `RequestValidationError`
- `ResponseDecodingError`
- `NetworkError`
- `APIError`
- `AuthenticationError`
- `AuthorizationError`
- `ValidationError`
- `NotFoundError`
- `ConflictError`
- `RateLimitError`

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

import "time"

// String returns a pointer to v. Use it for optional request fields.
func String(v string) *string { return &v }

// Int returns a pointer to v. Use it for optional request fields.
func Int(v int) *int { return &v }

// Float64 returns a pointer to v. Use it for optional request fields.
func Float64(v float64) *float64 { return &v }

// Bool returns a pointer to v. Use it for optional request fields.
func Bool(v bool) *bool { return &v }

// Time returns a pointer to v. Use it for optional request fields.
func Time(v time.Time) *time.Time { return &v }

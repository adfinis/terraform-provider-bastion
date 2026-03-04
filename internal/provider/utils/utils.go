// Copyright Adfinis 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package utils

func ToPtr[T any](v T) *T {
	return &v
}

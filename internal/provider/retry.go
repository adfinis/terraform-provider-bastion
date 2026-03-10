// Copyright Adfinis 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"errors"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/adfinis/bastion-go"
)

const (
	retryMaxAttempts = 5
	retryBaseDelay   = 1 * time.Second
	retryMaxDelay    = 30 * time.Second
)

// isRetryableError reports whether err is a transient Bastion API internal error
// that warrants a retry. The Bastion returns such errors under concurrent load
// (e.g. removing many group members or guest accesses simultaneously).
func isRetryableError(err error) bool {
	var apiErr *bastion.APIResponse
	retryableErrors := [...]string{
		"INTERNAL",
		"ERR_CHMOD_FAILED",
	}
	if errors.As(err, &apiErr) {
		code := strings.ToUpper(apiErr.ErrorCode)
		for _, retryable := range retryableErrors {
			if strings.Contains(code, retryable) {
				return true
			}
		}
	}
	return false
}

// retryOnInternalError runs op, retrying with exponential backoff whenever the
// operation fails with a retryable Bastion internal error. At most
// retryMaxAttempts attempts are made in total.
func retryOnInternalError(op func() error) error {
	delay := retryBaseDelay
	var err error
	for attempt := 0; attempt < retryMaxAttempts; attempt++ {
		err = op()
		if err == nil || !isRetryableError(err) {
			return err
		}
		if attempt < retryMaxAttempts-1 {
			// Add ±25 % jitter
			jitter := time.Duration(rand.Int64N(int64(delay)/2)) - delay/4
			time.Sleep(delay + jitter)
			delay = min(delay*2, retryMaxDelay)
		}
	}
	return err
}

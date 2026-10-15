// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"errors"
	"testing"
	"time"
)

func TestKindImageLoadDeadlineScalesWithSize(t *testing.T) {
	const gib = int64(1) << 30
	cases := []struct {
		name string
		size int64
		want time.Duration
	}{
		{"small runtime image", 200 << 20, 3 * time.Minute},
		{"exactly one GiB", gib, 3 * time.Minute},
		{"derived ollama image", 2787075596, 9 * time.Minute},
		{"far past the cap", 50 * gib, 20 * time.Minute},
	}
	for _, tc := range cases {
		if got := kindImageLoadDeadline(tc.size); got != tc.want {
			t.Errorf("%s: kindImageLoadDeadline(%d) = %s, want %s", tc.name, tc.size, got, tc.want)
		}
	}
}

func swapImageSizeProbe(t *testing.T, probe func(string) (int64, error)) {
	t.Helper()
	orig := imageSizeProbe
	imageSizeProbe = probe
	t.Cleanup(func() { imageSizeProbe = orig })
}

func TestImageLoadDeadlineUsesProbedSize(t *testing.T) {
	swapImageSizeProbe(t, func(string) (int64, error) { return 2787075596, nil })
	if got := imageLoadDeadline("localhost/ollama:test"); got != 9*time.Minute {
		t.Fatalf("imageLoadDeadline = %s, want 9m for the 2.8 GB image", got)
	}
}

func TestImageLoadDeadlineFailedProbeUsesCap(t *testing.T) {
	swapImageSizeProbe(t, func(string) (int64, error) { return 0, errors.New("no such image") })
	if got := imageLoadDeadline("localhost/missing:test"); got != kindImageLoadCap {
		t.Fatalf("imageLoadDeadline = %s, want the %s cap when the probe fails", got, kindImageLoadCap)
	}
}

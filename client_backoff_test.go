package opencode

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestRetryBackoffDelay(t *testing.T) {
	tests := []struct {
		name    string
		attempt int
		wantMin time.Duration
		wantMax time.Duration
	}{
		{name: "attempt 0", attempt: 0, wantMin: 250 * time.Millisecond, wantMax: 500 * time.Millisecond},
		{name: "attempt 1", attempt: 1, wantMin: 500 * time.Millisecond, wantMax: 1 * time.Second},
		{name: "attempt 2", attempt: 2, wantMin: 1 * time.Second, wantMax: 2 * time.Second},
		{name: "attempt 3", attempt: 3, wantMin: 2 * time.Second, wantMax: 4 * time.Second},
		{name: "attempt 4 capped", attempt: 4, wantMin: 4 * time.Second, wantMax: 8 * time.Second},
		{name: "attempt 8 capped", attempt: 8, wantMin: 4 * time.Second, wantMax: 8 * time.Second},
	}

	originalRand := retryBackoffRandInt63n
	t.Cleanup(func() {
		retryBackoffRandInt63n = originalRand
	})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			retryBackoffRandInt63n = func(int64) int64 { return 0 }
			gotMin := retryBackoffDelay(tc.attempt)
			if gotMin != tc.wantMin {
				t.Fatalf("retryBackoffDelay(%d) with min jitter = %s, want %s", tc.attempt, gotMin, tc.wantMin)
			}

			retryBackoffRandInt63n = func(n int64) int64 { return n - 1 }
			gotMax := retryBackoffDelay(tc.attempt)
			if gotMax != tc.wantMax {
				t.Fatalf("retryBackoffDelay(%d) with max jitter = %s, want %s", tc.attempt, gotMax, tc.wantMax)
			}
		})
	}
}

func TestParseRetryAfterDelay(t *testing.T) {
	now := time.Date(2026, time.March, 2, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name   string
		header string
		want   time.Duration
		ok     bool
	}{
		{name: "seconds value", header: "5", want: 5 * time.Second, ok: true},
		{name: "http date value", header: "Mon, 02 Mar 2026 12:00:03 GMT", want: 3 * time.Second, ok: true},
		{name: "past date clamps to zero", header: "Mon, 02 Mar 2026 11:59:00 GMT", want: 0, ok: true},
		{name: "invalid value", header: "soon", want: 0, ok: false},
		{name: "negative value", header: "-1", want: 0, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseRetryAfterDelay(tt.header, now)
			if ok != tt.ok {
				t.Fatalf("parseRetryAfterDelay ok=%v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("parseRetryAfterDelay delay=%v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetryDelayWithServerGuidance(t *testing.T) {
	now := time.Date(2026, time.March, 2, 12, 0, 0, 0, time.UTC)
	ctx := context.Background()
	originalRand := retryBackoffRandInt63n
	retryBackoffRandInt63n = func(int64) int64 { return 0 }
	t.Cleanup(func() {
		retryBackoffRandInt63n = originalRand
	})

	t.Run("uses retry-after header value when valid", func(t *testing.T) {
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"3"}}}
		delay := retryDelayWithServerGuidance(0, resp, ctx, now)
		if delay != 3*time.Second {
			t.Fatalf("delay=%v, want %v", delay, 3*time.Second)
		}
	})

	t.Run("honors retry-after longer than max backoff", func(t *testing.T) {
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"999"}}}
		delay := retryDelayWithServerGuidance(0, resp, ctx, now)
		if delay != 999*time.Second {
			t.Fatalf("delay=%v, want %v", delay, 999*time.Second)
		}
	})

	t.Run("falls back to exponential backoff for invalid header", func(t *testing.T) {
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"invalid"}}}
		delay := retryDelayWithServerGuidance(2, resp, ctx, now)
		if delay != 1*time.Second {
			t.Fatalf("delay=%v, want %v", delay, 1*time.Second)
		}
	})

	t.Run("bounds delay by context deadline", func(t *testing.T) {
		deadlineCtx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
		t.Cleanup(cancel)
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"5"}}}
		delay := retryDelayWithServerGuidance(0, resp, deadlineCtx, time.Now())
		if delay <= 0 || delay > 1200*time.Millisecond {
			t.Fatalf("delay=%v, want within (0, %v]", delay, 1200*time.Millisecond)
		}
	})
}

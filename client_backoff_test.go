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
		want    time.Duration
	}{
		{name: "attempt 0", attempt: 0, want: 500 * time.Millisecond},
		{name: "attempt 1", attempt: 1, want: 1 * time.Second},
		{name: "attempt 2", attempt: 2, want: 2 * time.Second},
		{name: "attempt 3", attempt: 3, want: 4 * time.Second},
		{name: "attempt 4 capped", attempt: 4, want: 8 * time.Second},
		{name: "attempt 8 capped", attempt: 8, want: 8 * time.Second},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := retryBackoffDelay(tc.attempt); got != tc.want {
				t.Fatalf("retryBackoffDelay(%d) = %s, want %s", tc.attempt, got, tc.want)
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

	t.Run("uses retry-after header value when valid", func(t *testing.T) {
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"3"}}}
		delay := retryDelayWithServerGuidance(0, resp, ctx, now)
		if delay != 3*time.Second {
			t.Fatalf("delay=%v, want %v", delay, 3*time.Second)
		}
	})

	t.Run("caps retry-after to max backoff", func(t *testing.T) {
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"999"}}}
		delay := retryDelayWithServerGuidance(0, resp, ctx, now)
		if delay != maxBackoff {
			t.Fatalf("delay=%v, want %v", delay, maxBackoff)
		}
	})

	t.Run("falls back to exponential backoff for invalid header", func(t *testing.T) {
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"invalid"}}}
		delay := retryDelayWithServerGuidance(2, resp, ctx, now)
		if delay != 2*time.Second {
			t.Fatalf("delay=%v, want %v", delay, 2*time.Second)
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

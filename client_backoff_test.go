package opencode

import (
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

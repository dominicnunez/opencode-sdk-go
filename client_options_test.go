package opencode_test

import (
	"testing"
	"time"

	"github.com/dominicnunez/opencode-sdk-go"
)

func TestWithTimeout_BoundaryValues(t *testing.T) {
	tests := []struct {
		name    string
		d       time.Duration
		wantErr bool
	}{
		{"negative", -1 * time.Second, true},
		{"zero", 0, true},
		{"positive", 1 * time.Second, false},
		{"one_nanosecond", 1 * time.Nanosecond, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := opencode.NewClient(opencode.WithTimeout(tt.d))
			if (err != nil) != tt.wantErr {
				t.Errorf("WithTimeout(%v): err=%v, wantErr=%v", tt.d, err, tt.wantErr)
			}
		})
	}
}

func TestWithHTTPClient_Nil(t *testing.T) {
	_, err := opencode.NewClient(opencode.WithHTTPClient(nil))
	if err == nil {
		t.Fatal("WithHTTPClient(nil): expected error, got nil")
	}
}

func TestWithMaxRetries_BoundaryValues(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		wantErr bool
	}{
		{"negative", -1, true},
		{"zero", 0, false},
		{"max_allowed", 10, false},
		{"exceeds_cap", 11, true},
		{"positive", 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := opencode.NewClient(opencode.WithMaxRetries(tt.n))
			if (err != nil) != tt.wantErr {
				t.Errorf("WithMaxRetries(%d): err=%v, wantErr=%v", tt.n, err, tt.wantErr)
			}
		})
	}
}

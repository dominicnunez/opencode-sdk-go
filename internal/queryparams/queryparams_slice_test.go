package queryparams

import (
	"strings"
	"testing"
)

func TestMarshal_UnsupportedSliceElementType_ReturnsError(t *testing.T) {
	type params struct {
		IDs []int `query:"ids"`
	}

	_, err := Marshal(params{IDs: []int{1, 2, 3}})
	if err == nil {
		t.Fatal("expected error for []int slice, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported slice element type") {
		t.Errorf("expected error containing %q, got %q", "unsupported slice element type", err.Error())
	}
}

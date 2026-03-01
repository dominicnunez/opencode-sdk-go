package queryparams

import (
	"strings"
	"testing"
)

func TestMarshal_SimpleString(t *testing.T) {
	type params struct {
		Query string `query:"query,required"`
	}

	result, err := Marshal(params{Query: "hello world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := result.Get("query"); got != "hello world" {
		t.Errorf("expected query=%q, got %q", "hello world", got)
	}
}

func TestMarshal_OptionalString(t *testing.T) {
	type params struct {
		Directory *string `query:"directory,omitempty"`
	}

	// Test with value
	dir := "src/main"
	result, err := Marshal(params{Directory: &dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := result.Get("directory"); got != "src/main" {
		t.Errorf("expected directory=%q, got %q", "src/main", got)
	}

	// Test with nil (should be omitted)
	result, err = Marshal(params{Directory: nil})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := result.Get("directory"); got != "" {
		t.Errorf("expected directory to be omitted, got %q", got)
	}
}

func TestMarshal_EmptyString(t *testing.T) {
	type params struct {
		Query string `query:"query"`
	}

	result, err := Marshal(params{Query: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty strings should be omitted
	if got := result.Get("query"); got != "" {
		t.Errorf("expected query to be omitted for empty string, got %q", got)
	}
}

func TestMarshal_MultipleFields(t *testing.T) {
	type params struct {
		Query     string  `query:"query,required"`
		Directory *string `query:"directory,omitempty"`
		Pattern   string  `query:"pattern,required"`
	}

	dir := "src"
	result, err := Marshal(params{
		Query:     "test",
		Directory: &dir,
		Pattern:   "*.go",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := result.Get("query"); got != "test" {
		t.Errorf("expected query=%q, got %q", "test", got)
	}
	if got := result.Get("directory"); got != "src" {
		t.Errorf("expected directory=%q, got %q", "src", got)
	}
	if got := result.Get("pattern"); got != "*.go" {
		t.Errorf("expected pattern=%q, got %q", "*.go", got)
	}
}

func TestMarshal_NoQueryTags(t *testing.T) {
	type params struct {
		Query string // no query tag
	}

	result, err := Marshal(params{Query: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return empty params since no query tags
	if len(result) != 0 {
		t.Errorf("expected empty params, got %v", result)
	}
}

func TestMarshal_DashTag(t *testing.T) {
	type params struct {
		Ignored string `query:"-"`
		Query   string `query:"query"`
	}

	result, err := Marshal(params{Ignored: "ignore", Query: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := result.Get("ignored"); got != "" {
		t.Errorf("expected ignored field to be omitted, got %q", got)
	}
	if got := result.Get("query"); got != "test" {
		t.Errorf("expected query=%q, got %q", "test", got)
	}
}

func TestMarshal_NilInput(t *testing.T) {
	result, err := Marshal(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty params for nil input, got %v", result)
	}
}

func TestMarshal_PointerToStruct(t *testing.T) {
	type params struct {
		Query string `query:"query"`
	}

	p := &params{Query: "test"}
	result, err := Marshal(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := result.Get("query"); got != "test" {
		t.Errorf("expected query=%q, got %q", "test", got)
	}
}

func TestMarshal_IntFields(t *testing.T) {
	type params struct {
		Count  int   `query:"count"`
		Limit  int64 `query:"limit"`
		Offset *int  `query:"offset,omitempty"`
	}

	offset := 10
	result, err := Marshal(params{
		Count:  5,
		Limit:  100,
		Offset: &offset,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := result.Get("count"); got != "5" {
		t.Errorf("expected count=5, got %q", got)
	}
	if got := result.Get("limit"); got != "100" {
		t.Errorf("expected limit=100, got %q", got)
	}
	if got := result.Get("offset"); got != "10" {
		t.Errorf("expected offset=10, got %q", got)
	}

	// Test nil optional int
	result, err = Marshal(params{Count: 5, Limit: 100, Offset: nil})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result.Get("offset"); got != "" {
		t.Errorf("expected offset to be omitted for nil, got %q", got)
	}
}

func TestMarshal_BoolFields(t *testing.T) {
	type params struct {
		Enabled *bool `query:"enabled,omitempty"`
	}

	enabled := true
	result, err := Marshal(params{Enabled: &enabled})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := result.Get("enabled"); got != "true" {
		t.Errorf("expected enabled=true, got %q", got)
	}

	disabled := false
	result, err = Marshal(params{Enabled: &disabled})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := result.Get("enabled"); got != "false" {
		t.Errorf("expected enabled=false, got %q", got)
	}
}

func TestMarshal_StringSlice(t *testing.T) {
	type params struct {
		Tags []string `query:"tags"`
	}

	result, err := Marshal(params{Tags: []string{"go", "rust", "python"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	values := result["tags"]
	if len(values) != 3 {
		t.Fatalf("expected 3 tag values, got %d", len(values))
	}

	expected := []string{"go", "rust", "python"}
	for i, want := range expected {
		if values[i] != want {
			t.Errorf("expected tags[%d]=%q, got %q", i, want, values[i])
		}
	}
}

func TestMarshal_UnexportedFields(t *testing.T) {
	type params struct {
		Query      string `query:"query"`
		unexported string `query:"unexported"` // should be ignored
	}

	result, err := Marshal(params{Query: "test", unexported: "ignore"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := result.Get("unexported"); got != "" {
		t.Errorf("expected unexported field to be ignored, got %q", got)
	}
	if got := result.Get("query"); got != "test" {
		t.Errorf("expected query=%q, got %q", "test", got)
	}
}

func TestParseTag(t *testing.T) {
	tests := []struct {
		tag      string
		wantName string
		wantReq  bool
		wantOmit bool
		wantErr  bool
	}{
		{"query,required", "query", true, false, false},
		{"directory,omitempty", "directory", false, true, false},
		{"pattern", "pattern", false, false, false},
		{"-", "-", false, false, false},
		{"count,required,omitempty", "count", true, true, true},
		{"query,omiempty", "query", false, false, true},
		{"query,reqired", "query", false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			name, req, omit, err := parseTag(tt.tag)
			if name != tt.wantName {
				t.Errorf("name: got %q, want %q", name, tt.wantName)
			}
			if req != tt.wantReq {
				t.Errorf("required: got %v, want %v", req, tt.wantReq)
			}
			if omit != tt.wantOmit {
				t.Errorf("omitempty: got %v, want %v", omit, tt.wantOmit)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("error: got %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarshal_RealSDKParams(t *testing.T) {
	// Test with actual SDK param struct patterns
	type SessionListParams struct {
		Directory *string `query:"directory,omitempty"`
	}

	type FindTextParams struct {
		Query     string  `query:"query,required"`
		Directory *string `query:"directory,omitempty"`
	}

	// Test SessionListParams with directory
	dir := "/home/user/project"
	result, err := Marshal(SessionListParams{Directory: &dir})
	if err != nil {
		t.Fatalf("SessionListParams: unexpected error: %v", err)
	}
	if got := result.Get("directory"); got != dir {
		t.Errorf("SessionListParams: expected directory=%q, got %q", dir, got)
	}

	// Test SessionListParams without directory
	result, err = Marshal(SessionListParams{Directory: nil})
	if err != nil {
		t.Fatalf("SessionListParams nil: unexpected error: %v", err)
	}
	if got := result.Get("directory"); got != "" {
		t.Errorf("SessionListParams nil: expected directory to be omitted, got %q", got)
	}

	// Test FindTextParams
	result, err = Marshal(FindTextParams{Query: "func main", Directory: &dir})
	if err != nil {
		t.Fatalf("FindTextParams: unexpected error: %v", err)
	}
	if got := result.Get("query"); got != "func main" {
		t.Errorf("FindTextParams: expected query=%q, got %q", "func main", got)
	}
	if got := result.Get("directory"); got != dir {
		t.Errorf("FindTextParams: expected directory=%q, got %q", dir, got)
	}
}

func TestMarshal_RequiredValidation(t *testing.T) {
	type params struct {
		Query string `query:"query,required"`
	}

	_, err := Marshal(params{Query: ""})
	if err == nil {
		t.Fatal("expected error for empty required field")
	}

	_, err = Marshal(params{Query: "valid"})
	if err != nil {
		t.Fatalf("unexpected error for valid required field: %v", err)
	}
}

func TestMarshal_RequiredSliceValidation(t *testing.T) {
	type params struct {
		Tags []string `query:"tags,required"`
	}

	_, err := Marshal(params{Tags: []string{}})
	if err == nil {
		t.Fatal("expected error for empty required slice")
	}

	_, err = Marshal(params{Tags: nil})
	if err == nil {
		t.Fatal("expected error for nil required slice")
	}

	result, err := Marshal(params{Tags: []string{"go"}})
	if err != nil {
		t.Fatalf("unexpected error for valid required slice: %v", err)
	}
	if got := result["tags"]; len(got) != 1 || got[0] != "go" {
		t.Errorf("expected tags=[go], got %v", got)
	}
}

func TestMarshal_EmbeddedStruct(t *testing.T) {
	type pagination struct {
		Limit  int `query:"limit,omitempty"`
		Offset int `query:"offset,omitempty"`
	}
	type params struct {
		pagination
		Query string `query:"query,required"`
	}

	result, err := Marshal(params{
		pagination: pagination{Limit: 10, Offset: 20},
		Query:      "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := result.Get("query"); got != "test" {
		t.Errorf("expected query=%q, got %q", "test", got)
	}
	if got := result.Get("limit"); got != "10" {
		t.Errorf("expected limit=10, got %q", got)
	}
	if got := result.Get("offset"); got != "20" {
		t.Errorf("expected offset=20, got %q", got)
	}
}

func TestMarshal_EmbeddedPointerToStruct(t *testing.T) {
	type pagination struct {
		Limit  int `query:"limit,omitempty"`
		Offset int `query:"offset,omitempty"`
	}
	type params struct {
		*pagination
		Query string `query:"query,required"`
	}

	// Non-nil embedded pointer
	result, err := Marshal(params{
		pagination: &pagination{Limit: 10, Offset: 20},
		Query:      "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result.Get("query"); got != "test" {
		t.Errorf("expected query=%q, got %q", "test", got)
	}
	if got := result.Get("limit"); got != "10" {
		t.Errorf("expected limit=10, got %q", got)
	}
	if got := result.Get("offset"); got != "20" {
		t.Errorf("expected offset=20, got %q", got)
	}

	// Nil embedded pointer — fields should be silently skipped
	result, err = Marshal(params{
		pagination: nil,
		Query:      "test",
	})
	if err != nil {
		t.Fatalf("unexpected error with nil embedded: %v", err)
	}
	if got := result.Get("query"); got != "test" {
		t.Errorf("expected query=%q, got %q", "test", got)
	}
	if got := result.Get("limit"); got != "" {
		t.Errorf("expected limit to be omitted for nil embedded, got %q", got)
	}
}

func TestMarshal_RequiredIntValidation(t *testing.T) {
	type params struct {
		Count int `query:"count,required"`
	}

	// Zero value should error
	_, err := Marshal(params{Count: 0})
	if err == nil {
		t.Fatal("expected error for zero-value required int")
	}

	// Non-zero should succeed
	result, err := Marshal(params{Count: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result.Get("count"); got != "5" {
		t.Errorf("expected count=5, got %q", got)
	}

	// Pointer-to-zero should succeed (intentional zero)
	type ptrParams struct {
		Count *int `query:"count,required"`
	}
	zero := 0
	result, err = Marshal(ptrParams{Count: &zero})
	if err != nil {
		t.Fatalf("unexpected error for pointer-to-zero required int: %v", err)
	}
	if got := result.Get("count"); got != "0" {
		t.Errorf("expected count=0, got %q", got)
	}
}

func TestMarshal_RequiredBoolValidation(t *testing.T) {
	type params struct {
		Active bool `query:"active,required"`
	}

	// False (zero value) should error
	_, err := Marshal(params{Active: false})
	if err == nil {
		t.Fatal("expected error for false required bool")
	}

	// True should succeed
	result, err := Marshal(params{Active: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result.Get("active"); got != "true" {
		t.Errorf("expected active=true, got %q", got)
	}

	// Pointer-to-false should succeed (intentional false)
	type ptrParams struct {
		Active *bool `query:"active,required"`
	}
	f := false
	result, err = Marshal(ptrParams{Active: &f})
	if err != nil {
		t.Fatalf("unexpected error for pointer-to-false required bool: %v", err)
	}
	if got := result.Get("active"); got != "false" {
		t.Errorf("expected active=false, got %q", got)
	}
}

func TestMarshal_RequiredOmitemptyConflict(t *testing.T) {
	type params struct {
		Count int `query:"count,required,omitempty"`
	}

	_, err := Marshal(params{Count: 5})
	if err == nil {
		t.Fatal("expected error for contradictory required,omitempty tag")
	}
}

func TestMarshal_UnrecognizedTagOption(t *testing.T) {
	type params struct {
		Query string `query:"query,omiempty"`
	}

	_, err := Marshal(params{Query: "test"})
	if err == nil {
		t.Fatal("expected error for unrecognized tag option")
	}
	if !strings.Contains(err.Error(), "unrecognized option") {
		t.Errorf("expected error to mention unrecognized option, got: %v", err)
	}
}

func TestMarshal_URLEncoding(t *testing.T) {
	type params struct {
		Query string `query:"query"`
	}

	result, err := Marshal(params{Query: "hello world & special=chars"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	encoded := result.Encode()
	expected := "query=hello+world+%26+special%3Dchars"
	if encoded != expected {
		t.Errorf("expected encoded query=%q, got %q", expected, encoded)
	}
}

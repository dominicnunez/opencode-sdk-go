package opencode

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestClientDo_NilResultResponseExceedsSizeLimit(t *testing.T) {
	const maxBodySize = int64(32)
	oversizedBody := strings.Repeat("x", int(maxBodySize)+1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(oversizedBody))
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL),
		WithMaxSuccessBodySize(maxBodySize),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	err = client.do(context.Background(), http.MethodGet, "/sessions", nil, nil)
	if err == nil {
		t.Fatal("expected size limit error for oversized response body with nil result")
	}
	if !strings.Contains(err.Error(), "response body exceeds") {
		t.Fatalf("expected body limit error, got: %v", err)
	}
}

func TestShouldMarshalRequestBody_CachesStructMetadata(t *testing.T) {
	type queryOnlyParams struct {
		Directory *string `query:"directory" json:"-"`
	}

	requestBodyFieldCache = sync.Map{}
	defer func() { requestBodyFieldCache = sync.Map{} }()

	paramsType := reflect.TypeOf(queryOnlyParams{})
	if _, exists := requestBodyFieldCache.Load(paramsType); exists {
		t.Fatal("expected empty request body field cache at start of test")
	}

	if shouldMarshalRequestBody(queryOnlyParams{}) {
		t.Fatal("expected query-only params to skip JSON body marshaling")
	}

	cachedValue, exists := requestBodyFieldCache.Load(paramsType)
	if !exists {
		t.Fatal("expected request body field decision to be cached")
	}
	if cachedValue.(bool) {
		t.Fatal("expected cached decision to indicate no JSON body fields")
	}
}

func TestShouldMarshalRequestBody_CacheIncludesEmbeddedJSONFields(t *testing.T) {
	type embeddedBody struct {
		Title string `json:"title"`
	}
	type paramsWithEmbeddedBody struct {
		embeddedBody
	}

	requestBodyFieldCache = sync.Map{}
	defer func() { requestBodyFieldCache = sync.Map{} }()

	paramsType := reflect.TypeOf(paramsWithEmbeddedBody{})
	if !shouldMarshalRequestBody(paramsWithEmbeddedBody{}) {
		t.Fatal("expected embedded json fields to trigger JSON body marshaling")
	}

	cachedValue, exists := requestBodyFieldCache.Load(paramsType)
	if !exists {
		t.Fatal("expected embedded body decision to be cached")
	}
	if !cachedValue.(bool) {
		t.Fatal("expected cached decision to include JSON body fields")
	}
}

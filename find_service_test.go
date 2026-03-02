package opencode

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindSymbols_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/find/symbol" {
			t.Errorf("expected path /find/symbol, got %s", r.URL.Path)
		}
		if q := r.URL.Query().Get("query"); q != "main" {
			t.Errorf("expected query param 'main', got %s", q)
		}

		resp := []Symbol{
			{
				Name: "main",
				Kind: SymbolKindFunction,
				Location: SymbolLocation{
					Uri: "file:///src/main.go",
					Range: SymbolLocationRange{
						Start: SymbolPosition{Line: 10, Character: 0},
						End:   SymbolPosition{Line: 10, Character: 4},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	symbols, err := client.Find.Symbols(context.Background(), &FindSymbolsParams{
		Query: "main",
	})
	if err != nil {
		t.Fatalf("Symbols failed: %v", err)
	}
	if len(symbols) != 1 {
		t.Fatalf("expected 1 symbol, got %d", len(symbols))
	}
	if symbols[0].Name != "main" {
		t.Errorf("expected name 'main', got %s", symbols[0].Name)
	}
	if symbols[0].Kind != SymbolKindFunction {
		t.Errorf("expected kind Function (%d), got %d", SymbolKindFunction, symbols[0].Kind)
	}
	if symbols[0].Location.Uri != "file:///src/main.go" {
		t.Errorf("expected URI file:///src/main.go, got %s", symbols[0].Location.Uri)
	}
}

func TestFindSymbols_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Find.Symbols(context.Background(), &FindSymbolsParams{Query: "main"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, apiErr.StatusCode)
	}
}

func TestFindText_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/find" {
			t.Errorf("expected path /find, got %s", r.URL.Path)
		}
		if q := r.URL.Query().Get("pattern"); q != "TODO" {
			t.Errorf("expected pattern query param 'TODO', got %s", q)
		}

		resp := []FindTextResponse{
			{
				LineNumber:     42,
				AbsoluteOffset: 1024,
				Lines:          FindTextResponseLines{Text: "// TODO: fix this"},
				Path:           FindTextResponsePath{Text: "src/main.go"},
				Submatches: []FindTextResponseSubmatch{
					{
						Start: 3,
						End:   7,
						Match: FindTextResponseSubmatchesMatch{Text: "TODO"},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	results, err := client.Find.Text(context.Background(), &FindTextParams{
		Pattern: "TODO",
	})
	if err != nil {
		t.Fatalf("Text failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].LineNumber != 42 {
		t.Errorf("expected line number 42, got %d", results[0].LineNumber)
	}
	if results[0].Path.Text != "src/main.go" {
		t.Errorf("expected path src/main.go, got %s", results[0].Path.Text)
	}
	if len(results[0].Submatches) != 1 {
		t.Fatalf("expected 1 submatch, got %d", len(results[0].Submatches))
	}
	if results[0].Submatches[0].Match.Text != "TODO" {
		t.Errorf("expected match text 'TODO', got %s", results[0].Submatches[0].Match.Text)
	}
}

func TestFindText_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Find.Text(context.Background(), &FindTextParams{Pattern: "TODO"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, apiErr.StatusCode)
	}
}

func TestFindFiles_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/find/file" {
			t.Errorf("expected path /find/file, got %s", r.URL.Path)
		}
		if q := r.URL.Query().Get("query"); q != "main" {
			t.Errorf("expected query param 'main', got %s", q)
		}

		resp := []string{
			"src/main.go",
			"src/main_test.go",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	files, err := client.Find.Files(context.Background(), &FindFilesParams{Query: "main"})
	if err != nil {
		t.Fatalf("Files failed: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if files[0] != "src/main.go" {
		t.Fatalf("expected first file src/main.go, got %s", files[0])
	}
	if files[1] != "src/main_test.go" {
		t.Fatalf("expected second file src/main_test.go, got %s", files[1])
	}
}

func TestFindFiles_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":"bad gateway"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Find.Files(context.Background(), &FindFilesParams{Query: "main"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status %d, got %d", http.StatusBadGateway, apiErr.StatusCode)
	}
}

func TestFindService_WhitespaceOnlyRequiredFieldsRejected(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Find.Files(context.Background(), &FindFilesParams{Query: "   "})
	if !errors.Is(err, &MissingRequiredParameterError{Parameter: "query"}) {
		t.Fatalf("expected missing required query error, got %v", err)
	}

	_, err = client.Find.Symbols(context.Background(), &FindSymbolsParams{Query: "\t\n"})
	if !errors.Is(err, &MissingRequiredParameterError{Parameter: "query"}) {
		t.Fatalf("expected missing required query error, got %v", err)
	}

	_, err = client.Find.Text(context.Background(), &FindTextParams{Pattern: "   "})
	if !errors.Is(err, &MissingRequiredParameterError{Parameter: "pattern"}) {
		t.Fatalf("expected missing required pattern error, got %v", err)
	}
}

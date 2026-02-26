package opencode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionDiff_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/session/ses123/diff" {
			t.Errorf("Expected path /session/ses123/diff, got %s", r.URL.Path)
		}

		diffs := []FileDiff{
			{
				File:      "main.go",
				Before:    "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}",
				After:     "package main\n\nfunc main() {\n\tprintln(\"hello world\")\n}",
				Additions: 1,
				Deletions: 1,
			},
			{
				File:      "README.md",
				Before:    "# Project",
				After:     "# Project\n\nDescription here",
				Additions: 2,
				Deletions: 0,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(diffs)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL))
	ctx := context.Background()

	result, err := client.Session.Diff(ctx, "ses123", &SessionDiffParams{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 diffs, got %d", len(result))
	}

	if result[0].File != "main.go" {
		t.Errorf("Expected file main.go, got %s", result[0].File)
	}
	if result[0].Additions != 1 {
		t.Errorf("Expected 1 addition, got %d", result[0].Additions)
	}
	if result[0].Deletions != 1 {
		t.Errorf("Expected 1 deletion, got %d", result[0].Deletions)
	}

	if result[1].File != "README.md" {
		t.Errorf("Expected file README.md, got %s", result[1].File)
	}
	if result[1].Additions != 2 {
		t.Errorf("Expected 2 additions, got %d", result[1].Additions)
	}
	if result[1].Deletions != 0 {
		t.Errorf("Expected 0 deletions, got %d", result[1].Deletions)
	}
}

func TestSessionDiff_WithMessageID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		messageID := r.URL.Query().Get("messageID")
		if messageID != "msg456" {
			t.Errorf("Expected messageID=msg456 in query params, got %s", messageID)
		}

		diffs := []FileDiff{
			{
				File:      "test.go",
				Before:    "",
				After:     "package test",
				Additions: 1,
				Deletions: 0,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(diffs)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL))
	ctx := context.Background()

	msgID := "msg456"
	result, err := client.Session.Diff(ctx, "ses123", &SessionDiffParams{
		MessageID: &msgID,
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 diff, got %d", len(result))
	}

	if result[0].File != "test.go" {
		t.Errorf("Expected file test.go, got %s", result[0].File)
	}
}

func TestSessionDiff_WithDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		directory := r.URL.Query().Get("directory")
		if directory != "/home/user/project" {
			t.Errorf("Expected directory=/home/user/project in query params, got %s", directory)
		}

		diffs := []FileDiff{}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(diffs)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL))
	ctx := context.Background()

	dir := "/home/user/project"
	result, err := client.Session.Diff(ctx, "ses123", &SessionDiffParams{
		Directory: &dir,
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("Expected 0 diffs, got %d", len(result))
	}
}

func TestSessionDiff_EmptyDiffs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FileDiff{})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL))
	ctx := context.Background()

	result, err := client.Session.Diff(ctx, "ses123", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("Expected 0 diffs, got %d", len(result))
	}
}

func TestSessionDiff_MissingID(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	_, err := client.Session.Diff(ctx, "", &SessionDiffParams{})
	if err == nil {
		t.Fatal("Expected error for missing id, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("Expected 'missing required id parameter' error, got: %v", err)
	}
}

func TestSessionDiff_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL))
	ctx := context.Background()

	_, err := client.Session.Diff(ctx, "ses123", &SessionDiffParams{})
	if err == nil {
		t.Fatal("Expected error for server error, got nil")
	}
}

func TestSessionDiff_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL))
	ctx := context.Background()

	_, err := client.Session.Diff(ctx, "ses123", &SessionDiffParams{})
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestFileDiff_Unmarshal(t *testing.T) {
	jsonData := `{
		"file": "example.go",
		"before": "old content",
		"after": "new content",
		"additions": 10,
		"deletions": 5
	}`

	var diff FileDiff
	err := json.Unmarshal([]byte(jsonData), &diff)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if diff.File != "example.go" {
		t.Errorf("Expected file example.go, got %s", diff.File)
	}
	if diff.Before != "old content" {
		t.Errorf("Expected before 'old content', got %s", diff.Before)
	}
	if diff.After != "new content" {
		t.Errorf("Expected after 'new content', got %s", diff.After)
	}
	if diff.Additions != 10 {
		t.Errorf("Expected 10 additions, got %d", diff.Additions)
	}
	if diff.Deletions != 5 {
		t.Errorf("Expected 5 deletions, got %d", diff.Deletions)
	}
}

func TestFileDiff_Marshal(t *testing.T) {
	diff := FileDiff{
		File:      "test.go",
		Before:    "before",
		After:     "after",
		Additions: 3,
		Deletions: 2,
	}

	data, err := json.Marshal(diff)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	var result FileDiff
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unexpected error on unmarshal: %v", err)
	}

	if result.File != diff.File {
		t.Errorf("Expected file %s, got %s", diff.File, result.File)
	}
	if result.Additions != diff.Additions {
		t.Errorf("Expected %d additions, got %d", diff.Additions, result.Additions)
	}
}

func TestSessionDiffParams_URLQuery(t *testing.T) {
	msgID := "msg123"
	dir := "/home/user"
	params := SessionDiffParams{
		MessageID: &msgID,
		Directory: &dir,
	}

	values, err := params.URLQuery()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if values.Get("messageID") != "msg123" {
		t.Errorf("Expected messageID=msg123, got %s", values.Get("messageID"))
	}
	if values.Get("directory") != "/home/user" {
		t.Errorf("Expected directory=/home/user, got %s", values.Get("directory"))
	}
}

func TestSessionDiffParams_URLQuery_Empty(t *testing.T) {
	params := SessionDiffParams{}

	values, err := params.URLQuery()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(values) != 0 {
		t.Errorf("Expected empty query params, got %v", values)
	}
}

package opencode

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFileRead_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/file/content" {
			t.Errorf("expected path /file/content, got %s", r.URL.Path)
		}
		if q := r.URL.Query().Get("path"); q != "/src/main.go" {
			t.Errorf("expected path query param /src/main.go, got %s", q)
		}

		resp := FileReadResponse{
			Content: "package main\n",
			Type:    FileReadResponseTypeText,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.File.Read(context.Background(), &FileReadParams{
		Path: "/src/main.go",
	})
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if result.Content != "package main\n" {
		t.Errorf("expected content %q, got %q", "package main\n", result.Content)
	}
	if result.Type != FileReadResponseTypeText {
		t.Errorf("expected type %q, got %q", FileReadResponseTypeText, result.Type)
	}
}

func TestFileRead_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.File.Read(context.Background(), &FileReadParams{Path: "/missing.go"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, apiErr.StatusCode)
	}
}

func TestFileStatus_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/file/status" {
			t.Errorf("expected path /file/status, got %s", r.URL.Path)
		}

		resp := []File{
			{
				Path:    "src/main.go",
				Status:  FileStatusModified,
				Added:   5,
				Removed: 2,
			},
			{
				Path:   "src/new.go",
				Status: FileStatusAdded,
				Added:  20,
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

	files, err := client.File.Status(context.Background(), nil)
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if files[0].Path != "src/main.go" {
		t.Errorf("expected path src/main.go, got %s", files[0].Path)
	}
	if files[0].Status != FileStatusModified {
		t.Errorf("expected status modified, got %s", files[0].Status)
	}
	if files[1].Status != FileStatusAdded {
		t.Errorf("expected status added, got %s", files[1].Status)
	}
}

func TestFileStatus_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.File.Status(context.Background(), nil)
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

func TestFileList_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/file" {
			t.Errorf("expected path /file, got %s", r.URL.Path)
		}
		if q := r.URL.Query().Get("path"); q != "/src" {
			t.Errorf("expected path query param /src, got %s", q)
		}

		resp := []FileNode{
			{
				Name: "main.go",
				Type: FileNodeTypeFile,
				Path: "/src/main.go",
			},
			{
				Name: "pkg",
				Type: FileNodeTypeDirectory,
				Path: "/src/pkg",
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

	nodes, err := client.File.List(context.Background(), &FileListParams{Path: "/src"})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	if nodes[0].Path != "/src/main.go" || nodes[0].Type != FileNodeTypeFile {
		t.Fatalf("unexpected first node: %+v", nodes[0])
	}
	if nodes[1].Path != "/src/pkg" || nodes[1].Type != FileNodeTypeDirectory {
		t.Fatalf("unexpected second node: %+v", nodes[1])
	}
}

func TestFileList_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"forbidden"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.File.List(context.Background(), &FileListParams{Path: "/src"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, apiErr.StatusCode)
	}
}

func TestFileService_WhitespaceOnlyPathRejected(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.File.List(context.Background(), &FileListParams{Path: "   "})
	if !errors.Is(err, &MissingRequiredParameterError{Parameter: "path"}) {
		t.Fatalf("expected missing required path error for List, got %v", err)
	}

	_, err = client.File.Read(context.Background(), &FileReadParams{Path: "\t"})
	if !errors.Is(err, &MissingRequiredParameterError{Parameter: "path"}) {
		t.Fatalf("expected missing required path error for Read, got %v", err)
	}
}

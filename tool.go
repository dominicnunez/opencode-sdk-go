package opencode

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
)

// ToolService handles experimental tool endpoints
type ToolService struct {
	client *Client
}

// ToolIDs represents a list of tool IDs
type ToolIDs []string

// ToolList represents a list of tools with their schemas
type ToolList []ToolListItem

// ToolListItem represents a single tool with its ID, description, and parameter schema
type ToolListItem struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

// ToolIDsParams are the parameters for GET /experimental/tool/ids
type ToolIDsParams struct {
	Directory *string `json:"-" query:"directory,omitempty"`
}

// URLQuery returns the query parameters for ToolIDsParams
func (p ToolIDsParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(p)
}

// ToolListParams are the parameters for GET /experimental/tool
type ToolListParams struct {
	Directory *string `json:"-" query:"directory,omitempty"`
	Provider  string  `query:"provider,required"`
	Model     string  `query:"model,required"`
}

// URLQuery returns the query parameters for ToolListParams
func (p ToolListParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(p)
}

// IDs retrieves all tool IDs (including built-in and dynamically registered)
// GET /experimental/tool/ids
func (s *ToolService) IDs(ctx context.Context, params *ToolIDsParams) (*ToolIDs, error) {
	if params == nil {
		params = &ToolIDsParams{}
	}

	var result ToolIDs
	err := s.client.do(ctx, http.MethodGet, "experimental/tool/ids", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// List retrieves tools with JSON schema parameters for a provider/model
// GET /experimental/tool
func (s *ToolService) List(ctx context.Context, params *ToolListParams) (*ToolList, error) {
	if params == nil {
		return nil, errors.New("params is required")
	}

	var result ToolList
	err := s.client.do(ctx, http.MethodGet, "experimental/tool", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

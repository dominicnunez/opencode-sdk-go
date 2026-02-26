package opencode

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
)

// McpService handles communication with MCP-related endpoints
type McpService struct {
	client *Client
}

// McpStatus represents the MCP server status response
type McpStatus map[string]interface{}

// McpStatusParams contains parameters for the Status method
type McpStatusParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes McpStatusParams into URL query parameters
func (r McpStatusParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

// Status retrieves MCP server status
func (s *McpService) Status(ctx context.Context, params *McpStatusParams) (*McpStatus, error) {
	if params == nil {
		params = &McpStatusParams{}
	}

	var result McpStatus
	if err := s.client.do(ctx, http.MethodGet, "mcp", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

package opencode

import (
	"context"
	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
	"net/http"
	"net/url"
)

type PathService struct {
	client *Client
}

func (s *PathService) Get(ctx context.Context, params *PathGetParams) (*Path, error) {
	if params == nil {
		params = &PathGetParams{}
	}
	var result Path
	err := s.client.do(ctx, http.MethodGet, "path", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type Path struct {
	Config    string `json:"config"`
	Directory string `json:"directory"`
	State     string `json:"state"`
	Worktree  string `json:"worktree"`
}

type PathGetParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r PathGetParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

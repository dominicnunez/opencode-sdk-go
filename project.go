package opencode

import (
	"context"
	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
	"net/http"
	"net/url"
)

type ProjectService struct {
	client *Client
}

func (s *ProjectService) List(ctx context.Context, params *ProjectListParams) ([]Project, error) {
	if params == nil {
		params = &ProjectListParams{}
	}
	var result []Project
	err := s.client.do(ctx, http.MethodGet, "project", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) Current(ctx context.Context, params *ProjectCurrentParams) (*Project, error) {
	if params == nil {
		params = &ProjectCurrentParams{}
	}
	var result Project
	err := s.client.do(ctx, http.MethodGet, "project/current", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type Project struct {
	ID       string      `json:"id"`
	Time     ProjectTime `json:"time"`
	Worktree string      `json:"worktree"`
	Vcs      *ProjectVcs `json:"vcs,omitempty"`
}

type ProjectTime struct {
	Created     float64  `json:"created"`
	Initialized *float64 `json:"initialized,omitempty"`
}

type ProjectVcs string

const (
	ProjectVcsGit ProjectVcs = "git"
)

func (r ProjectVcs) IsKnown() bool {
	switch r {
	case ProjectVcsGit:
		return true
	}
	return false
}

type ProjectListParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r ProjectListParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type ProjectCurrentParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r ProjectCurrentParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

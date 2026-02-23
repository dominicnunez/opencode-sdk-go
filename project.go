package opencode

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/apijson"
	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
	"github.com/dominicnunez/opencode-sdk-go/internal/param"
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
	ID       string      `json:"id,required"`
	Time     ProjectTime `json:"time,required"`
	Worktree string      `json:"worktree,required"`
	Vcs      ProjectVcs  `json:"vcs"`
	JSON     projectJSON `json:"-"`
}

type projectJSON struct {
	ID          apijson.Field
	Time        apijson.Field
	Worktree    apijson.Field
	Vcs         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Project) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r projectJSON) RawJSON() string {
	return r.raw
}

type ProjectTime struct {
	Created     float64         `json:"created,required"`
	Initialized float64         `json:"initialized"`
	JSON        projectTimeJSON `json:"-"`
}

type projectTimeJSON struct {
	Created     apijson.Field
	Initialized apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProjectTime) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r projectTimeJSON) RawJSON() string {
	return r.raw
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
	Directory param.Field[string] `query:"directory"`
}

func (r ProjectListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type ProjectCurrentParams struct {
	Directory param.Field[string] `query:"directory"`
}

func (r ProjectCurrentParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

package opencode

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/apijson"
	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
	"github.com/dominicnunez/opencode-sdk-go/internal/param"
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
	Config    string   `json:"config,required"`
	Directory string   `json:"directory,required"`
	State     string   `json:"state,required"`
	Worktree  string   `json:"worktree,required"`
	JSON      pathJSON `json:"-"`
}

type pathJSON struct {
	Config      apijson.Field
	Directory   apijson.Field
	State       apijson.Field
	Worktree    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Path) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r pathJSON) RawJSON() string {
	return r.raw
}

type PathGetParams struct {
	Directory param.Field[string] `query:"directory"`
}

func (r PathGetParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

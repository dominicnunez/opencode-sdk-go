package opencode

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/apijson"
	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
	"github.com/dominicnunez/opencode-sdk-go/internal/param"
)

type CommandService struct {
	client *Client
}

func (s *CommandService) List(ctx context.Context, params *CommandListParams) ([]Command, error) {
	if params == nil {
		params = &CommandListParams{}
	}
	var result []Command
	err := s.client.do(ctx, http.MethodGet, "command", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type Command struct {
	Name        string      `json:"name,required"`
	Template    string      `json:"template,required"`
	Agent       string      `json:"agent"`
	Description string      `json:"description"`
	Model       string      `json:"model"`
	Subtask     bool        `json:"subtask"`
	JSON        commandJSON `json:"-"`
}

type commandJSON struct {
	Name        apijson.Field
	Template    apijson.Field
	Agent       apijson.Field
	Description apijson.Field
	Model       apijson.Field
	Subtask     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Command) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r commandJSON) RawJSON() string {
	return r.raw
}

type CommandListParams struct {
	Directory param.Field[string] `query:"directory"`
}

func (r CommandListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

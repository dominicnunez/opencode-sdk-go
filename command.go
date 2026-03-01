package opencode

import (
	"context"
	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
	"net/http"
	"net/url"
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
	Name        string `json:"name"`
	Template    string `json:"template"`
	Agent       string `json:"agent"`
	Description string `json:"description"`
	Model       string `json:"model"`
	Subtask     bool   `json:"subtask"`
}

type CommandListParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r CommandListParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

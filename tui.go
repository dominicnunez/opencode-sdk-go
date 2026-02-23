package opencode

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
)

type TuiService struct {
	client *Client
}

func (s *TuiService) AppendPrompt(ctx context.Context, params *TuiAppendPromptParams) (bool, error) {
	if params == nil {
		params = &TuiAppendPromptParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "tui/append-prompt", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *TuiService) ClearPrompt(ctx context.Context, params *TuiClearPromptParams) (bool, error) {
	if params == nil {
		params = &TuiClearPromptParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "tui/clear-prompt", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *TuiService) ExecuteCommand(ctx context.Context, params *TuiExecuteCommandParams) (bool, error) {
	if params == nil {
		params = &TuiExecuteCommandParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "tui/execute-command", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *TuiService) OpenHelp(ctx context.Context, params *TuiOpenHelpParams) (bool, error) {
	if params == nil {
		params = &TuiOpenHelpParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "tui/open-help", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *TuiService) OpenModels(ctx context.Context, params *TuiOpenModelsParams) (bool, error) {
	if params == nil {
		params = &TuiOpenModelsParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "tui/open-models", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *TuiService) OpenSessions(ctx context.Context, params *TuiOpenSessionsParams) (bool, error) {
	if params == nil {
		params = &TuiOpenSessionsParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "tui/open-sessions", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *TuiService) OpenThemes(ctx context.Context, params *TuiOpenThemesParams) (bool, error) {
	if params == nil {
		params = &TuiOpenThemesParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "tui/open-themes", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *TuiService) ShowToast(ctx context.Context, params *TuiShowToastParams) (bool, error) {
	if params == nil {
		params = &TuiShowToastParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "tui/show-toast", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *TuiService) SubmitPrompt(ctx context.Context, params *TuiSubmitPromptParams) (bool, error) {
	if params == nil {
		params = &TuiSubmitPromptParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "tui/submit-prompt", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

type ToastVariant string

const (
	ToastVariantInfo    ToastVariant = "info"
	ToastVariantSuccess ToastVariant = "success"
	ToastVariantWarning ToastVariant = "warning"
	ToastVariantError   ToastVariant = "error"
)

type TuiAppendPromptParams struct {
	Text      string  `json:"text"`
	Directory *string `query:"directory,omitempty"`
}

func (r TuiAppendPromptParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type TuiClearPromptParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r TuiClearPromptParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type TuiExecuteCommandParams struct {
	Command   string  `json:"command"`
	Directory *string `query:"directory,omitempty"`
}

func (r TuiExecuteCommandParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type TuiOpenHelpParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r TuiOpenHelpParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type TuiOpenModelsParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r TuiOpenModelsParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type TuiOpenSessionsParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r TuiOpenSessionsParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type TuiOpenThemesParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r TuiOpenThemesParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type TuiShowToastParams struct {
	Message   string       `json:"message"`
	Variant   ToastVariant `json:"variant"`
	Title     string       `json:"title,omitempty"`
	Directory *string      `query:"directory,omitempty"`
}

func (r TuiShowToastParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type TuiSubmitPromptParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r TuiSubmitPromptParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

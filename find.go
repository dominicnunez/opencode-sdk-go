// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package opencode

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/apijson"
	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
)

type FindService struct {
	client *Client
}

func (s *FindService) Files(ctx context.Context, params *FindFilesParams) ([]string, error) {
	if params == nil {
		params = &FindFilesParams{}
	}
	var result []string
	err := s.client.do(ctx, http.MethodGet, "find/file", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *FindService) Symbols(ctx context.Context, params *FindSymbolsParams) ([]Symbol, error) {
	if params == nil {
		params = &FindSymbolsParams{}
	}
	var result []Symbol
	err := s.client.do(ctx, http.MethodGet, "find/symbol", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *FindService) Text(ctx context.Context, params *FindTextParams) ([]FindTextResponse, error) {
	if params == nil {
		params = &FindTextParams{}
	}
	var result []FindTextResponse
	err := s.client.do(ctx, http.MethodGet, "find", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type Symbol struct {
	Kind     float64        `json:"kind,required"`
	Location SymbolLocation `json:"location,required"`
	Name     string         `json:"name,required"`
}

func (r *Symbol) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type SymbolLocation struct {
	Range SymbolLocationRange `json:"range,required"`
	Uri   string              `json:"uri,required"`
}

func (r *SymbolLocation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type SymbolLocationRange struct {
	End   SymbolLocationRangeEnd   `json:"end,required"`
	Start SymbolLocationRangeStart `json:"start,required"`
}

func (r *SymbolLocationRange) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type SymbolLocationRangeEnd struct {
	Character float64 `json:"character,required"`
	Line      float64 `json:"line,required"`
}

func (r *SymbolLocationRangeEnd) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type SymbolLocationRangeStart struct {
	Character float64 `json:"character,required"`
	Line      float64 `json:"line,required"`
}

func (r *SymbolLocationRangeStart) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type FindTextResponse struct {
	AbsoluteOffset float64                    `json:"absolute_offset,required"`
	LineNumber     float64                    `json:"line_number,required"`
	Lines          FindTextResponseLines      `json:"lines,required"`
	Path           FindTextResponsePath       `json:"path,required"`
	Submatches     []FindTextResponseSubmatch `json:"submatches,required"`
}

func (r *FindTextResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type FindTextResponseLines struct {
	Text string `json:"text,required"`
}

func (r *FindTextResponseLines) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type FindTextResponsePath struct {
	Text string `json:"text,required"`
}

func (r *FindTextResponsePath) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type FindTextResponseSubmatch struct {
	End   float64                         `json:"end,required"`
	Match FindTextResponseSubmatchesMatch `json:"match,required"`
	Start float64                         `json:"start,required"`
}

func (r *FindTextResponseSubmatch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type FindTextResponseSubmatchesMatch struct {
	Text string `json:"text,required"`
}

func (r *FindTextResponseSubmatchesMatch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type FindFilesParams struct {
	Query     string  `query:"query,required"`
	Directory *string `query:"directory,omitempty"`
}

func (r FindFilesParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type FindSymbolsParams struct {
	Query     string  `query:"query,required"`
	Directory *string `query:"directory,omitempty"`
}

func (r FindSymbolsParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type FindTextParams struct {
	Pattern   string  `query:"pattern,required"`
	Directory *string `query:"directory,omitempty"`
}

func (r FindTextParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

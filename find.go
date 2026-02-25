// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package opencode

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
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
	Kind     float64        `json:"kind"`
	Location SymbolLocation `json:"location"`
	Name     string         `json:"name"`
}

type SymbolLocation struct {
	Range SymbolLocationRange `json:"range"`
	Uri   string              `json:"uri"`
}

type SymbolLocationRange struct {
	End   SymbolLocationRangeEnd   `json:"end"`
	Start SymbolLocationRangeStart `json:"start"`
}

type SymbolLocationRangeEnd struct {
	Character float64 `json:"character"`
	Line      float64 `json:"line"`
}

type SymbolLocationRangeStart struct {
	Character float64 `json:"character"`
	Line      float64 `json:"line"`
}

type FindTextResponse struct {
	AbsoluteOffset float64                    `json:"absolute_offset"`
	LineNumber     float64                    `json:"line_number"`
	Lines          FindTextResponseLines      `json:"lines"`
	Path           FindTextResponsePath       `json:"path"`
	Submatches     []FindTextResponseSubmatch `json:"submatches"`
}

type FindTextResponseLines struct {
	Text string `json:"text"`
}

type FindTextResponsePath struct {
	Text string `json:"text"`
}

type FindTextResponseSubmatch struct {
	End   float64                         `json:"end"`
	Match FindTextResponseSubmatchesMatch `json:"match"`
	Start float64                         `json:"start"`
}

type FindTextResponseSubmatchesMatch struct {
	Text string `json:"text"`
}

type FindFilesParams struct {
	Query     string  `query:"query,required"`
	Directory *string `query:"directory,omitempty"`
}

func (r FindFilesParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type FindSymbolsParams struct {
	Query     string  `query:"query,required"`
	Directory *string `query:"directory,omitempty"`
}

func (r FindSymbolsParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type FindTextParams struct {
	Pattern   string  `query:"pattern,required"`
	Directory *string `query:"directory,omitempty"`
}

func (r FindTextParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

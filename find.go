package opencode

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
)

type FindService struct {
	client *Client
}

func (s *FindService) Files(ctx context.Context, params *FindFilesParams) ([]string, error) {
	if params == nil {
		return nil, errors.New("params is required")
	}
	if params.Query == "" {
		return nil, errors.New("missing required Query parameter")
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
		return nil, errors.New("params is required")
	}
	if params.Query == "" {
		return nil, errors.New("missing required Query parameter")
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
		return nil, errors.New("params is required")
	}
	if params.Pattern == "" {
		return nil, errors.New("missing required Pattern parameter")
	}
	var result []FindTextResponse
	err := s.client.do(ctx, http.MethodGet, "find", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SymbolKind represents the kind of a symbol as defined by LSP.
type SymbolKind int64

const (
	SymbolKindFile          SymbolKind = 1
	SymbolKindModule        SymbolKind = 2
	SymbolKindNamespace     SymbolKind = 3
	SymbolKindPackage       SymbolKind = 4
	SymbolKindClass         SymbolKind = 5
	SymbolKindMethod        SymbolKind = 6
	SymbolKindProperty      SymbolKind = 7
	SymbolKindField         SymbolKind = 8
	SymbolKindConstructor   SymbolKind = 9
	SymbolKindEnum          SymbolKind = 10
	SymbolKindInterface     SymbolKind = 11
	SymbolKindFunction      SymbolKind = 12
	SymbolKindVariable      SymbolKind = 13
	SymbolKindConstant      SymbolKind = 14
	SymbolKindString        SymbolKind = 15
	SymbolKindNumber        SymbolKind = 16
	SymbolKindBoolean       SymbolKind = 17
	SymbolKindArray         SymbolKind = 18
	SymbolKindObject        SymbolKind = 19
	SymbolKindKey           SymbolKind = 20
	SymbolKindNull          SymbolKind = 21
	SymbolKindEnumMember    SymbolKind = 22
	SymbolKindStruct        SymbolKind = 23
	SymbolKindEvent         SymbolKind = 24
	SymbolKindOperator      SymbolKind = 25
	SymbolKindTypeParameter SymbolKind = 26
)

func (r SymbolKind) IsKnown() bool {
	return r >= SymbolKindFile && r <= SymbolKindTypeParameter
}

type Symbol struct {
	Kind     SymbolKind     `json:"kind"`
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
	Character int64 `json:"character"`
	Line      int64 `json:"line"`
}

type SymbolLocationRangeStart struct {
	Character int64 `json:"character"`
	Line      int64 `json:"line"`
}

type FindTextResponse struct {
	AbsoluteOffset int64                      `json:"absolute_offset"`
	LineNumber     int64                      `json:"line_number"`
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
	End   int64                           `json:"end"`
	Match FindTextResponseSubmatchesMatch `json:"match"`
	Start int64                           `json:"start"`
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

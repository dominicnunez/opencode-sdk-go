// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package opencode

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/apijson"
	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
	"github.com/dominicnunez/opencode-sdk-go/internal/param"
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
	JSON     symbolJSON     `json:"-"`
}

type symbolJSON struct {
	Kind        apijson.Field
	Location    apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Symbol) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r symbolJSON) RawJSON() string {
	return r.raw
}

type SymbolLocation struct {
	Range SymbolLocationRange `json:"range,required"`
	Uri   string              `json:"uri,required"`
	JSON  symbolLocationJSON  `json:"-"`
}

type symbolLocationJSON struct {
	Range       apijson.Field
	Uri         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SymbolLocation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r symbolLocationJSON) RawJSON() string {
	return r.raw
}

type SymbolLocationRange struct {
	End   SymbolLocationRangeEnd   `json:"end,required"`
	Start SymbolLocationRangeStart `json:"start,required"`
	JSON  symbolLocationRangeJSON  `json:"-"`
}

type symbolLocationRangeJSON struct {
	End         apijson.Field
	Start       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SymbolLocationRange) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r symbolLocationRangeJSON) RawJSON() string {
	return r.raw
}

type SymbolLocationRangeEnd struct {
	Character float64                    `json:"character,required"`
	Line      float64                    `json:"line,required"`
	JSON      symbolLocationRangeEndJSON `json:"-"`
}

type symbolLocationRangeEndJSON struct {
	Character   apijson.Field
	Line        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SymbolLocationRangeEnd) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r symbolLocationRangeEndJSON) RawJSON() string {
	return r.raw
}

type SymbolLocationRangeStart struct {
	Character float64                      `json:"character,required"`
	Line      float64                      `json:"line,required"`
	JSON      symbolLocationRangeStartJSON `json:"-"`
}

type symbolLocationRangeStartJSON struct {
	Character   apijson.Field
	Line        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SymbolLocationRangeStart) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r symbolLocationRangeStartJSON) RawJSON() string {
	return r.raw
}

type FindTextResponse struct {
	AbsoluteOffset float64                    `json:"absolute_offset,required"`
	LineNumber     float64                    `json:"line_number,required"`
	Lines          FindTextResponseLines      `json:"lines,required"`
	Path           FindTextResponsePath       `json:"path,required"`
	Submatches     []FindTextResponseSubmatch `json:"submatches,required"`
	JSON           findTextResponseJSON       `json:"-"`
}

type findTextResponseJSON struct {
	AbsoluteOffset apijson.Field
	LineNumber     apijson.Field
	Lines          apijson.Field
	Path           apijson.Field
	Submatches     apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *FindTextResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r findTextResponseJSON) RawJSON() string {
	return r.raw
}

type FindTextResponseLines struct {
	Text string                    `json:"text,required"`
	JSON findTextResponseLinesJSON `json:"-"`
}

type findTextResponseLinesJSON struct {
	Text        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FindTextResponseLines) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r findTextResponseLinesJSON) RawJSON() string {
	return r.raw
}

type FindTextResponsePath struct {
	Text string                   `json:"text,required"`
	JSON findTextResponsePathJSON `json:"-"`
}

type findTextResponsePathJSON struct {
	Text        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FindTextResponsePath) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r findTextResponsePathJSON) RawJSON() string {
	return r.raw
}

type FindTextResponseSubmatch struct {
	End   float64                         `json:"end,required"`
	Match FindTextResponseSubmatchesMatch `json:"match,required"`
	Start float64                         `json:"start,required"`
	JSON  findTextResponseSubmatchJSON    `json:"-"`
}

type findTextResponseSubmatchJSON struct {
	End         apijson.Field
	Match       apijson.Field
	Start       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FindTextResponseSubmatch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r findTextResponseSubmatchJSON) RawJSON() string {
	return r.raw
}

type FindTextResponseSubmatchesMatch struct {
	Text string                              `json:"text,required"`
	JSON findTextResponseSubmatchesMatchJSON `json:"-"`
}

type findTextResponseSubmatchesMatchJSON struct {
	Text        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FindTextResponseSubmatchesMatch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r findTextResponseSubmatchesMatchJSON) RawJSON() string {
	return r.raw
}

type FindFilesParams struct {
	Query     param.Field[string] `query:"query,required"`
	Directory param.Field[string] `query:"directory"`
}

func (r FindFilesParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type FindSymbolsParams struct {
	Query     param.Field[string] `query:"query,required"`
	Directory param.Field[string] `query:"directory"`
}

func (r FindSymbolsParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type FindTextParams struct {
	Pattern   param.Field[string] `query:"pattern,required"`
	Directory param.Field[string] `query:"directory"`
}

func (r FindTextParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

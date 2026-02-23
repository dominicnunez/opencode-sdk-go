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

type AppService struct {
	client *Client
}

func (s *AppService) Log(ctx context.Context, params *AppLogParams) (bool, error) {
	if params == nil {
		params = &AppLogParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "log", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *AppService) Providers(ctx context.Context, params *AppProvidersParams) (*AppProvidersResponse, error) {
	if params == nil {
		params = &AppProvidersParams{}
	}
	var result AppProvidersResponse
	err := s.client.do(ctx, http.MethodGet, "config/providers", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type Model struct {
	ID           string                 `json:"id,required"`
	Attachment   bool                   `json:"attachment,required"`
	Cost         ModelCost              `json:"cost,required"`
	Limit        ModelLimit             `json:"limit,required"`
	Name         string                 `json:"name,required"`
	Options      map[string]interface{} `json:"options,required"`
	Reasoning    bool                   `json:"reasoning,required"`
	ReleaseDate  string                 `json:"release_date,required"`
	Temperature  bool                   `json:"temperature,required"`
	ToolCall     bool                   `json:"tool_call,required"`
	Experimental bool                   `json:"experimental"`
	Modalities   ModelModalities        `json:"modalities"`
	Provider     ModelProvider          `json:"provider"`
	Status       ModelStatus            `json:"status"`
	JSON         modelJSON              `json:"-"`
}

type modelJSON struct {
	ID           apijson.Field
	Attachment   apijson.Field
	Cost         apijson.Field
	Limit        apijson.Field
	Name         apijson.Field
	Options      apijson.Field
	Reasoning    apijson.Field
	ReleaseDate  apijson.Field
	Temperature  apijson.Field
	ToolCall     apijson.Field
	Experimental apijson.Field
	Modalities   apijson.Field
	Provider     apijson.Field
	Status       apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *Model) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelJSON) RawJSON() string {
	return r.raw
}

type ModelCost struct {
	Input      float64       `json:"input,required"`
	Output     float64       `json:"output,required"`
	CacheRead  float64       `json:"cache_read"`
	CacheWrite float64       `json:"cache_write"`
	JSON       modelCostJSON `json:"-"`
}

type modelCostJSON struct {
	Input       apijson.Field
	Output      apijson.Field
	CacheRead   apijson.Field
	CacheWrite  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ModelCost) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelCostJSON) RawJSON() string {
	return r.raw
}

type ModelLimit struct {
	Context float64        `json:"context,required"`
	Output  float64        `json:"output,required"`
	JSON    modelLimitJSON `json:"-"`
}

type modelLimitJSON struct {
	Context     apijson.Field
	Output      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ModelLimit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelLimitJSON) RawJSON() string {
	return r.raw
}

type ModelModalities struct {
	Input  []ModelModalitiesInput  `json:"input,required"`
	Output []ModelModalitiesOutput `json:"output,required"`
	JSON   modelModalitiesJSON     `json:"-"`
}

type modelModalitiesJSON struct {
	Input       apijson.Field
	Output      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ModelModalities) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelModalitiesJSON) RawJSON() string {
	return r.raw
}

type ModelModalitiesInput string

const (
	ModelModalitiesInputText  ModelModalitiesInput = "text"
	ModelModalitiesInputAudio ModelModalitiesInput = "audio"
	ModelModalitiesInputImage ModelModalitiesInput = "image"
	ModelModalitiesInputVideo ModelModalitiesInput = "video"
	ModelModalitiesInputPdf   ModelModalitiesInput = "pdf"
)

func (r ModelModalitiesInput) IsKnown() bool {
	switch r {
	case ModelModalitiesInputText, ModelModalitiesInputAudio, ModelModalitiesInputImage, ModelModalitiesInputVideo, ModelModalitiesInputPdf:
		return true
	}
	return false
}

type ModelModalitiesOutput string

const (
	ModelModalitiesOutputText  ModelModalitiesOutput = "text"
	ModelModalitiesOutputAudio ModelModalitiesOutput = "audio"
	ModelModalitiesOutputImage ModelModalitiesOutput = "image"
	ModelModalitiesOutputVideo ModelModalitiesOutput = "video"
	ModelModalitiesOutputPdf   ModelModalitiesOutput = "pdf"
)

func (r ModelModalitiesOutput) IsKnown() bool {
	switch r {
	case ModelModalitiesOutputText, ModelModalitiesOutputAudio, ModelModalitiesOutputImage, ModelModalitiesOutputVideo, ModelModalitiesOutputPdf:
		return true
	}
	return false
}

type ModelProvider struct {
	Npm  string            `json:"npm,required"`
	JSON modelProviderJSON `json:"-"`
}

type modelProviderJSON struct {
	Npm         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ModelProvider) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelProviderJSON) RawJSON() string {
	return r.raw
}

type ModelStatus string

const (
	ModelStatusAlpha ModelStatus = "alpha"
	ModelStatusBeta  ModelStatus = "beta"
)

func (r ModelStatus) IsKnown() bool {
	switch r {
	case ModelStatusAlpha, ModelStatusBeta:
		return true
	}
	return false
}

type Provider struct {
	ID     string           `json:"id,required"`
	Env    []string         `json:"env,required"`
	Models map[string]Model `json:"models,required"`
	Name   string           `json:"name,required"`
	API    string           `json:"api"`
	Npm    string           `json:"npm"`
	JSON   providerJSON     `json:"-"`
}

type providerJSON struct {
	ID          apijson.Field
	Env         apijson.Field
	Models      apijson.Field
	Name        apijson.Field
	API         apijson.Field
	Npm         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Provider) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r providerJSON) RawJSON() string {
	return r.raw
}

type AppProvidersResponse struct {
	Default   map[string]string        `json:"default,required"`
	Providers []Provider               `json:"providers,required"`
	JSON      appProvidersResponseJSON `json:"-"`
}

type appProvidersResponseJSON struct {
	Default     apijson.Field
	Providers   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AppProvidersResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r appProvidersResponseJSON) RawJSON() string {
	return r.raw
}

type AppLogParams struct {
	Level     param.Field[AppLogParamsLevel] `json:"level,required"`
	Message   param.Field[string]            `json:"message,required"`
	Service   param.Field[string]            `json:"service,required"`
	Directory param.Field[string]            `query:"directory"`
	Extra     param.Field[map[string]interface{}] `json:"extra"`
}

func (r AppLogParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r AppLogParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AppLogParamsLevel string

const (
	AppLogParamsLevelDebug AppLogParamsLevel = "debug"
	AppLogParamsLevelInfo  AppLogParamsLevel = "info"
	AppLogParamsLevelError AppLogParamsLevel = "error"
	AppLogParamsLevelWarn  AppLogParamsLevel = "warn"
)

func (r AppLogParamsLevel) IsKnown() bool {
	switch r {
	case AppLogParamsLevelDebug, AppLogParamsLevelInfo, AppLogParamsLevelError, AppLogParamsLevelWarn:
		return true
	}
	return false
}

type AppProvidersParams struct {
	Directory param.Field[string] `query:"directory"`
}

func (r AppProvidersParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

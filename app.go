package opencode

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
)

type AppService struct {
	client *Client
}

func (s *AppService) Log(ctx context.Context, params *AppLogParams) (bool, error) {
	if params == nil {
		return false, errors.New("params is required")
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "log", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

func (r LogLevel) IsKnown() bool {
	switch r {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return true
	}
	return false
}

type Model struct {
	ID           string                 `json:"id"`
	Attachment   bool                   `json:"attachment"`
	Cost         ModelCost              `json:"cost"`
	Experimental bool                   `json:"experimental,omitempty"`
	Limit        ModelLimit             `json:"limit"`
	Modalities   ModelModalities        `json:"modalities,omitempty"`
	Name         string                 `json:"name"`
	Options      map[string]interface{} `json:"options"`
	Provider     ModelProvider          `json:"provider,omitempty"`
	Reasoning    bool                   `json:"reasoning"`
	ReleaseDate  string                 `json:"release_date"`
	Status       ModelStatus            `json:"status,omitempty"`
	Temperature  bool                   `json:"temperature"`
	ToolCall     bool                   `json:"tool_call"`
}

type ModelCost struct {
	Input      float64 `json:"input"`
	Output     float64 `json:"output"`
	CacheRead  float64 `json:"cache_read,omitempty"`
	CacheWrite float64 `json:"cache_write,omitempty"`
}

type ModelLimit struct {
	Context float64 `json:"context"`
	Output  float64 `json:"output"`
}

type ModelModalities struct {
	Input  []ModelModalityInput  `json:"input"`
	Output []ModelModalityOutput `json:"output"`
}

type ModelModalityInput string

const (
	ModelModalityInputText  ModelModalityInput = "text"
	ModelModalityInputAudio ModelModalityInput = "audio"
	ModelModalityInputImage ModelModalityInput = "image"
	ModelModalityInputVideo ModelModalityInput = "video"
	ModelModalityInputPdf   ModelModalityInput = "pdf"
)

func (r ModelModalityInput) IsKnown() bool {
	switch r {
	case ModelModalityInputText, ModelModalityInputAudio, ModelModalityInputImage, ModelModalityInputVideo, ModelModalityInputPdf:
		return true
	}
	return false
}

type ModelModalityOutput string

const (
	ModelModalityOutputText  ModelModalityOutput = "text"
	ModelModalityOutputAudio ModelModalityOutput = "audio"
	ModelModalityOutputImage ModelModalityOutput = "image"
	ModelModalityOutputVideo ModelModalityOutput = "video"
	ModelModalityOutputPdf   ModelModalityOutput = "pdf"
)

func (r ModelModalityOutput) IsKnown() bool {
	switch r {
	case ModelModalityOutputText, ModelModalityOutputAudio, ModelModalityOutputImage, ModelModalityOutputVideo, ModelModalityOutputPdf:
		return true
	}
	return false
}

type ModelProvider struct {
	Npm string `json:"npm"`
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
	ID     string           `json:"id"`
	API    string           `json:"api,omitempty"`
	Env    []string         `json:"env"`
	Models map[string]Model `json:"models"`
	Name   string           `json:"name"`
	Npm    string           `json:"npm,omitempty"`
}

type AppLogParams struct {
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Service   string                 `json:"service"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
	Directory *string                `json:"-" query:"directory,omitempty"`
}

func (r AppLogParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

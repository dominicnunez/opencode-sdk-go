package opencode

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
)

type FileService struct {
	client *Client
}

func (s *FileService) List(ctx context.Context, params *FileListParams) ([]FileNode, error) {
	if params == nil {
		return nil, errors.New("params is required")
	}
	if params.Path == "" {
		return nil, errors.New("missing required Path parameter")
	}
	var result []FileNode
	err := s.client.do(ctx, http.MethodGet, "file", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *FileService) Read(ctx context.Context, params *FileReadParams) (*FileReadResponse, error) {
	if params == nil {
		return nil, errors.New("params is required")
	}
	if params.Path == "" {
		return nil, errors.New("missing required Path parameter")
	}
	var result FileReadResponse
	err := s.client.do(ctx, http.MethodGet, "file/content", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *FileService) Status(ctx context.Context, params *FileStatusParams) ([]File, error) {
	if params == nil {
		params = &FileStatusParams{}
	}
	var result []File
	err := s.client.do(ctx, http.MethodGet, "file/status", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type File struct {
	Added   int64      `json:"added"`
	Path    string     `json:"path"`
	Removed int64      `json:"removed"`
	Status  FileStatus `json:"status"`
}

type FileStatus string

const (
	FileStatusAdded    FileStatus = "added"
	FileStatusDeleted  FileStatus = "deleted"
	FileStatusModified FileStatus = "modified"
)

func (r FileStatus) IsKnown() bool {
	switch r {
	case FileStatusAdded, FileStatusDeleted, FileStatusModified:
		return true
	}
	return false
}

type FileNode struct {
	Absolute string       `json:"absolute"`
	Ignored  bool         `json:"ignored"`
	Name     string       `json:"name"`
	Path     string       `json:"path"`
	Type     FileNodeType `json:"type"`
}

type FileNodeType string

const (
	FileNodeTypeFile      FileNodeType = "file"
	FileNodeTypeDirectory FileNodeType = "directory"
)

func (r FileNodeType) IsKnown() bool {
	switch r {
	case FileNodeTypeFile, FileNodeTypeDirectory:
		return true
	}
	return false
}

type FileReadResponse struct {
	Content  string                   `json:"content"`
	Type     FileReadResponseType     `json:"type"`
	Diff     string                   `json:"diff,omitempty"`
	Encoding FileReadResponseEncoding `json:"encoding,omitempty"`
	MimeType string                   `json:"mimeType,omitempty"`
	Patch    FileReadResponsePatch    `json:"patch,omitempty"`
}

type FileReadResponseType string

const (
	FileReadResponseTypeText FileReadResponseType = "text"
)

func (r FileReadResponseType) IsKnown() bool {
	switch r {
	case FileReadResponseTypeText:
		return true
	}
	return false
}

type FileReadResponseEncoding string

const (
	FileReadResponseEncodingBase64 FileReadResponseEncoding = "base64"
)

func (r FileReadResponseEncoding) IsKnown() bool {
	switch r {
	case FileReadResponseEncodingBase64:
		return true
	}
	return false
}

type FileReadResponsePatch struct {
	Hunks       []FileReadResponsePatchHunk `json:"hunks"`
	NewFileName string                      `json:"newFileName"`
	OldFileName string                      `json:"oldFileName"`
	Index       string                      `json:"index,omitempty"`
	NewHeader   string                      `json:"newHeader,omitempty"`
	OldHeader   string                      `json:"oldHeader,omitempty"`
}

type FileReadResponsePatchHunk struct {
	Lines    []string `json:"lines"`
	NewLines int64    `json:"newLines"`
	NewStart int64    `json:"newStart"`
	OldLines int64    `json:"oldLines"`
	OldStart int64    `json:"oldStart"`
}

type FileListParams struct {
	Path      string  `query:"path,required"`
	Directory *string `query:"directory,omitempty"`
}

func (r FileListParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type FileReadParams struct {
	Path      string  `query:"path,required"`
	Directory *string `query:"directory,omitempty"`
}

func (r FileReadParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type FileStatusParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r FileStatusParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

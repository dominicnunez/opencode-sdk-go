// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package opencode

import (
	"context"
	"net/http"
	"net/url"
	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
)

type FileService struct {
	client *Client
}

func (s *FileService) List(ctx context.Context, params *FileListParams) ([]FileNode, error) {
	if params == nil {
		params = &FileListParams{}
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
		params = &FileReadParams{}
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
	Added   int64      `json:"added,required"`
	Path    string     `json:"path,required"`
	Removed int64      `json:"removed,required"`
	Status  FileStatus `json:"status,required"`
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
	Absolute string       `json:"absolute,required"`
	Ignored  bool         `json:"ignored,required"`
	Name     string       `json:"name,required"`
	Path     string       `json:"path,required"`
	Type     FileNodeType `json:"type,required"`
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
	Content  string                   `json:"content,required"`
	Type     FileReadResponseType     `json:"type,required"`
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
	Hunks       []FileReadResponsePatchHunk `json:"hunks,required"`
	NewFileName string                      `json:"newFileName,required"`
	OldFileName string                      `json:"oldFileName,required"`
	Index       string                      `json:"index,omitempty"`
	NewHeader   string                      `json:"newHeader,omitempty"`
	OldHeader   string                      `json:"oldHeader,omitempty"`
}

type FileReadResponsePatchHunk struct {
	Lines    []string `json:"lines,required"`
	NewLines float64  `json:"newLines,required"`
	NewStart float64  `json:"newStart,required"`
	OldLines float64  `json:"oldLines,required"`
	OldStart float64  `json:"oldStart,required"`
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

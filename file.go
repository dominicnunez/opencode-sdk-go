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
	JSON    fileJSON   `json:"-"`
}

type fileJSON struct {
	Added       apijson.Field
	Path        apijson.Field
	Removed     apijson.Field
	Status      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *File) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileJSON) RawJSON() string {
	return r.raw
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
	JSON     fileNodeJSON `json:"-"`
}

type fileNodeJSON struct {
	Absolute    apijson.Field
	Ignored     apijson.Field
	Name        apijson.Field
	Path        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileNode) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileNodeJSON) RawJSON() string {
	return r.raw
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
	Diff     string                   `json:"diff"`
	Encoding FileReadResponseEncoding `json:"encoding"`
	MimeType string                   `json:"mimeType"`
	Patch    FileReadResponsePatch    `json:"patch"`
	JSON     fileReadResponseJSON     `json:"-"`
}

type fileReadResponseJSON struct {
	Content     apijson.Field
	Type        apijson.Field
	Diff        apijson.Field
	Encoding    apijson.Field
	MimeType    apijson.Field
	Patch       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileReadResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileReadResponseJSON) RawJSON() string {
	return r.raw
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
	Index       string                      `json:"index"`
	NewHeader   string                      `json:"newHeader"`
	OldHeader   string                      `json:"oldHeader"`
	JSON        fileReadResponsePatchJSON   `json:"-"`
}

type fileReadResponsePatchJSON struct {
	Hunks       apijson.Field
	NewFileName apijson.Field
	OldFileName apijson.Field
	Index       apijson.Field
	NewHeader   apijson.Field
	OldHeader   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileReadResponsePatch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileReadResponsePatchJSON) RawJSON() string {
	return r.raw
}

type FileReadResponsePatchHunk struct {
	Lines    []string                      `json:"lines,required"`
	NewLines float64                       `json:"newLines,required"`
	NewStart float64                       `json:"newStart,required"`
	OldLines float64                       `json:"oldLines,required"`
	OldStart float64                       `json:"oldStart,required"`
	JSON     fileReadResponsePatchHunkJSON `json:"-"`
}

type fileReadResponsePatchHunkJSON struct {
	Lines       apijson.Field
	NewLines    apijson.Field
	NewStart    apijson.Field
	OldLines    apijson.Field
	OldStart    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileReadResponsePatchHunk) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileReadResponsePatchHunkJSON) RawJSON() string {
	return r.raw
}

type FileListParams struct {
	Path      param.Field[string] `query:"path,required"`
	Directory param.Field[string] `query:"directory"`
}

func (r FileListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type FileReadParams struct {
	Path      param.Field[string] `query:"path,required"`
	Directory param.Field[string] `query:"directory"`
}

func (r FileReadParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type FileStatusParams struct {
	Directory param.Field[string] `query:"directory"`
}

func (r FileStatusParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

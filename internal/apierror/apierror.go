// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package apierror

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/anomalyco/opencode-sdk-go/internal/apijson"
)

// Error represents an error that originates from the API, i.e. when a request is
// made and the API returns a response with a HTTP status code. Other errors are
// not wrapped by this SDK.
type Error struct {
	JSON       errorJSON `json:"-"`
	StatusCode int
	Request    *http.Request
	Response   *http.Response
}

// errorJSON contains the JSON metadata for the struct [Error]
type errorJSON struct {
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Error) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r errorJSON) RawJSON() string {
	return r.raw
}

func (r *Error) Error() string {
	var method, url string
	var statusCode int
	var statusText string

	if r.Request != nil {
		method = r.Request.Method
		if r.Request.URL != nil {
			url = r.Request.URL.String()
		}
	}
	if r.Response != nil {
		statusCode = r.Response.StatusCode
		statusText = http.StatusText(r.Response.StatusCode)
	}

	return fmt.Sprintf("%s \"%s\": %d %s %s", method, url, statusCode, statusText, r.JSON.RawJSON())
}

func (r *Error) DumpRequest(body bool) []byte {
	if r.Request == nil {
		return nil
	}
	if r.Request.GetBody != nil {
		r.Request.Body, _ = r.Request.GetBody()
	}
	out, _ := httputil.DumpRequestOut(r.Request, body)
	return out
}

func (r *Error) DumpResponse(body bool) []byte {
	if r.Response == nil {
		return nil
	}
	out, _ := httputil.DumpResponse(r.Response, body)
	return out
}

# Misreads

> Findings where the audit misread the code or described behavior that doesn't occur.
> Managed by sfk willie. Follow the entry format below.
>
> Entry format:
> ### Plain language description
> **Location:** `file/path:line` — optional context
> **Date:** YYYY-MM-DD
> **Reason:** Explanation (can be multiple lines)

### custom roundtripper returning nil response and nil error does not panic doRaw

**Location:** `client.go:473` — `doRaw` success-status branch
**Date:** 2026-03-02

**Reason:** The client stores `httpClient` as `*http.Client` (`client.go:53`) and only accepts `*http.Client` in `WithHTTPClient` (`client.go:184`). Go's standard library `(*http.Client).Do` converts a non-conforming transport `(nil, nil)` result into a non-nil error (`http: RoundTripper implementation ... returned a nil *Response with a nil error`) before returning to `doRaw`. Because `lastErr` is non-nil, `resp.StatusCode` is not dereferenced in that path. This behavior is now covered by `TestClientDo_NilResponseWithoutErrorFromTransport` in `client_do_test.go`.

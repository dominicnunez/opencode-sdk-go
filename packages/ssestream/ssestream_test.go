package ssestream

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
)

type mockDecoder struct{}

func (m *mockDecoder) Event() Event { return Event{} }
func (m *mockDecoder) Next() bool   { return false }
func (m *mockDecoder) Close() error { return nil }
func (m *mockDecoder) Err() error   { return nil }

// saveAndRestoreDecoders snapshots the global decoder map and restores it
// when the test completes, preventing registration leaks between tests.
func saveAndRestoreDecoders(t *testing.T) {
	t.Helper()
	decoderTypesMu.Lock()
	snapshot := make(map[string]func(io.ReadCloser) Decoder, len(decoderTypes))
	for k, v := range decoderTypes {
		snapshot[k] = v
	}
	decoderTypesMu.Unlock()
	t.Cleanup(func() {
		decoderTypesMu.Lock()
		decoderTypes = snapshot
		decoderTypesMu.Unlock()
	})
}

func TestRegisterDecoderConcurrent(t *testing.T) {
	saveAndRestoreDecoders(t)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			contentType := "application/x-test-stream"
			if i%2 == 0 {
				contentType = "application/x-alt-stream"
			}
			RegisterDecoder(contentType, func(rc io.ReadCloser) Decoder {
				return &mockDecoder{}
			})
		}(i)
	}
	wg.Wait()
}

// newSSEResponse wraps a raw SSE string into an *http.Response with
// content-type "text/event-stream" and an io.NopCloser body.
func newSSEResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestEventStreamDecoder_NoTrailingBlankLine_DispatchesFinalEvent(t *testing.T) {
	// Stream ends with data but no trailing blank line after the last event.
	raw := "event: message\ndata: {\"ok\":true}\n"
	dec := NewDecoder(newSSEResponse(raw))
	defer func() { _ = dec.Close() }()

	if !dec.Next() {
		t.Fatal("expected Next() to return true for buffered event at EOF")
	}

	evt := dec.Event()
	if evt.Type != "message" {
		t.Fatalf("expected event type %q, got %q", "message", evt.Type)
	}
	if string(evt.Data) != `{"ok":true}` {
		t.Fatalf("expected data %q, got %q", `{"ok":true}`, string(evt.Data))
	}

	if dec.Next() {
		t.Fatal("expected Next() to return false after final event")
	}
	if err := dec.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEventStreamDecoder_TrailingBlankLine_DispatchesEvent(t *testing.T) {
	// Baseline: stream with a proper trailing blank line dispatches the event.
	raw := "event: done\ndata: {\"id\":1}\n\n"
	dec := NewDecoder(newSSEResponse(raw))
	defer func() { _ = dec.Close() }()

	if !dec.Next() {
		t.Fatal("expected Next() to return true")
	}

	evt := dec.Event()
	if evt.Type != "done" {
		t.Fatalf("expected event type %q, got %q", "done", evt.Type)
	}
	if string(evt.Data) != `{"id":1}` {
		t.Fatalf("expected data %q, got %q", `{"id":1}`, string(evt.Data))
	}

	if dec.Next() {
		t.Fatal("expected Next() to return false after all events consumed")
	}
	if err := dec.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEventStreamDecoder_EmptyStream_ReturnsFalse(t *testing.T) {
	// EOF with no data at all should return false immediately.
	dec := NewDecoder(newSSEResponse(""))
	defer func() { _ = dec.Close() }()

	if dec.Next() {
		t.Fatal("expected Next() to return false on empty stream")
	}
	if err := dec.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEventStreamDecoder_EmptyDataFieldAtEOF_DispatchesEvent(t *testing.T) {
	// Stream ends with "data: " (data field with empty value after colon+space).
	// Per SSE spec, "data:" seen = hasData = true, so event should dispatch at EOF
	// even though the data value is empty.
	raw := "event: foo\ndata: \n"
	dec := NewDecoder(newSSEResponse(raw))
	defer func() { _ = dec.Close() }()

	if !dec.Next() {
		t.Fatal("expected Next() to return true for event with empty data field")
	}

	evt := dec.Event()
	if evt.Type != "foo" {
		t.Fatalf("expected event type %q, got %q", "foo", evt.Type)
	}
	if len(evt.Data) != 0 {
		t.Fatalf("expected empty data, got %q", string(evt.Data))
	}

	if dec.Next() {
		t.Fatal("expected Next() to return false after final event")
	}
	if err := dec.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEventStreamDecoder_EventTypeOnlyNoDataField_NoDispatchAtEOF(t *testing.T) {
	// Stream ends with event type only, no "data:" field at all.
	// Per SSE spec, no data field seen = hasData = false, so no event should dispatch.
	raw := "event: foo\n"
	dec := NewDecoder(newSSEResponse(raw))
	defer func() { _ = dec.Close() }()

	if dec.Next() {
		t.Fatal("expected Next() to return false when no data field present")
	}
	if err := dec.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEventStreamDecoder_EventTypeOnlyWithBlankLine_DropsEvent(t *testing.T) {
	// An event block with only "event:" and no "data:" field, followed by
	// a blank line, should be silently dropped per SSE spec §9.2.6.
	// The next valid event (with data) should still be dispatched.
	raw := "event: ping\n\nevent: msg\ndata: {\"ok\":true}\n\n"
	dec := NewDecoder(newSSEResponse(raw))
	defer func() { _ = dec.Close() }()

	if !dec.Next() {
		t.Fatal("expected Next() to return true for the second event")
	}
	evt := dec.Event()
	if evt.Type != "msg" {
		t.Fatalf("expected event type %q, got %q", "msg", evt.Type)
	}
	if string(evt.Data) != `{"ok":true}` {
		t.Fatalf("expected data %q, got %q", `{"ok":true}`, string(evt.Data))
	}

	if dec.Next() {
		t.Fatal("expected Next() to return false after all events consumed")
	}
	if err := dec.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEventStreamDecoder_EventTypeResetAfterDrop(t *testing.T) {
	// After dropping an event-only block (no data), the event type must
	// be reset. A subsequent data-only block should have an empty type.
	raw := "event: ping\n\ndata: hello\n\n"
	dec := NewDecoder(newSSEResponse(raw))
	defer func() { _ = dec.Close() }()

	if !dec.Next() {
		t.Fatal("expected Next() to return true")
	}
	evt := dec.Event()
	if evt.Type != "" {
		t.Fatalf("expected empty event type after reset, got %q", evt.Type)
	}
	if string(evt.Data) != "hello" {
		t.Fatalf("expected data %q, got %q", "hello", string(evt.Data))
	}
}

func TestStream_SkipsEmptyDataEvents(t *testing.T) {
	// An empty-data event (e.g. keep-alive) between two valid events should
	// be silently skipped instead of killing the stream with an unmarshal error.
	raw := "event: ping\ndata: \n\nevent: msg\ndata: {\"ok\":true}\n\n"
	dec := NewDecoder(newSSEResponse(raw))

	type payload struct {
		OK bool `json:"ok"`
	}
	stream := NewStream[payload](dec, nil)
	defer func() { _ = stream.Close() }()

	if !stream.Next() {
		t.Fatalf("expected Next() to return true, got false; err=%v", stream.Err())
	}
	if !stream.Current().OK {
		t.Fatal("expected Current().OK to be true")
	}

	if stream.Next() {
		t.Fatal("expected Next() to return false after all events consumed")
	}
	if stream.Err() != nil {
		t.Fatalf("unexpected error: %v", stream.Err())
	}
}

func TestStream_DoubleClose_BothReturnNil(t *testing.T) {
	// After first Close(), decoder is set to nil. Second Close() should
	// return nil (not panic or error) since decoder is nil.
	dec := &mockDecoder{}
	stream := NewStream[interface{}](dec, nil)

	// First Close() should succeed.
	err1 := stream.Close()
	if err1 != nil {
		t.Fatalf("expected first Close() to return nil, got %v", err1)
	}

	// Verify decoder is nil after first Close().
	if stream.decoder != nil {
		t.Fatal("expected decoder to be nil after Close()")
	}

	// Second Close() should also return nil (not panic).
	err2 := stream.Close()
	if err2 != nil {
		t.Fatalf("expected second Close() to return nil, got %v", err2)
	}
}

func TestStream_DecoderErrorPropagation(t *testing.T) {
	readErr := errors.New("connection reset by peer")

	dec := &errorDecoder{err: readErr}
	stream := NewStream[interface{}](dec, nil)

	if stream.Next() {
		t.Fatal("expected Next() to return false when decoder errors")
	}

	if !errors.Is(stream.Err(), readErr) {
		t.Fatalf("expected stream.Err() to be %v, got %v", readErr, stream.Err())
	}
}

// errorDecoder simulates a decoder that fails mid-stream with a read error.
type errorDecoder struct {
	err error
}

func (d *errorDecoder) Event() Event { return Event{} }
func (d *errorDecoder) Next() bool   { return false }
func (d *errorDecoder) Close() error { return nil }
func (d *errorDecoder) Err() error   { return d.err }

func TestRegisterDecoder_LookupByContentType(t *testing.T) {
	saveAndRestoreDecoders(t)
	customType := "application/x-custom-stream"
	called := false

	RegisterDecoder(customType, func(rc io.ReadCloser) Decoder {
		called = true
		return &mockDecoder{}
	})

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{customType}},
		Body:       io.NopCloser(strings.NewReader("")),
	}

	dec := NewDecoder(resp)
	if dec == nil {
		t.Fatal("expected non-nil decoder")
	}
	if !called {
		t.Error("expected custom decoder factory to be called for matching content-type")
	}
}

func TestNewDecoder_UnknownContentType_FallsBackToSSE(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/x-unknown-stream-type"}},
		Body:       io.NopCloser(strings.NewReader("event: test\ndata: {\"ok\":true}\n\n")),
	}

	dec := NewDecoder(resp)
	if dec == nil {
		t.Fatal("expected non-nil decoder for unknown content-type")
	}

	// Should fall back to SSE decoder and parse the event
	if !dec.Next() {
		t.Fatal("expected SSE fallback decoder to parse the event")
	}
	evt := dec.Event()
	if evt.Type != "test" {
		t.Errorf("expected event type %q, got %q", "test", evt.Type)
	}
}

func TestStream_CloseMiddleOfIteration_NextReturnsErrNilDecoder(t *testing.T) {
	// Closing a stream mid-iteration (from the same goroutine) should cause
	// subsequent Next() calls to return false with ErrNilDecoder, not panic.
	raw := "event: a\ndata: {\"ok\":true}\n\nevent: b\ndata: {\"ok\":true}\n\n"
	dec := NewDecoder(newSSEResponse(raw))

	type payload struct {
		OK bool `json:"ok"`
	}
	stream := NewStream[payload](dec, nil)

	// Consume first event
	if !stream.Next() {
		t.Fatalf("expected first Next() to return true, err=%v", stream.Err())
	}

	// Close mid-iteration
	if err := stream.Close(); err != nil {
		t.Fatalf("Close() returned error: %v", err)
	}

	// Next() after Close should return false with ErrNilDecoder
	if stream.Next() {
		t.Fatal("expected Next() to return false after Close()")
	}
	if !errors.Is(stream.Err(), ErrNilDecoder) {
		t.Fatalf("expected ErrNilDecoder, got %v", stream.Err())
	}
}

func TestEventStreamDecoder_ConnectionDropMidEvent(t *testing.T) {
	// Simulate connection drop mid-event: reader returns partial data then error.
	connErr := errors.New("connection reset by peer")
	r := &failingReader{
		data: "event: msg\ndata: {\"partial\":",
		err:  connErr,
	}
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:       io.NopCloser(r),
	}

	dec := NewDecoder(resp)
	defer func() { _ = dec.Close() }()

	if dec.Next() {
		t.Fatal("expected Next() to return false on connection drop mid-event")
	}
	if dec.Err() == nil {
		t.Fatal("expected non-nil error from decoder after connection drop")
	}
}

// failingReader returns data on first read, then an error.
type failingReader struct {
	data string
	err  error
	read bool
}

func (r *failingReader) Read(p []byte) (int, error) {
	if !r.read {
		r.read = true
		n := copy(p, r.data)
		return n, nil
	}
	return 0, r.err
}

func TestEventStreamDecoder_InterleavedComments(t *testing.T) {
	// SSE comments (lines starting with ":") should be silently ignored.
	raw := ": keep-alive\nevent: msg\n: another comment\ndata: {\"ok\":true}\n\n"
	dec := NewDecoder(newSSEResponse(raw))
	defer func() { _ = dec.Close() }()

	if !dec.Next() {
		t.Fatal("expected Next() to return true")
	}
	evt := dec.Event()
	if evt.Type != "msg" {
		t.Errorf("expected event type %q, got %q", "msg", evt.Type)
	}
	if string(evt.Data) != `{"ok":true}` {
		t.Errorf("expected data %q, got %q", `{"ok":true}`, string(evt.Data))
	}
}

func TestEventStreamDecoder_TokenExceedsBufferLimit(t *testing.T) {
	// A single line exceeding the scanner buffer causes a scanner error.
	// Use a small custom scanner to avoid allocating 32MB in tests.
	const smallLimit = 256
	body := "data: " + strings.Repeat("x", smallLimit+1) + "\n\n"
	r := io.NopCloser(strings.NewReader(body))

	scn := bufio.NewScanner(r)
	scn.Buffer(nil, smallLimit)
	dec := &eventStreamDecoder{rc: r, scn: scn}

	if dec.Next() {
		t.Fatal("expected Next() to return false when token exceeds buffer")
	}
	if dec.Err() == nil {
		t.Fatal("expected scanner error for oversized token")
	}
}

func TestEventStreamDecoder_IncompleteEventNoBlankLine(t *testing.T) {
	// Stream has data but connection closes without a blank line separator.
	// The decoder should still dispatch the buffered event at EOF.
	raw := "event: update\ndata: {\"id\":1}"
	dec := NewDecoder(newSSEResponse(raw))
	defer func() { _ = dec.Close() }()

	if !dec.Next() {
		t.Fatal("expected Next() to return true for event at EOF without trailing newline")
	}
	evt := dec.Event()
	if evt.Type != "update" {
		t.Errorf("expected event type %q, got %q", "update", evt.Type)
	}
	if string(evt.Data) != `{"id":1}` {
		t.Errorf("expected data %q, got %q", `{"id":1}`, string(evt.Data))
	}
}

package ssestream

import (
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

func TestRegisterDecoderConcurrent(t *testing.T) {
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

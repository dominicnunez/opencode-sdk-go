// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package ssestream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
)

// Multiplier for bufio.MaxScanTokenSize to handle large SSE events.
// At 9x the default (64KB), provides ~576KB buffer for event data.
const sseBufferMultiplier = 9

// ErrNilDecoder is returned by Stream.Next when the SSE decoder is nil,
// indicating the response body was not available.
var ErrNilDecoder = errors.New("ssestream: decoder is nil")

type Decoder interface {
	Event() Event
	Next() bool
	Close() error
	Err() error
}

func NewDecoder(res *http.Response) Decoder {
	if res == nil || res.Body == nil {
		return nil
	}

	var decoder Decoder
	contentType := res.Header.Get("content-type")
	decoderTypesMu.RLock()
	t, ok := decoderTypes[contentType]
	decoderTypesMu.RUnlock()
	if ok {
		decoder = t(res.Body)
	} else {
		scn := bufio.NewScanner(res.Body)
		scn.Buffer(nil, bufio.MaxScanTokenSize<<sseBufferMultiplier)
		decoder = &eventStreamDecoder{rc: res.Body, scn: scn}
	}
	return decoder
}

var (
	decoderTypes   = map[string](func(io.ReadCloser) Decoder){}
	decoderTypesMu sync.RWMutex
)

func RegisterDecoder(contentType string, decoder func(io.ReadCloser) Decoder) {
	decoderTypesMu.Lock()
	decoderTypes[strings.ToLower(contentType)] = decoder
	decoderTypesMu.Unlock()
}

type Event struct {
	Type string
	Data []byte
}

// A base implementation of a Decoder for text/event-stream.
type eventStreamDecoder struct {
	evt Event
	rc  io.ReadCloser
	scn *bufio.Scanner
	err error
}

func (s *eventStreamDecoder) Next() bool {
	if s.err != nil {
		return false
	}

	event := ""
	data := bytes.NewBuffer(nil)

	for s.scn.Scan() {
		txt := s.scn.Bytes()

		// Dispatch event on an empty line
		if len(txt) == 0 {
			s.evt = Event{
				Type: event,
				Data: data.Bytes(),
			}
			return true
		}

		// Split a string like "event: bar" into name="event" and value=" bar".
		name, value, _ := bytes.Cut(txt, []byte(":"))

		// Consume an optional space after the colon if it exists.
		if len(value) > 0 && value[0] == ' ' {
			value = value[1:]
		}

		switch string(name) {
		case "":
			// SSE comment lines (starting with ":") are intentionally ignored.
		case "event":
			event = string(value)
		case "data":
			_, s.err = data.Write(value)
			if s.err != nil {
				return false
			}
			_, s.err = data.WriteRune('\n')
			if s.err != nil {
				return false
			}
		}
	}

	if s.scn.Err() != nil {
		s.err = s.scn.Err()
	}

	return false
}

func (s *eventStreamDecoder) Event() Event {
	return s.evt
}

func (s *eventStreamDecoder) Close() error {
	return s.rc.Close()
}

func (s *eventStreamDecoder) Err() error {
	return s.err
}

type Stream[T any] struct {
	decoder Decoder
	cur     T
	err     error
}

func NewStream[T any](decoder Decoder, err error) *Stream[T] {
	return &Stream[T]{
		decoder: decoder,
		err:     err,
	}
}

// Next returns false if the stream has ended or an error occurred.
// Call Stream.Current() to get the current value.
// Call Stream.Err() to get the error.
//
//		for stream.Next() {
//			data := stream.Current()
//		}
//
//	 	if stream.Err() != nil {
//			...
//	 	}
func (s *Stream[T]) Next() bool {
	if s.err != nil {
		return false
	}

	if s.decoder == nil {
		s.err = ErrNilDecoder
		return false
	}

	for s.decoder.Next() {
		var nxt T
		s.err = json.Unmarshal(s.decoder.Event().Data, &nxt)
		if s.err != nil {
			return false
		}
		s.cur = nxt
		return true //nolint:staticcheck // SA4004: intentional iterator pattern - process one event per call
	}

	// decoder.Next() may be false because of an error
	s.err = s.decoder.Err()

	return false
}

func (s *Stream[T]) Current() T {
	return s.cur
}

func (s *Stream[T]) Err() error {
	return s.err
}

func (s *Stream[T]) Close() error {
	if s.decoder == nil {
		// already closed
		return nil
	}
	return s.decoder.Close()
}

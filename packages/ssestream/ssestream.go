package ssestream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"sync"
)

const (
	// maxSSEDataSize is the maximum accumulated data payload per SSE event.
	maxSSEDataSize = 32 * 1024 * 1024
	// maxSSELineSize bounds any single SSE line to avoid large transient
	// allocations from oversized data lines.
	maxSSELineSize = 256 * 1024
	// defaultSSEReaderSize keeps internal read buffering small and stable.
	defaultSSEReaderSize = 4 * 1024
)

// ErrNilDecoder is returned by Stream.Next when the SSE decoder is nil,
// indicating the response body was not available.
var ErrNilDecoder = errors.New("ssestream: decoder is nil")

// nilDecoder is returned by NewDecoder when the response or body is nil.
// It always reports ErrNilDecoder, making the error discoverable without Stream.
type nilDecoder struct{}

func (d *nilDecoder) Event() Event { return Event{} }
func (d *nilDecoder) Next() bool   { return false }
func (d *nilDecoder) Close() error { return nil }
func (d *nilDecoder) Err() error   { return ErrNilDecoder }

type Decoder interface {
	Event() Event
	Next() bool
	Close() error
	Err() error
}

func NewDecoder(res *http.Response) Decoder {
	if res == nil || res.Body == nil {
		return &nilDecoder{}
	}

	var decoder Decoder
	mediaType, _, _ := mime.ParseMediaType(res.Header.Get("content-type"))
	decoderTypesMu.RLock()
	t, ok := decoderTypes[mediaType]
	decoderTypesMu.RUnlock()
	if ok {
		decoder = t(res.Body)
	}
	if decoder == nil {
		reader := bufio.NewReaderSize(res.Body, defaultSSEReaderSize)
		decoder = &eventStreamDecoder{rc: res.Body, reader: reader}
	}
	return decoder
}

var (
	decoderTypes   = map[string](func(io.ReadCloser) Decoder){}
	decoderTypesMu sync.RWMutex
)

func RegisterDecoder(contentType string, decoder func(io.ReadCloser) Decoder) {
	if decoder == nil {
		panic("ssestream: RegisterDecoder decoder cannot be nil")
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType == "" {
		panic(fmt.Sprintf("ssestream: RegisterDecoder invalid content type %q", contentType))
	}
	mediaType = strings.ToLower(mediaType)

	decoderTypesMu.Lock()
	decoderTypes[mediaType] = decoder
	decoderTypesMu.Unlock()
}

type Event struct {
	Type string
	Data []byte
}

// A base implementation of a Decoder for text/event-stream.
type eventStreamDecoder struct {
	evt          Event
	rc           io.ReadCloser
	reader       *bufio.Reader
	err          error
	maxDataBytes int // max accumulated data size per event; 0 uses maxSSEDataSize
	maxLineBytes int // max bytes in a single SSE line; 0 uses maxSSELineSize
}

func (s *eventStreamDecoder) readLine(dst *bytes.Buffer) ([]byte, bool, error) {
	if s.reader == nil {
		return nil, false, ErrNilDecoder
	}

	lineLimit := s.maxLineBytes
	if lineLimit == 0 {
		lineLimit = maxSSELineSize
	}

	dst.Reset()

	for {
		fragment, readErr := s.reader.ReadSlice('\n')
		if len(fragment) > 0 {
			if dst.Len()+len(fragment) > lineLimit {
				return nil, false, fmt.Errorf("event line exceeds %d bytes", lineLimit)
			}
			if _, err := dst.Write(fragment); err != nil {
				return nil, false, err
			}
		}

		switch {
		case readErr == nil:
			line := dst.Bytes()
			if len(line) > 0 && line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
			}
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			return line, false, nil
		case errors.Is(readErr, bufio.ErrBufferFull):
			continue
		case errors.Is(readErr, io.EOF):
			if dst.Len() == 0 {
				return nil, true, nil
			}
			line := dst.Bytes()
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			return line, false, nil
		default:
			return nil, false, readErr
		}
	}
}

func (s *eventStreamDecoder) Next() bool {
	if s.err != nil {
		return false
	}

	event := ""
	data := bytes.NewBuffer(nil)
	hasData := false

	var lineBuf bytes.Buffer
	for {
		txt, eof, err := s.readLine(&lineBuf)
		if err != nil {
			s.err = err
			return false
		}
		if eof {
			break
		}

		// Per SSE spec §9.2.6: on a blank line, dispatch only if data
		// was received. Otherwise reset the event type and continue.
		if len(txt) == 0 {
			if !hasData {
				event = ""
				continue
			}
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
			hasData = true
			if data.Len() > 0 {
				_, s.err = data.WriteRune('\n')
				if s.err != nil {
					return false
				}
			}
			_, s.err = data.Write(value)
			if s.err != nil {
				return false
			}
			limit := s.maxDataBytes
			if limit == 0 {
				limit = maxSSEDataSize
			}
			if data.Len() > limit {
				s.err = fmt.Errorf("event data exceeds %d bytes", limit)
				return false
			}
		}
	}

	// Per the SSE spec, dispatch any buffered event when the connection
	// closes without a trailing blank line. Use hasData (not data.Len)
	// because "data:" with an empty value is a valid event with empty data.
	if hasData {
		s.evt = Event{
			Type: event,
			Data: data.Bytes(),
		}
		return true
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
		data := s.decoder.Event().Data
		if len(data) == 0 {
			continue
		}
		var nxt T
		s.err = json.Unmarshal(data, &nxt)
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
		return nil
	}
	err := s.decoder.Close()
	s.decoder = nil
	return err
}

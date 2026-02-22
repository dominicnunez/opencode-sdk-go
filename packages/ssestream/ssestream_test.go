package ssestream

import (
	"io"
	"sync"
	"testing"
)

type mockDecoder struct{}

func (m *mockDecoder) Event() Event        { return Event{} }
func (m *mockDecoder) Next() bool          { return false }
func (m *mockDecoder) Close() error        { return nil }
func (m *mockDecoder) Err() error          { return nil }

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

package stream

import (
	"fmt"
)

func NewWriter[T any](factory func(body []byte) T, stream WriteStream[T]) *Writer[T] {
	return &Writer[T]{
		factory: factory,
		stream:  stream,
	}
}

type WriteStream[T any] interface {
	Send(sendRequest T) error
}

type Writer[T any] struct {
	factory func(body []byte) T
	stream  WriteStream[T]
}

func (w *Writer[T]) Write(p []byte) (n int, err error) {
	for i := 0; i < len(p); i += 1024 {
		start, end := i, i+1024
		if len(p) <= end {
			end = len(p)
		}
		if err := w.stream.Send(w.factory(p[start:end])); err != nil {
			return 0, fmt.Errorf("w.stream.Send: %w", err)
		}
	}
	return len(p), nil
}

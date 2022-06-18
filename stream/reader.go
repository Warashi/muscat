package stream

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

func NewReader[T ReadResponse](stream ReadStream[T]) *Reader[T] {
	return &Reader[T]{
		stream: stream,
	}
}

type ReadResponse interface {
	GetBody() []byte
}

type ReadStream[T ReadResponse] interface {
	Recv() (T, error)
}

type Reader[T ReadResponse] struct {
	stream ReadStream[T]
	buffer bytes.Buffer
}

func (r *Reader[T]) Read(p []byte) (n int, err error) {
	for r.buffer.Len() < len(p) {
		b, err := r.stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, fmt.Errorf("r.stream.Recv: %w", err)
		}
		if _, err := r.buffer.Write(b.GetBody()); err != nil {
			return 0, fmt.Errorf("r.buffer.Write: %w", err)
		}
	}
	return r.buffer.Read(p)
}

package stream

import (
	"errors"
	"io"
)

func NewBidiReader[T ReadResponse](stream BidiReadStream[T]) *BidiReader[T] {
	return &BidiReader[T]{
		stream: stream,
	}
}

type BidiReadStream[T ReadResponse] interface {
	Receive() (T, error)
}

type BidiReader[T ReadResponse] struct {
	stream BidiReadStream[T]
	buf    []byte
}

func (r *BidiReader[T]) Read(p []byte) (n int, err error) {
	end := false
	for len(r.buf) < len(p) {
		msg, err := r.stream.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				end = true
				break
			}
			return 0, err
		}
		r.buf = append(r.buf, msg.GetBody()...)
	}
	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	if end && n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

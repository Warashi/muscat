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
	closed bool
}

func (r *BidiReader[T]) Read(p []byte) (n int, err error) {
	if r.closed {
		n := copy(p, r.buf)
		r.buf = r.buf[n:]
		if len(r.buf) == 0 {
			return n, io.EOF
		}
		return n, nil
	}
	if len(r.buf) < len(p) {
		msg, err := r.stream.Receive()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return 0, err
			}
			// errors.Is(err, io.EOF)
			r.closed = true
		}
		r.buf = append(r.buf, msg.GetBody()...)
	}
	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

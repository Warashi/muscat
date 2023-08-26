package stream

import "io"

func NewReader[T ReadResponse](stream ReadStream[T]) *Reader[T] {
	return &Reader[T]{
		stream: stream,
	}
}

type ReadResponse interface {
	GetBody() []byte
}

type ReadStream[T ReadResponse] interface {
	Receive() bool
	Msg() T
}

type Reader[T ReadResponse] struct {
	stream ReadStream[T]
	buf    []byte
	closed bool
}

func (r *Reader[T]) Read(p []byte) (n int, err error) {
	if r.closed {
		n := copy(p, r.buf)
		r.buf = r.buf[n:]
		if len(r.buf) == 0 {
			return n, io.EOF
		}
		return n, nil
	}
	if len(r.buf) < len(p) {
		if !r.stream.Receive() {
			r.closed = true
		}
		r.buf = append(r.buf, r.stream.Msg().GetBody()...)
	}
	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

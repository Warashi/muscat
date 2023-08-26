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
}

func (r *Reader[T]) Read(p []byte) (n int, err error) {
	end := false
	for len(r.buf) < len(p) {
		if !r.stream.Receive() {
			end = true
			break
		}
		r.buf = append(r.buf, r.stream.Msg().GetBody()...)
	}
	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	if end && n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

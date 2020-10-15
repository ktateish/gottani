package lib

import "bytes"

func NewT() *T {
	return &T{new(bytes.Buffer)}
}

func (t *T) Write(b []byte) (int, error) {
	return t.buf.Write(b)
}

func (t *T) String() string {
	return t.buf.String()
}

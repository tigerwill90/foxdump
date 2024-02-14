package foxdump

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

type shortResponseWriter struct{}

func (m *shortResponseWriter) Status() int {
	return http.StatusOK
}

func (m *shortResponseWriter) Written() bool {
	return true
}

func (m *shortResponseWriter) Size() int {
	return 0
}

func (m *shortResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *shortResponseWriter) Write(p []byte) (n int, err error) {
	return len(p) - 5, nil
}

func (m *shortResponseWriter) WriteString(s string) (n int, err error) {
	return len(s) - 5, nil
}

func (m *shortResponseWriter) WriteHeader(int) {}

func TestMultiWriter_RwErrShortWrite(t *testing.T) {
	mw := multiWriter{new(shortResponseWriter), bytes.NewBuffer(nil)}
	_, err := mw.Write([]byte("foobar"))
	assert.ErrorIs(t, err, io.ErrShortWrite)
	_, err = mw.WriteString("foobar")
	assert.ErrorIs(t, err, io.ErrShortWrite)
}

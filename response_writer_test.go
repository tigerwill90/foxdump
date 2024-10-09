package foxdump

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/tigerwill90/fox"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

var _ fox.ResponseWriter = (*shortResponseWriter)(nil)

type shortResponseWriter struct{}

func (m shortResponseWriter) Status() int {
	return http.StatusOK
}

func (m shortResponseWriter) Written() bool {
	return true
}

func (m shortResponseWriter) Size() int {
	return 0
}

func (m shortResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m shortResponseWriter) Write(p []byte) (n int, err error) {
	return len(p) - 5, nil
}

func (m shortResponseWriter) WriteString(s string) (n int, err error) {
	return len(s) - 5, nil
}

func (m shortResponseWriter) WriteHeader(int) {}

func (m shortResponseWriter) ReadFrom(r io.Reader) (n int64, err error) { return 0, nil }

func (m shortResponseWriter) FlushError() error { return nil }

func (m shortResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) { return nil, nil, nil }

func (m shortResponseWriter) Push(target string, opts *http.PushOptions) error { return nil }

func (m shortResponseWriter) SetReadDeadline(deadline time.Time) error {
	return nil
}

func (m shortResponseWriter) SetWriteDeadline(deadline time.Time) error {
	return nil
}

func TestMultiWriter_RwErrShortWrite(t *testing.T) {
	mw := multiWriter{new(shortResponseWriter), bytes.NewBuffer(nil)}
	_, err := mw.Write([]byte("foobar"))
	assert.ErrorIs(t, err, io.ErrShortWrite)
	_, err = mw.WriteString("foobar")
	assert.ErrorIs(t, err, io.ErrShortWrite)
}

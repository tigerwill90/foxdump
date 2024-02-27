package foxdump

import (
	"bytes"
	"github.com/tigerwill90/fox"
	"io"
	"net/http"
	"sync"
)

var copyBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 32*1024)
		return &b
	},
}

type onlyWrite struct {
	io.Writer
}

type multiWriter struct {
	fox.ResponseWriter
	buf *bytes.Buffer
}

func (w multiWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w multiWriter) Write(p []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(p)
	if err != nil {
		return
	}
	if n != len(p) {
		err = io.ErrShortWrite
		return
	}

	n, err = w.buf.Write(p)
	if err != nil {
		return
	}

	return n, nil
}

func (w multiWriter) WriteString(s string) (n int, err error) {
	n, err = w.ResponseWriter.WriteString(s)
	if err != nil {
		return
	}
	if n != len(s) {
		err = io.ErrShortWrite
		return
	}

	n, err = w.buf.WriteString(s)
	if err != nil {
		return
	}

	return n, nil
}

func (w multiWriter) ReadFrom(src io.Reader) (n int64, err error) {
	bufp := copyBufPool.Get().(*[]byte)
	buf := *bufp
	// onlyWrite hide "ReadFrom" from w.
	n, err = io.CopyBuffer(onlyWrite{w}, src, buf)
	copyBufPool.Put(bufp)
	return
}

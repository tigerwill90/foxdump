package foxdump

import (
	"crypto/rand"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tigerwill90/fox"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

type repeatReader struct {
	data []byte
	eof  bool
}

func (r *repeatReader) Read(p []byte) (n int, err error) {
	if !r.eof {
		n = copy(p, r.data)
		r.eof = true
		return
	}
	r.eof = false
	return 0, io.EOF
}

func (r *repeatReader) Close() error { return nil }

func BenchmarkMiddleware(b *testing.B) {
	f := fox.New(fox.WithMiddleware(Middleware(func(c fox.Context, buf []byte) {

	}, func(c fox.Context, buf []byte) {

	})))

	buf := make([]byte, 1*1024)
	rand.Read(buf)

	f.MustHandle(http.MethodPost, "/foo/bar", func(c fox.Context) {
		_, _ = c.Writer().Write(buf)
	})

	body := &repeatReader{data: buf}
	req := httptest.NewRequest(http.MethodPost, "/foo/bar", nil)
	w := new(mockResponseWriter)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req.Body = body
		f.ServeHTTP(w, req)
	}
}

func BenchmarkEchoDumpMiddleware(b *testing.B) {
	e := echo.New()
	e.Use(middleware.BodyDump(func(context echo.Context, bytes []byte, bytes2 []byte) {

	}))

	buf := make([]byte, 1*1024)
	rand.Read(buf)

	e.POST("/foo/bar", func(c echo.Context) error {
		_, _ = c.Response().Write(buf)
		return nil
	})

	body := &repeatReader{data: buf}
	req := httptest.NewRequest(http.MethodPost, "/foo/bar", nil)
	w := new(mockResponseWriter)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req.Body = body
		e.ServeHTTP(w, req)
	}
}

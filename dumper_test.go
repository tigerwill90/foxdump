package dumpfox

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tigerwill90/fox"
	"net/http"
	"net/http/httptest"
	"strings"
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

func BenchmarkMiddleware(b *testing.B) {
	f := fox.New(fox.WithMiddleware(Middleware(func(c fox.Context, req, res []byte) {

	})))

	buf := []byte("foo bar")

	f.MustHandle(http.MethodPost, "/foo/bar", func(c fox.Context) {
		_, _ = c.Writer().Write(buf)
	})

	req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"foo":"bar"}`))
	w := new(mockResponseWriter)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f.ServeHTTP(w, req)
	}
}

func BenchmarkEchoDumpMiddleware(b *testing.B) {
	e := echo.New()
	e.Use(middleware.BodyDump(func(context echo.Context, bytes []byte, bytes2 []byte) {

	}))

	buf := []byte("foo bar")

	e.POST("/foo/bar", func(c echo.Context) error {
		_, _ = c.Response().Write(buf)
		return nil
	})

	req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"foo":"bar"}`))
	w := new(mockResponseWriter)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e.ServeHTTP(w, req)
	}
}

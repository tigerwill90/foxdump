package foxdump

import (
	"bytes"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func BenchmarkFoxDumpMiddleware(b *testing.B) {
	f := fox.New(fox.WithMiddleware(Middleware(func(c fox.Context, buf []byte) {

	}, func(c fox.Context, buf []byte) {

	})))

	buf := make([]byte, 1*1024)
	_, _ = rand.Read(buf)

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

func TestBodyDumper_DumpBody(t *testing.T) {
	var nilBodyHandler = func(t *testing.T, want []byte) BodyHandler {
		return nil
	}

	cases := []struct {
		name string
		req  func(t *testing.T, want []byte) BodyHandler
		res  func(t *testing.T, want []byte) BodyHandler
	}{
		{
			name: "dump request only",
			req: func(t *testing.T, want []byte) BodyHandler {
				return func(c fox.Context, buf []byte) {
					assert.Equal(t, want, buf)
				}
			},
			res: nilBodyHandler,
		},
		{
			name: "dump response only",
			req:  nilBodyHandler,
			res: func(t *testing.T, want []byte) BodyHandler {
				return func(c fox.Context, buf []byte) {
					assert.Equal(t, want, buf)
				}
			},
		},
		{
			name: "dump request and response",
			req: func(t *testing.T, want []byte) BodyHandler {
				return func(c fox.Context, buf []byte) {
					assert.Equal(t, want, buf)
				}
			},
			res: func(t *testing.T, want []byte) BodyHandler {
				return func(c fox.Context, buf []byte) {
					assert.Equal(t, want, buf)
				}
			},
		},
		{
			name: "both nil req and res",
			req:  nilBodyHandler,
			res:  nilBodyHandler,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			buf := make([]byte, 1*1024*1024)
			_, err := rand.Read(buf)
			require.NoError(t, err)

			f := fox.New(fox.WithMiddleware(Middleware(tc.req(t, buf), tc.res(t, buf))))
			require.NoError(t, f.Handle(http.MethodPost, "/foo", func(c fox.Context) {
				assert.NoError(t, c.Blob(http.StatusOK, fox.MIMEOctetStream, buf))
			}))

			req := httptest.NewRequest(http.MethodPost, "/foo", bytes.NewReader(buf))
			w := httptest.NewRecorder()
			f.ServeHTTP(w, req)
			assert.Equal(t, buf, w.Body.Bytes())
		})
	}
}

func TestWithFilter(t *testing.T) {
	var failBodyHandler = func(t *testing.T, want []byte) BodyHandler {
		return func(c fox.Context, buf []byte) {
			t.Error("should not be call")
		}
	}

	cases := []struct {
		name   string
		filter Filter
		req    func(t *testing.T, want []byte) BodyHandler
		res    func(t *testing.T, want []byte) BodyHandler
	}{
		{
			name: "filter match req",
			filter: func(r *http.Request) bool {
				return r.URL.Path != "/foo"
			},
			res: failBodyHandler,
			req: failBodyHandler,
		},
		{
			name: "filter does not match req",
			filter: func(r *http.Request) bool {
				return r.URL.Path != "/bar"
			},
			res: func(t *testing.T, want []byte) BodyHandler {
				return func(c fox.Context, buf []byte) {
					assert.Equal(t, want, buf)
				}
			},
			req: func(t *testing.T, want []byte) BodyHandler {
				return func(c fox.Context, buf []byte) {
					assert.Equal(t, want, buf)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			buf := make([]byte, 1*1024*1024)
			_, err := rand.Read(buf)
			require.NoError(t, err)

			f := fox.New(fox.WithMiddleware(Middleware(tc.req(t, buf), tc.res(t, buf), WithFilter(tc.filter))))
			require.NoError(t, f.Handle(http.MethodPost, "/foo", func(c fox.Context) {
				assert.NoError(t, c.Blob(http.StatusOK, fox.MIMEOctetStream, buf))
			}))

			req := httptest.NewRequest(http.MethodPost, "/foo", bytes.NewReader(buf))
			w := httptest.NewRecorder()
			f.ServeHTTP(w, req)
			assert.Equal(t, buf, w.Body.Bytes())
		})
	}
}

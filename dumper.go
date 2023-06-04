package dumpfox

import (
	"bytes"
	"github.com/tigerwill90/fox"
	"io"
	"net/http"
	"sync"
)

var p = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(nil)
	},
}

// BodyDumpHandler is a callback that receives the HTTP request and response bodies.
// The 'req' and 'res' byte slices are transient, and their data is only guaranteed to be valid
// during the execution of the BodyDumpHandler function. If the data needs to be persisted or
// used outside the scope of this function, it should be copied to a new byte slice (e.g., using 'copy').
// Furthermore, these slices should be treated as read-only to prevent any unintended side effects.
type BodyDumpHandler func(c fox.Context, req, res []byte)

// Middleware is a convenience function that creates a new BodyDumper middleware instance 0with the
// given BodyDumpHandler and returns the DumpBody function. Options can be provided to configure the tracer.
func Middleware(fn BodyDumpHandler, opts ...Option) fox.MiddlewareFunc {
	return NewBodyDumper(fn, opts...).DumpBody
}

// BodyDumper is a middleware that dumps the HTTP request and response bodies.
// It calls a BodyDumpHandler function with the body content.
type BodyDumper struct {
	fn  BodyDumpHandler
	cfg *config
}

// NewBodyDumper creates a new BodyDumper with the given BodyDumpHandler function and options.
func NewBodyDumper(fn BodyDumpHandler, opts ...Option) *BodyDumper {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return &BodyDumper{
		fn:  fn,
		cfg: cfg,
	}
}

// DumpBody is the middleware function that gets called for each HTTP request. It reads the request and response
// bodies, if configured to do so, and then calls the BodyDumpHandler function with the body content.
// See the BodyDumpHandler documentation for guidelines on the correct usage of the body content.
func (d *BodyDumper) DumpBody(next fox.HandlerFunc) fox.HandlerFunc {
	return func(c fox.Context) {
		if !d.cfg.res && !d.cfg.req {
			next(c)
			return
		}

		for _, f := range d.cfg.filters {
			if !f(c.Request()) {
				next(c)
				return
			}
		}

		buf := p.Get().(*bytes.Buffer)
		buf.Reset()
		defer p.Put(buf)

		var offset int64
		if d.cfg.req {
			var err error
			offset, err = buf.ReadFrom(c.Request().Body)
			if err != nil {
				http.Error(c.Writer(), http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			cpBuf := p.Get().(*bytes.Buffer)
			cpBuf.Reset()
			defer p.Put(cpBuf)

			// Safe as Buffer.Writer make a copy of p
			cpBuf.Write(buf.Bytes())

			c.Request().Body = nopCloser{cpBuf}
		}

		if d.cfg.res {
			c.TeeWriter(buf)
		}

		next(c)

		d.fn(c, buf.Bytes()[:offset], buf.Bytes()[offset:])
	}
}

type nopCloser struct {
	*bytes.Buffer
}

func (nopCloser) Close() error { return nil }

func (c nopCloser) WriteTo(w io.Writer) (n int64, err error) {
	return c.Buffer.WriteTo(w)
}

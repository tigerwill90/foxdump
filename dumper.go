package foxdump

import (
	"bytes"
	"github.com/tigerwill90/fox"
	"io"
	"log"
	"sync"
)

var p = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(nil)
	},
}

// BodyHandler is a callback that receives the HTTP request or response body. The buf slice is transient, and its
// data is only guaranteed to be valid during the execution of the BodyHandler function. If the data needs to be
// persisted or used outside the scope of this function, it should be copied to a new byte slice (e.g., using 'copy').
// Furthermore, 'buf' should be treated as read-only slice to prevent any unintended side effects.
type BodyHandler func(c fox.Context, buf []byte)

// Middleware is a convenience function that creates a new BodyDumper middleware instance with the
// given BodyHandler functions and returns the DumpBody function. Options can be provided to configure the dumper.
func Middleware(req BodyHandler, res BodyHandler, opts ...Option) fox.MiddlewareFunc {
	return NewBodyDumper(req, res, opts...).DumpBody
}

// BodyDumper is a middleware that dumps the HTTP request and response bodies.
// It calls a BodyHandler function with the body content.
type BodyDumper struct {
	req BodyHandler
	res BodyHandler
	cfg *config
}

// NewBodyDumper creates a new BodyDumper with the given BodyHandler functions and options.
func NewBodyDumper(req BodyHandler, res BodyHandler, opts ...Option) *BodyDumper {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return &BodyDumper{
		req: req,
		res: res,
		cfg: cfg,
	}
}

// DumpBody is the middleware function that gets called for each HTTP request. It reads the request and response
// bodies, if configured to do so, and then calls the BodyHandler function with the body content.
// See the BodyHandler documentation for guidelines on the correct usage of the body content.
func (d *BodyDumper) DumpBody(next fox.HandlerFunc) fox.HandlerFunc {
	return func(c fox.Context) {
		if d.req == nil && d.res == nil {
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

		if d.req != nil {
			_, err := buf.ReadFrom(c.Request().Body)
			if err != nil {
				log.Println("body dumper: unexpected error while reading request body")
				buf.Reset()
				goto RespFallback
			}

			cpBuf := p.Get().(*bytes.Buffer)
			cpBuf.Reset()
			defer p.Put(cpBuf)

			// Safe as Buffer.Writer make a copy of p
			cpBuf.Write(buf.Bytes())

			d.req(c, buf.Bytes())
			buf.Reset()

			c.Request().Body = nopCloser{cpBuf}
		}

	RespFallback:
		if d.res != nil {
			c.TeeWriter(buf)
			next(c)
			d.res(c, buf.Bytes())
			return
		}

		next(c)
	}
}

type nopCloser struct {
	*bytes.Buffer
}

func (nopCloser) Close() error { return nil }

func (c nopCloser) WriteTo(w io.Writer) (n int64, err error) {
	return c.Buffer.WriteTo(w)
}

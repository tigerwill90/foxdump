package dumpfox

import "net/http"

type config struct {
	req     bool
	res     bool
	filters []Filter
}

func defaultConfig() *config {
	return &config{
		req: true,
		res: true,
	}
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) {
	f(c)
}

type Filter func(r *http.Request) bool

// DisableRequestDump disables the dumping of HTTP requests. When this option is set, the BodyDumpHandler
// will receive an empty 'req' slice.
func DisableRequestDump() Option {
	return optionFunc(func(c *config) {
		c.req = false
	})
}

// DisableResponseDump disables the dumping of HTTP responses. When this option is set, the BodyDumpHandler
// will receive an empty 'res' slice.
func DisableResponseDump() Option {
	return optionFunc(func(c *config) {
		c.res = false
	})
}

// WithFilter appends the provided filters to the middleware's filter list.
// A filter returning false will exclude the request from being traced. If no filters
// are provided, all requests will be traced. Keep in mind that filters are invoked for each request,
// so they should be simple and efficient.
func WithFilter(f ...Filter) Option {
	return optionFunc(func(c *config) {
		c.filters = f
	})
}

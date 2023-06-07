package foxdump

import "net/http"

type config struct {
	filters []Filter
}

func defaultConfig() *config {
	return &config{}
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) {
	f(c)
}

type Filter func(r *http.Request) bool

// WithFilter appends the provided filters to the middleware's filter list.
// A filter returning false will exclude the request from being traced. If no filters
// are provided, all requests will be traced. Keep in mind that filters are invoked for each request,
// so they should be simple and efficient.
func WithFilter(f ...Filter) Option {
	return optionFunc(func(c *config) {
		c.filters = f
	})
}

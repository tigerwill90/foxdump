// Copyright 2023 Sylvain Müller. All rights reserved.
// Mount of this source code is governed by a MIT license that can be found
// at https://github.com/tigerwill90/foxdump/blob/master/LICENSE.txt.

package foxdump

import (
	"github.com/tigerwill90/fox"
)

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

type Filter func(c fox.Context) (skip bool)

// WithFilter appends the provided filters to the middleware's filter list.
// A filter returning true will exclude the request from being dumped. If no filters
// are provided, all requests will be dumped. Keep in mind that filters are invoked for each request,
// so they should be simple and efficient.
func WithFilter(f ...Filter) Option {
	return optionFunc(func(c *config) {
		c.filters = f
	})
}

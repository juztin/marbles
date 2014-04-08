// Copyright 2014 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

type static struct {
	path    string
	handler HandlerFunc
}

func (r *static) Path() string {
	return r.path
}

func (r *static) IsCanonical(url string) (string, bool) {
	return IsCanonical(url)
}

func (r *static) Matches(path string) bool {
	return r.path == path
}

func (r *static) Execute(ctx Context) {
	r.handler(ctx)
}

func Static(expr string) func(h HandlerFunc) Route {
	return func(h HandlerFunc) Route {
		return NewStatic(expr, h)
	}
}

func NewStatic(path string, h HandlerFunc) Route {
	rt := new(static)
	rt.path = path
	rt.handler = h

	return rt
}

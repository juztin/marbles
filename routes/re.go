// Copyright 2014 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

import "regexp"

type re struct {
	path    string
	params  func() map[string]string
	expr    *regexp.Regexp
	handler Handler
}

func (r *re) Path() string {
	return r.path
}

func (r *re) IsCanonical(url string) (string, bool) {
	return IsCanonical(url)
}

func (r *re) Matches(url string) bool {
	return r.expr.MatchString(url)
}

func (r *re) Execute(ctx Context) {
	ctx.Params = params(r.expr, ctx.URL.Path)
	r.handler(ctx)
}

func params(expr *regexp.Regexp, url string) func() map[string]string {
	return func() map[string]string {
		data := make(map[string]string)
		//matches := r.expr.FindAllStringSubmatch(url, -1)
		matches := expr.FindStringSubmatch(url)

		for i, n := range expr.SubexpNames() {
			if i == 0 {
				continue
			}
			data[n] = matches[i]
		}

		return data
	}
}

func Re(expr string) func(h Handler) Route {
	return func(h Handler) Route {
		return NewRe(expr, h)
	}
}

func NewRe(expr string, h Handler) Route {
	rt := new(re)
	rt.path = expr
	rt.expr = regexp.MustCompile(expr)
	rt.handler = h

	return rt
}

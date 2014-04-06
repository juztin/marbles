// Copyright 2014 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

type wildcard struct {
	handler Handler
}

func (r *wildcard) Path() string {
	return ""
}

func (r *wildcard) IsCanonical(url string) (string, bool) {
	return url, true
}

func (r *wildcard) Matches(path string) bool {
	return true
}

func (r *wildcard) Execute(ctx Context) {
	r.handler(ctx)
}

func Wildcard() func(h Handler) Route {
	return func(h Handler) Route {
		return NewWildcard(h)
	}
}

func NewWildcard(h Handler) Route {
	rt := new(wildcard)
	rt.handler = h

	return rt
}

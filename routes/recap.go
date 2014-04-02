// Copyright 2014 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

import (
	"reflect"
	"regexp"
)

type reCap struct {
	path    string
	expr    *regexp.Regexp
	handler reflect.Value
}

func (r *reCap) Path() string {
	return r.path
}

func (r reCap) IsCanonical(url string) (string, bool) {
	return IsCanonical(url)
}

func (r *reCap) Matches(url string) bool {
	return r.expr.MatchString(url)
}

func (r *reCap) Execute(ctx Context) {
	// TODO it would be nice if we could detect numbers and cast them as such prior to invoking the func
	args := []reflect.Value{reflect.ValueOf(ctx)}
	matches := r.expr.FindStringSubmatch(ctx.URL.Path)
	params := make(map[string]string)
	for i, a := range matches[1:] {
		args = append(args, reflect.ValueOf(a))
		params[r.expr.SubexpNames()[i+1]] = a
	}
	ctx.Params = func() map[string]string {
		return params
	}
	r.handler.Call(args)
}

func ReCap(expr string) func(h interface{}) Route {
	return func(h interface{}) Route {
		return NewReCap(expr, h)
	}
}

func NewReCap(expr string, h interface{}) Route {
	r := new(reCap)
	r.path = expr
	r.expr = regexp.MustCompile(expr)

	if fn, ok := h.(reflect.Value); ok {
		r.handler = fn
	} else {
		r.handler = reflect.ValueOf(h)
	}

	return r
}

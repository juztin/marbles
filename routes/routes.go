// Copyright 2014 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(ctx Context)

type Params func() map[string]string

type Context struct {
	*http.Request
	Response http.ResponseWriter
	Params   Params
}

type Route interface {
	Path() string
	IsCanonical(path string) (string, bool)
	Matches(path string) bool
	Execute(ctx Context)
}

type Routes struct {
	col map[string][]Route
}

func (c *Context) Redirect(path string) {
	http.Redirect(c.Response, c.Request, path, http.StatusFound)
}

func (c *Context) RedirectPerm(path string) {
	r := c.Response
	r.Header().Set("Location", path)
	r.WriteHeader(http.StatusMovedPermanently)
}

func (c *Context) Error(status int) {
	c.Response.WriteHeader(status)
}

func (s *Routes) Add(r Route, methods ...string) error {
	for _, m := range methods {
		c, ok := s.col[strings.ToUpper(m)]
		if !ok {
			return errors.New(fmt.Sprintf("Invalid HTTP method: %s", m))
		}
		s.col[m] = append(c, r)
	}
	return nil
}
func (s *Routes) Options(r Route) *Routes {
	s.Add(r, "OPTIONS")
	return s
}
func (s *Routes) Head(r Route) *Routes {
	s.Add(r, "HEAD")
	return s
}
func (s *Routes) Get(r Route) *Routes {
	s.Add(r, "GET")
	return s
}
func (s *Routes) Post(r Route) *Routes {
	s.Add(r, "POST")
	return s
}
func (s *Routes) Put(r Route) *Routes {
	s.Add(r, "PUT")
	return s
}
func (s *Routes) Delete(r Route) *Routes {
	s.Add(r, "DELETE")
	return s
}
func (s *Routes) Trace(r Route) *Routes {
	s.Add(r, "TRACE")
	return s
}
func (s *Routes) Connect(r Route) *Routes {
	s.Add(r, "CONNECT")
	return s
}

func (s *Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	c, ok := s.col[r.Method]
	if !ok {
		ctx.Error(http.StatusNotFound)
		return
	}

	for _, rt := range c {
		if !rt.Matches(r.URL.Path) {
			continue
		}

		if path, ok := rt.IsCanonical(r.URL.Path); !ok {
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}
			ctx.RedirectPerm(path)
		} else {
			//defer _500Handler(ctx)
			rt.Execute(ctx)
		}
		return
	}

	ctx.Error(http.StatusNotFound)
}

func emptyParams() map[string]string {
	return make(map[string]string)
}

func IsCanonical(p string) (string, bool) {
	if len(p) == 0 {
		return "/", false
	} else if p[0] != '/' {
		return "/" + p, false
	}

	cp := path.Clean(p)

	if cp[len(cp)-1] != '/' {
		cp = cp + "/"
		return cp, cp == p
	}

	return cp, cp == p
}

func Wrap(fn http.HandlerFunc) HandlerFunc {
	return func(ctx Context) {
		fn(ctx.Response, ctx.Request)
	}
}

func New() *Routes {
	r := new(Routes)
	r.col = make(map[string][]Route)
	r.col["ADD"] = []Route{}
	r.col["OPTIONS"] = []Route{}
	r.col["HEAD"] = []Route{}
	r.col["GET"] = []Route{}
	r.col["POST"] = []Route{}
	r.col["PUT"] = []Route{}
	r.col["DELETE"] = []Route{}
	r.col["TRACE"] = []Route{}
	r.col["CONNECT"] = []Route{}
	return r
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return Context{r, w, emptyParams}
}

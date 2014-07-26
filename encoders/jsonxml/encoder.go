// Copyright 2014 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package jsonxml implements some simple boilerplate for
// writing JSON/XML data.
package jsonxml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
)

// Signature that both the encoding/xml and encoding/json match.
type marshalFunc func(v interface{}) ([]byte, error)
type unmarshalFunc func(data []byte, v interface{}) error

// A wrapper for encoding/json to add padding to the resulting JSON data.
func jsonpMarshal(callback string) marshalFunc {
	return func(o interface{}) ([]byte, error) {
		data, err := json.Marshal(o)
		if err != nil {
			return nil, err
		}
		return pad(data, callback), nil
	}
}

// Adds padding to the given JSON data.
func pad(data []byte, callback string) []byte {
	l := len(callback)
	s := make([]byte, len(data)+l+2)
	copy(s[l+1:], data)
	copy(s[0:], []byte(callback))
	s[l] = '('
	s[len(s)-1] = ')'
	return s
}

// Marshal the given data, based on the requests Accept type, and write
// it to the response.
func Write(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return WriteStatus(w, r, data, http.StatusOK)
}

// Marshal the given data, based on the requests Accept type, and write
// it to the response with the status.
func WriteStatus(w http.ResponseWriter, r *http.Request, data interface{}, status int) error {
	var marshal marshalFunc
	ct := r.Header.Get("Accept")
	ct, _, _ = mime.ParseMediaType(ct)
	switch ct {
	default:
		callback := r.URL.Query().Get("callback")
		if len(callback) > 0 {
			ct = "application/javascript"
			marshal = jsonpMarshal(callback)
		} else {
			ct = "application/json"
			marshal = json.Marshal
		}
	case "application/json":
		marshal = json.Marshal
	case "application/xml":
		marshal = xml.Marshal
	}

	b, err := marshal(data)
	if err != nil {
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", ct)
	w.WriteHeader(status)
	if err != nil {
		return err
	}

	i, err := w.Write(b)
	if err != nil {
		err = fmt.Errorf("%d bytes written before error: %v", i, err)
	}
	return err
}

func read(r *http.Request, data interface{}, reader io.Reader) error {
	var unmarshal unmarshalFunc
	ct := r.Header.Get("Accept")
	ct, _, _ = mime.ParseMediaType(ct)
	switch ct {
	default:
		unmarshal = json.Unmarshal
	case "application/json":
		unmarshal = json.Unmarshal
	case "application/xml":
		unmarshal = xml.Unmarshal
	}

	if r.Method != "POST" && r.Method != "PUT" {
		return nil
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return unmarshal(body, data)
}

func Read(r *http.Request, data interface{}) error {
	return read(r, data, r.Body)
}

func ReadMax(r *http.Request, data interface{}, maxBytes int64) error {
	reader := io.LimitReader(r.Body, maxBytes)
	return read(r, data, reader)
}

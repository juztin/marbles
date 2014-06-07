package jsonxml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"mime"
	"net/http"
)

type encodeFunc func(o interface{}) ([]byte, error)

func jsonpEncoder(callback string) encodeFunc {
	return func(o interface{}) ([]byte, error) {
		data, err := json.Marshal(o)
		if err != nil {
			return nil, err
		}
		return pad(data, callback), nil
	}
}

func pad(data []byte, callback string) []byte {
	l := len(callback)
	s := make([]byte, len(data)+l+2)
	copy(s[l+1:], data)
	copy(s[0:], []byte(callback))
	s[l] = '('
	s[len(s)-1] = ')'
	return s
}

func Write(w http.ResponseWriter, r *http.Request, o interface{}) error {
	return WriteStatus(w, r, o, http.StatusOK)
}

func WriteStatus(w http.ResponseWriter, r *http.Request, o interface{}, status int) error {
	var encode encodeFunc
	ct := r.Header.Get("Accept")
	ct, _, _ = mime.ParseMediaType(ct)
	switch ct {
	default:
		callback := r.URL.Query().Get("callback")
		if len(callback) > 0 {
			ct = "application/javascript"
			encode = jsonpEncoder(callback)
		} else {
			ct = "application/json"
			encode = json.Marshal
		}
	case "application/json":
		encode = json.Marshal
	case "application/xml":
		encode = xml.Marshal
	}

	data, err := encode(o)
	if err != nil {
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", ct)
	w.WriteHeader(status)
	if err != nil {
		return err
	}

	i, err := w.Write(data)
	if err != nil {
		err = fmt.Errorf("%d bytes written before error: %v", i, err)
	}
	return err
}

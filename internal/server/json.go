package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type JSONWriter struct {
	w  http.ResponseWriter
	jw *json.Encoder
}

func newJSONWriter(w http.ResponseWriter) *JSONWriter {
	return &JSONWriter{
		w:  w,
		jw: json.NewEncoder(w),
	}
}

func (j *JSONWriter) Header() http.Header {
	return j.w.Header()
}

func (j *JSONWriter) Write(p []byte) (int, error) {

	var ans struct {
		Result string `json:"result"`
	}

	ans.Result = string(p)

	err := j.jw.Encode(ans)
	if err != nil {
		log.Println(err)
	}
	return len(p), err
}

func (j *JSONWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		j.w.Header().Set("Content-Type", "application/json")
	}
	j.w.WriteHeader(statusCode)
}

type JSONReader struct {
	r  io.ReadCloser
	jr *json.Decoder
}

func newJSONReader(r io.ReadCloser) *JSONReader {
	return &JSONReader{
		r:  r,
		jr: json.NewDecoder(r),
	}
}

func (j JSONReader) Read(p []byte) (int, error) {

	var req struct {
		URL string `json:"url"`
	}

	p2 := make([]byte, 0)

	for j.jr.More() {

		err := j.jr.Decode(&req)
		if err != nil {
			return 0, err
		}

		p2 = []byte(req.URL)
	}

	copy(p, p2)

	return len(p2), io.EOF
}

func (j *JSONReader) Close() error {
	if err := j.r.Close(); err != nil {
		return err
	}
	return nil
}

func JSONHandle(next http.Handler) http.Handler {
	jsonFn := func(w http.ResponseWriter, r *http.Request) {
		ow := w

		contentTypes := r.Header.Get("Content-Type")
		supportsJSON := strings.Contains(contentTypes, "application/json")

		if supportsJSON {
			jr := newJSONReader(r.Body)
			r.Body = jr
			defer jr.Close()
		}

		contentAccept := r.Header.Get("Accept")
		AcceptJSON := strings.Contains(contentAccept, "application/json")

		if AcceptJSON {
			jw := newJSONWriter(w)
			ow = jw
			ow.Header().Set("Content-Type", "application/json")
		}

		next.ServeHTTP(ow, r)
	}
	return http.HandlerFunc(jsonFn)
}

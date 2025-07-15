package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
)

func WithETag(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := &etagResponseRecorder{
			ResponseWriter: w,
			body:           []byte{},
			status:         http.StatusOK,
			header:         http.Header{},
		}

		next(rec, r)

		if rec.status == http.StatusOK {
			hash := md5.Sum(rec.body)
			etag := fmt.Sprintf(`W/"%s"`, hex.EncodeToString(hash[:]))

			if match := r.Header.Get("If-None-Match"); match == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}

			w.Header().Set("ETag", etag)
		}

		for k, vv := range rec.header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(rec.status)
		w.Write(rec.body)
	}
}

type etagResponseRecorder struct {
	http.ResponseWriter
	body   []byte
	status int
	header http.Header
}

func (r *etagResponseRecorder) WriteHeader(code int) {
	r.status = code
}

func (r *etagResponseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return len(b), nil
}

func (r *etagResponseRecorder) Header() http.Header {
	return r.header
}

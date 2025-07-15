package middleware

import (
	"net/http"
	"strconv"
)

func WithCacheControl(maxAge int, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(maxAge))
		next(w, r)
	}
}

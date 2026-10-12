package routes

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

func RequirePassword(password string, next http.Handler) http.Handler {
	want := sha256.Sum256([]byte(password))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}

		_, given, ok := r.BasicAuth()
		got := sha256.Sum256([]byte(given))
		if !ok || subtle.ConstantTimeCompare(got[:], want[:]) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="cs-stalker", charset="UTF-8"`)
			http.Error(w, "Password required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

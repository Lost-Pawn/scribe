package handlers

import (
	"crypto/sha256"
	"net/http"
)

func RequireAuth(username, password string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || sha256.Sum256([]byte(pass)) != sha256.Sum256([]byte(password)) {
			w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

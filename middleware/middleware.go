package middleware

import (
	"gfh.com/web/sessions"
	"net/http"
)

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := sessions.GetUserId(r)
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

package middleware

import (
	"gfh.com/web/session"
	"net/http"
)

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := session.Store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func SaveSession(w http.ResponseWriter, r *http.Request, username string) {
	session, _ := session.Store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
}

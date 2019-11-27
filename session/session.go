package session

import (
	"gfh.com/web/models"
	"github.com/gorilla/sessions"
	"net/http"
)

var Store = sessions.NewCookieStore([]byte("t0p-s3cr3t"))

func SaveSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	userId, err := user.GetId()
	if err != nil {
		return err
	}

	session, _ := Store.Get(r, "session")
	session.Values["user_id"] = userId
	session.Save(r, w)
	return nil
}

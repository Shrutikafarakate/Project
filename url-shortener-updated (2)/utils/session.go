package utils

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("very-secret-key"))

func GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "session")
	return session
}

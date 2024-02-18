package mypackage

import (
    "net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStore().Get(r, "session-name")

	session.Values["auth"] = false
	delete(session.Values, "username")
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
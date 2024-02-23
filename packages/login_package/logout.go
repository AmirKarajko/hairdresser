package login_package

import (
    "net/http"

	"hairdresser/packages/database_package"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")

	session.Values["auth"] = false
	delete(session.Values, "username")
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
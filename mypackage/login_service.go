package mypackage

import (
	"log"
    "net/http"
)

func LoginServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	DatabaseConnect()

	var user string
	var pass string
	var permissionDeleteBill bool

	username := r.FormValue("username")
	password := r.FormValue("password")

	DB.QueryRow("SELECT username, password, permission_delete_bill FROM users WHERE username = ? AND password = ?", username, password).Scan(&user, &pass, &permissionDeleteBill)

	if (user == username && pass == password) {

		session, _ := cookieStore().Get(r, "session-name")

		session.Values["username"] = user
		session.Values["auth"] = true
		session.Values["permission_delete_bill"] = permissionDeleteBill
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		log.Println("Username or password is incorrect")

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	DatabaseDisconnect()
}
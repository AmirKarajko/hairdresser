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

	username := r.FormValue("username")
	password := r.FormValue("password")

	DB.QueryRow("SELECT username, password FROM users WHERE username = ? AND password = ?", username, password).Scan(&user, &pass)

	if (user == username && pass == password) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		log.Println("Username or password is incorrect")

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	DatabaseDisconnect()
}
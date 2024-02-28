package users_package

import (
	"log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
)

type AddUserPageData struct {
	Title string
	Content string
	Username string
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")
	username := session.Values["username"].(string)

	if username != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if r.Method == http.MethodPost {
		database_package.DatabaseConnect()

		inputUsername := r.FormValue("username")
		inputPassword := r.FormValue("password")

		database_package.DB.QueryRow("INSERT INTO users (username, password, permission_delete_bill, permission_delete_service) VALUES (?, ?, 0, 0)", inputUsername, inputPassword)
	
		database_package.DatabaseDisconnect()
	}

	data := AddUserPageData {
		Title: "Hairdresser | Add User",
		Content: "This is a hairdresser web application.",
		Username: username,
	}

	tmpl, err := template.ParseFiles("pages/adduser.html", "pages/navbar.html")

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
package users_package

import (
	"log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
	"hairdresser/packages/utils_package"
)

type AddUserPageData struct {
	Title string
	Content string

	IsAdmin bool
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")
	isAuthenticated := session.Values["authenticated"].(bool)
	isAdmin := session.Values["is_admin"].(bool)

	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !isAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		inputUsername := r.FormValue("username")
		inputPassword := r.FormValue("password")

		AddUser(inputUsername, inputPassword)
	}

	data := AddUserPageData {
		Title: "Hairdresser | Add User",
		Content: "This is a hairdresser web application.",

		IsAdmin: isAdmin,
	}

	tmpl, err := template.ParseFiles("pages/users/add_user.html", "pages/navbar.html")

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

func AddUser(username string, password string) {
	hashedPassword := utils_package.HashPassword(password)

	database_package.DatabaseConnect()
	_, err := database_package.DB.Exec("INSERT INTO users (username, password, permission_delete_bill, permission_delete_service) VALUES (?, ?, 0, 0)", username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
	database_package.DatabaseDisconnect()
}
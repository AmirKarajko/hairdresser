package users_package

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
)

type UsersData struct {
	ID int
	USERNAME string
	PASSWORD string
	ISADMIN bool
}

var Users []UsersData

type UsersPageData struct {
	Title string
	Content string

	Users []UsersData
	IsAdmin bool
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")

	if session.Values["authenticated"] == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	isAuthenticated := session.Values["authenticated"].(bool)
	if !isAuthenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	isAdmin := session.Values["is_admin"].(bool)

	if !isAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	LoadUsersData()

	data := UsersPageData {
		Title: "Hairdresser | Users",
		Content: "This is a hairdresser web application.",
		
		Users: Users,
		IsAdmin: isAdmin,
	}

	tmpl, err := template.ParseFiles("pages/users/users.html", "pages/navbar.html")

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

func LoadUsersData() {
	Users = nil

	database_package.DatabaseConnect()

	rows, err := database_package.DB.Query("SELECT id, username, password, is_admin FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			userID int
			userUsername string
			userPassword string
			userIsAdmin bool
		)

		err := rows.Scan(&userID, &userUsername, &userPassword, &userIsAdmin)
		if err != nil {
			log.Fatal(err)
		}

		Users = append(Users, UsersData{userID, userUsername, userPassword, userIsAdmin})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}
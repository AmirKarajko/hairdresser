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
}

var Users []UsersData

type UsersPageData struct {
	Title string
	Content string
	Username string
	Users []UsersData
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")
	isAuthenticated := session.Values["authenticated"].(bool)
	username := session.Values["username"].(string)

	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	if (username != "admin") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	LoadUsersData()

	data := UsersPageData {
		Title: "Hairdresser | Users",
		Content: "This is a hairdresser web application.",
		Username: username,
		Users: Users,
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

	rows, err := database_package.DB.Query("SELECT id, username, password FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			userID int
			userUsername string
			userPassword string
		)

		err := rows.Scan(&userID, &userUsername, &userPassword)
		if err != nil {
			log.Fatal(err)
		}

		Users = append(Users, UsersData{userID, userUsername, userPassword})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}
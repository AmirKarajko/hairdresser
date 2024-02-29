package packages

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
)

type UsersData struct {
	ID int
	USERNAME string
}

var Users []UsersData

type AdminPageData struct {
	Title string
	Content string
	Username string
	Users []UsersData
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
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

	data := AdminPageData {
		Title: "Hairdresser | Admin",
		Content: "This is a hairdresser web application.",
		Username: username,
		Users: Users,
	}

	tmpl, err := template.ParseFiles("pages/admin.html", "pages/navbar.html")

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

	rows, err := database_package.DB.Query("SELECT id, username FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			userID int
			userUsername string
		)

		err := rows.Scan(&userID, &userUsername)
		if err != nil {
			log.Fatal(err)
		}

		Users = append(Users, UsersData{userID, userUsername})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}
package users_package

import (
	"log"
	"html/template"
    "net/http"
	"strconv"

	"hairdresser/packages/database_package"
	"hairdresser/packages/utils_package"
)

type EditUserPageData struct {
	Title string
	Content string
	
	Id int
	Username string
	Password string

	IsAdmin bool
}

var (
	id int
	username string
	password string
)

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
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

	idStr := r.URL.Path[len("/edit_user/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	GetUserData(id)

	if r.Method == http.MethodPost {
		inputUsername := r.FormValue("username")
		inputPassword := r.FormValue("password")

		UpdateUser(id, inputUsername, inputPassword)
	}	

	data := EditUserPageData {
		Title: "Hairdresser | Edit User",
		Content: "This is a hairdresser web application.",
		
		Id: id,
		Username: username,
		Password: password,

		IsAdmin: isAdmin,
	}

	tmpl, err := template.ParseFiles("pages/users/edit_user.html", "pages/navbar.html")

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

func GetUserData(id int) {
	database_package.DatabaseConnect()

	rows, err := database_package.DB.Query("SELECT username, password FROM users WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&username, &password)

		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}

func UpdateUser(id int, username string, password string) {
	hashedPassword := utils_package.HashPassword(password)

	database_package.DatabaseConnect()
	_, err := database_package.DB.Exec("UPDATE users SET username = ?, password = ? WHERE id = ?", username, hashedPassword, id)
	if err != nil {
		log.Fatal(err)
	}
	database_package.DatabaseDisconnect()
}
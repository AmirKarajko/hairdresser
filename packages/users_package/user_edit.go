package users_package

import (
	"log"
	"html/template"
    "net/http"
	"strconv"

	"hairdresser/packages/database_package"
)

type EditUserPageData struct {
	Title string
	Content string
	Username string
	Id int
	User string
	Password string
}

var user string
var password string

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")
	username := session.Values["username"].(string)

	if username != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	idStr := r.URL.Path[len("/edituser/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		database_package.DatabaseConnect()

		inputUsername := r.FormValue("username")
		inputPassword := r.FormValue("password")

		database_package.DB.QueryRow("UPDATE users SET username = ?, password = ? WHERE id = ?", inputUsername, inputPassword, id)
	
		database_package.DatabaseDisconnect()
	}

	GetUserData(id)

	data := EditUserPageData {
		Title: "Hairdresser | Edit User",
		Content: "This is a hairdresser web application.",
		Username: username,
		Id: id,
		User: user,
		Password: password,
	}

	tmpl, err := template.ParseFiles("pages/edituser.html", "pages/navbar.html")

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
		err := rows.Scan(&user, &password)

		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}
package login_package

import (
    "log"
	"html/template"
    "net/http"
	"database/sql"

	"hairdresser/packages/database_package"
	"hairdresser/packages/utils_package"
)

type LoginPageData struct {
	Title string
	Content string
	ErrorMessage string
}

var (
	errorMessage string
	permissionDeleteBill bool
	permissionDeleteService bool
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")

	if r.Method == http.MethodPost {
		inputUsername := r.FormValue("username")
		inputPassword := r.FormValue("password")

		if LoginUser(inputUsername, inputPassword) {
			session.Values["username"] = inputUsername
			session.Values["authenticated"] = true
			session.Values["permission_delete_bill"] = permissionDeleteBill
			session.Values["permission_delete_service"] = permissionDeleteService
			session.Save(r, w)
		}
	}

	isAuthenticated := session.Values["authenticated"].(bool)
	if isAuthenticated {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

	data := LoginPageData {
		Title: "Hairdresser | Login",
		Content: "This is a hairdresser web application.",
		ErrorMessage: errorMessage,
	}

	tmpl, err := template.ParseFiles("pages/login.html")

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

func LoginUser(username string, password string) bool {
	database_package.DatabaseConnect()
	row := database_package.DB.QueryRow("SELECT id, password, permission_delete_bill, permission_delete_service FROM users WHERE username = ?", username)

	var (
		id int
		pass string
	)

	err := row.Scan(&id, &pass, &permissionDeleteBill, &permissionDeleteService)
	if err != nil {
		if err == sql.ErrNoRows {
			errorMessage = "Error: User does not exist."
		} else {
			log.Fatal(err)
		}
		return false
	}

	hashedPassword := utils_package.HashPassword(password)

	if hashedPassword == pass {
		errorMessage = ""
		return true
	} else {
		errorMessage = "Error: Incorrect password"
		return false
	}

	database_package.DatabaseDisconnect()

	return false
}
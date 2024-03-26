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
	PageTitle string
	Title string
	ErrorMessage string
}

var (
	errorMessage string

	id int
	username string
	password string
	permissionDeleteBill bool
	permissionDeleteService bool
	isAdmin bool
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
			session.Values["is_admin"] = isAdmin
			session.Save(r, w)
		}
	}

	if session.Values["authenticated"] != nil {
		isAuthenticated := session.Values["authenticated"].(bool)
		if isAuthenticated {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	data := LoginPageData {
		PageTitle: "Hairdresser | Login",
		Title: "Welcome to Hairdresser",
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

func LoginUser(user string, pass string) bool {
	database_package.DatabaseConnect()
	row := database_package.DB.QueryRow("SELECT id, password, permission_delete_bill, permission_delete_service, is_admin FROM users WHERE username = ?", user)

	err := row.Scan(&id, &password, &permissionDeleteBill, &permissionDeleteService, &isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			errorMessage = "Error: User does not exist."
		} else {
			log.Fatal(err)
		}
		return false
	}

	hashedPassword := utils_package.HashPassword(pass)

	if hashedPassword == password {
		errorMessage = ""
		return true
	} else {
		errorMessage = "Error: Incorrect password"
		return false
	}

	database_package.DatabaseDisconnect()

	return false
}
package login_package

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
)

type LoginPageData struct {
	Title string
	Content string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		database_package.DatabaseConnect()

		var user string
		var pass string
		var permissionDeleteBill bool
		var permissionDeleteService bool
	
		username := r.FormValue("username")
		password := r.FormValue("password")
	
		database_package.DB.QueryRow("SELECT username, password, permission_delete_bill, permission_delete_service FROM users WHERE username = ? AND password = ?", username, password).Scan(&user, &pass, &permissionDeleteBill, &permissionDeleteService)
	
		if (user == username && pass == password) {
			session, _ := database_package.CookieStore().Get(r, "session-name")
	
			session.Values["username"] = user
			session.Values["auth"] = true
			session.Values["permission_delete_bill"] = permissionDeleteBill
			session.Values["permission_delete_service"] = permissionDeleteService
			session.Save(r, w)
	
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {
			log.Println("Username or password is incorrect")
	
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	
		database_package.DatabaseDisconnect()
		return
	}
	
	session, _ := database_package.CookieStore().Get(r, "session-name")
	authenticated := session.Values["auth"]

	if authenticated == true {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
	
	data := LoginPageData {
		Title: "Hairdresser | Login",
		Content: "This is a hairdresser web application.",
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
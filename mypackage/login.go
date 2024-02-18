package mypackage

import (
    "log"
	"html/template"
    "net/http"
)

type LoginPageData struct {
	Title string
	Content string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	
	session, _ := cookieStore().Get(r, "session-name")

	authenticated := session.Values["auth"]

	if authenticated == true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
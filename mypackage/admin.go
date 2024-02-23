package mypackage

import (
    "log"
	"html/template"
    "net/http"
)

type AdminPageData struct {
	Title string
	Content string
	Username string
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	
	session, _ := cookieStore().Get(r, "session-name")

	authenticated := session.Values["auth"]
	username := session.Values["username"].(string)

	if authenticated == false {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		if (username != "admin") {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			data := AdminPageData {
				Title: "Hairdresser | Admin",
				Content: "This is a hairdresser web application.",
				Username: username,
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
	}
}
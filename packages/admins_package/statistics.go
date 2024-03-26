package admins_package

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
)

type StatisticsPageData struct {
	PageTitle string
	Title string
	Content string

	IsAdmin bool
}

func StatisticsHandler(w http.ResponseWriter, r *http.Request) {
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

	data := StatisticsPageData {
		PageTitle: "Hairdresser | Statistics",
		Title: "Statistics",
		Content: "This is a hairdresser web application.",
		
		IsAdmin: isAdmin,
	}

	tmpl, err := template.ParseFiles("pages/admin/statistics.html", "pages/navbar.html")
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
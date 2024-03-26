package packages

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
	"hairdresser/packages/services_package"
)

type ServicesPageData struct {
	PageTitle string
	Title string
	Content string
	Services []services_package.ServicesData

	PermissionDeleteService bool
	IsAdmin bool
}

func ServicesHandler(w http.ResponseWriter, r *http.Request) {
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

	permissionDeleteService := session.Values["permission_delete_service"].(bool)
	isAdmin := session.Values["is_admin"].(bool)

	services_package.LoadServicesData()

	data := ServicesPageData {
		PageTitle: "Hairdresser | Services",
		Title: "Services",
		Content: "This is a hairdresser web application.",
		Services: services_package.Services,
		
		PermissionDeleteService: permissionDeleteService,
		IsAdmin: isAdmin,
	}

	tmpl, err := template.ParseFiles("pages/services.html", "pages/navbar.html")

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
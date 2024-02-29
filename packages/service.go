package packages

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
	"hairdresser/packages/services_package"
)

type ServicePageData struct {
	Title string
	Content string
	Services []services_package.ServicesData
	Username string
	PermissionDeleteService bool
}

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")
	isAuthenticated := session.Values["authenticated"].(bool)
	username := session.Values["username"].(string)
	permissionDeleteService := session.Values["permission_delete_service"].(bool)

	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	services_package.LoadServicesData()

	data := ServicePageData {
		Title: "Hairdresser | Service",
		Content: "This is a hairdresser web application.",
		Services: services_package.Services,
		Username: username,
		PermissionDeleteService: permissionDeleteService,
	}

	tmpl, err := template.ParseFiles("pages/service.html", "pages/navbar.html")

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
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

	PermissionDeleteService bool
	IsAdmin bool
}

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")
	isAuthenticated := session.Values["authenticated"].(bool)
	permissionDeleteService := session.Values["permission_delete_service"].(bool)
	isAdmin := session.Values["is_admin"].(bool)

	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	services_package.LoadServicesData()

	data := ServicePageData {
		Title: "Hairdresser | Service",
		Content: "This is a hairdresser web application.",
		Services: services_package.Services,
		
		PermissionDeleteService: permissionDeleteService,
		IsAdmin: isAdmin,
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
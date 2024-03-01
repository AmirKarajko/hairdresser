package packages

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
	"hairdresser/packages/bills_package"
	"hairdresser/packages/services_package"
)

type HomePageData struct {
	Title string
	Content string
	Username string
	PermissionDeleteBill bool
	Bills []bills_package.BillsData
	Services []services_package.ServicesData
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")
	isAuthenticated := session.Values["authenticated"].(bool)
	username := session.Values["username"].(string)
	permissionDeleteBill := session.Values["permission_delete_bill"].(bool)

	if !isAuthenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	bills_package.LoadBillsData()
	services_package.LoadServicesData()

	data := HomePageData {
		Title: "Hairdresser",
		Content: "This is a hairdresser web application.",
		Username: username,
		PermissionDeleteBill: permissionDeleteBill,
		Bills: bills_package.Bills,
		Services: services_package.Services,
	}

	tmpl, err := template.ParseFiles("pages/home.html", "pages/navbar.html")
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
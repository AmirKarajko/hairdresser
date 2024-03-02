package packages

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
	"hairdresser/packages/bills_package"
	"hairdresser/packages/services_package"
)

type DashboardPageData struct {
	Title string
	Content string
	Bills []bills_package.BillsData
	Services []services_package.ServicesData

	PermissionDeleteBill bool
	IsAdmin bool
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")
	isAuthenticated := session.Values["authenticated"].(bool)
	permissionDeleteBill := session.Values["permission_delete_bill"].(bool)
	isAdmin := session.Values["is_admin"].(bool)

	if !isAuthenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	bills_package.LoadBillsData()
	services_package.LoadServicesData()

	data := DashboardPageData {
		Title: "Hairdresser",
		Content: "This is a hairdresser web application.",
		Bills: bills_package.Bills,
		Services: services_package.Services,
		
		PermissionDeleteBill: permissionDeleteBill,
		IsAdmin: isAdmin,
	}

	tmpl, err := template.ParseFiles("pages/dashboard.html", "pages/navbar.html")
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
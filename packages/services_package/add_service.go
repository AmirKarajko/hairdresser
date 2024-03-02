package services_package

import (
	"net/http"
	
	"hairdresser/packages/database_package"
)

func AddServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	database_package.DatabaseConnect()

	serviceName := r.FormValue("serviceName")
	servicePrice := r.FormValue("servicePrice")

	database_package.DB.QueryRow("INSERT INTO services (name, price) VALUES (?, ?)", serviceName, servicePrice)

	database_package.DatabaseDisconnect()

	http.Redirect(w, r, "/services", http.StatusSeeOther)
}
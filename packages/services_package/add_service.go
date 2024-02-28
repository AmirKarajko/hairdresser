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

	sName := r.FormValue("service-name")
	sPrice := r.FormValue("service-price")

	query := "INSERT INTO services (name, price) VALUES (?, ?)"
	database_package.DB.QueryRow(query, sName, sPrice)

	database_package.DatabaseDisconnect()

	http.Redirect(w, r, "/service", http.StatusSeeOther)
}
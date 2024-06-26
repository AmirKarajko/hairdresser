package services_package

import (
	"net/http"
	"strconv"

	"hairdresser/packages/database_package"
)

func DeleteServiceHandler(w http.ResponseWriter, r *http.Request) {
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

	if permissionDeleteService || isAdmin {
		database_package.DatabaseConnect()

		idStr := r.URL.Path[len("/delete_service/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
	
		_, err = database_package.DB.Exec("DELETE FROM services WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Failed to delete service", http.StatusInternalServerError)
			return
		}
	
		database_package.DatabaseDisconnect()
	}

	http.Redirect(w, r, "/services", http.StatusSeeOther)
}

package bills_package

import (
	"net/http"
	"strconv"

	"hairdresser/packages/database_package"
)

func DeleteBillHandler(w http.ResponseWriter, r *http.Request) {
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

	permissionDeleteBill := session.Values["permission_delete_bill"].(bool)
	isAdmin := session.Values["is_admin"].(bool)

	if permissionDeleteBill || isAdmin {
		database_package.DatabaseConnect()

		idStr := r.URL.Path[len("/delete_bill/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
	
		_, err = database_package.DB.Exec("DELETE FROM bills WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Failed to delete bill", http.StatusInternalServerError)
			return
		}

		database_package.DatabaseDisconnect()
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
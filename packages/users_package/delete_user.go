package users_package

import (
    "net/http"
	"strconv"

	"hairdresser/packages/database_package"
)

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

	database_package.DatabaseConnect()

	idStr := r.URL.Path[len("/delete_user/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = database_package.DB.Exec("DELETE FROM users WHERE id = ? AND is_admin = 0", id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	database_package.DatabaseDisconnect()

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
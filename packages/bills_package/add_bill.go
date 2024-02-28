package bills_package

import (
	"net/http"
	
	"hairdresser/packages/database_package"
)

func AddBillHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	database_package.DatabaseConnect()

	session, _ := database_package.CookieStore().Get(r, "session-name")
	username := session.Values["username"].(string)
	sID := r.FormValue("service-list")

	database_package.DB.QueryRow("INSERT INTO bills (user, service) VALUES ((SELECT ID FROM users WHERE username = ?), ?)", username, sID)

	database_package.DatabaseDisconnect()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
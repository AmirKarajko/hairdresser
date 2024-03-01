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
	serviceId := r.FormValue("serviceList")

	database_package.DB.QueryRow("INSERT INTO bills (user, service) VALUES ((SELECT ID FROM users WHERE username = ?), ?)", username, serviceId)

	database_package.DatabaseDisconnect()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
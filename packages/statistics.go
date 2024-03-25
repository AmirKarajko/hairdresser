package packages

import (
    "log"
	"html/template"
    "net/http"
	"encoding/json"

	"hairdresser/packages/database_package"
)

type UsersData struct {
	ID string
	RESULT float32
}

var Users []UsersData

type StatisticsPageData struct {
	Title string
	Content string

	Users []UsersData
	IsAdmin bool
}

func StatisticsHandler(w http.ResponseWriter, r *http.Request) {
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

	data := StatisticsPageData {
		Title: "Hairdresser | Statistics",
		Content: "This is a hairdresser web application.",
		
		Users: Users,
		IsAdmin: isAdmin,
	}

	tmpl, err := template.ParseFiles("pages/statistics.html", "pages/navbar.html")
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

func LoadUsersResult() {
	Users = nil

	database_package.DatabaseConnect()

	rows, err := database_package.DB.Query("SELECT u.username, SUM(s.price) AS result FROM users u INNER JOIN bills b ON u.id = b.user INNER JOIN services s ON b.service = s.ID GROUP BY b.user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			userID string
			userResult float32
		)

		err := rows.Scan(&userID, &userResult)
		if err != nil {
			log.Fatal(err)
		}

		Users = append(Users, UsersData{userID, userResult})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}

func GetUsersResultData(w http.ResponseWriter, r *http.Request) {
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

	LoadUsersResult()

	jsonData, err := json.Marshal(Users)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "Failed to write JSON response", http.StatusInternalServerError)
		return
	}
}
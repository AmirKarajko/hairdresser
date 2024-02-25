package packages

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
)

type AdminPageData struct {
	Title string
	Content string
	Username string
	Users [][]interface{}
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")
	authenticated := session.Values["auth"]
	username := session.Values["username"].(string)

	if authenticated == false {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if (username != "admin") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		data := AdminPageData {
			Title: "Hairdresser | Admin",
			Content: "This is a hairdresser web application.",
			Username: username,
			Users: [][]interface{}{
			},
		}

		data.GetUsersData()
	
		tmpl, err := template.ParseFiles("pages/admin.html", "pages/navbar.html")
	
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
}

func (d *AdminPageData) GetUsersData() {
	database_package.DatabaseConnect()

	rows, err := database_package.DB.Query("SELECT id, username FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		var userUsername string

		err := rows.Scan(&userID, &userUsername)

		if err != nil {
			log.Fatal(err)
		}

		d.InsertUserIntoData(userID, userUsername)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}

func (d *AdminPageData) InsertUserIntoData(id int, username string) {
	row1 := []interface{}{id, username}

	d.Users = append(d.Users, row1)
}
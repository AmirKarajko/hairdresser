package packages

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
)

type ServicePageData struct {
	Title string
	Content string
	Services [][]interface{}
	Bills [][]interface{}
	Username string
	PermissionDeleteService bool
}

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")
	isAuthenticated := session.Values["authenticated"].(bool)
	username := session.Values["username"].(string)
	permissionDeleteService := session.Values["permission_delete_service"].(bool)

	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	data := ServicePageData {
		Title: "Hairdresser | Service",
		Content: "This is a hairdresser web application.",
		Services: [][]interface{}{
		},
		Bills: [][]interface{}{
		},
		Username: username,
		PermissionDeleteService: permissionDeleteService,
	}

	data.GetServicesData()

	tmpl, err := template.ParseFiles("pages/service.html", "pages/navbar.html")

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

func (d *ServicePageData) GetServicesData() {
	database_package.DatabaseConnect()

	rows, err := database_package.DB.Query("SELECT id, name, price FROM services")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var serviceID int
		var serviceName string
		var servicePrice float32

		err := rows.Scan(&serviceID, &serviceName, &servicePrice)

		if err != nil {
			log.Fatal(err)
		}

		d.InsertServiceIntoData(serviceID, serviceName, servicePrice)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}

func (d *ServicePageData) InsertServiceIntoData(id int, name string, price float32) {
	row1 := []interface{}{id, name, price}

	d.Services = append(d.Services, row1)
}
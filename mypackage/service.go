package mypackage

import (
    "log"
	"html/template"
    "net/http"
	"fmt"
)

type ServicePageData struct {
	Title string
	Content string
	Services [][]interface{}
	Bills [][]interface{}
}

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	data := ServicePageData {
		Title: "Hairdresser | Service",
		Content: "This is a hairdresser web application.",
		Services: [][]interface{}{
		},
		Bills: [][]interface{}{
		},
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
	DatabaseConnect()

	rows, err := DB.Query("SELECT id, name, price FROM services")
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

	DatabaseDisconnect()
}

func AddServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}


	DatabaseConnect()

	sName := r.FormValue("service-name")
	sPrice := r.FormValue("service-price")

	query := "INSERT INTO services (name, price) VALUES (?, ?)"
	DB.QueryRow(query, sName, sPrice)

	fmt.Println("Service added")

	DatabaseDisconnect()


	http.Redirect(w, r, "/service", http.StatusSeeOther)
}

func (d *ServicePageData) InsertServiceIntoData(id int, name string, price float32) {
	row1 := []interface{}{id, name, price}

	d.Services = append(d.Services, row1)
}
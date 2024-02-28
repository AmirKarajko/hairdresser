package packages

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
)

type HomePageData struct {
	Title string
	Content string
	Services [][]interface{}
	Bills [][]interface{}
	Username string
	PermissionDeleteBill bool
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database_package.CookieStore().Get(r, "session-name")
	authenticated := session.Values["auth"]
	username := session.Values["username"].(string)
	permissionDeleteBill := session.Values["permission_delete_bill"].(bool)

	if authenticated == false {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	data := HomePageData {
		Title: "Hairdresser",
		Content: "This is a hairdresser web application.",
		Services: [][]interface{}{
		},
		Bills: [][]interface{}{
		},
		Username: username,
		PermissionDeleteBill: permissionDeleteBill,
	}

	data.GetServicesData()
	data.GetBillsData()

	tmpl, err := template.ParseFiles("pages/home.html", "pages/navbar.html")

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

func (d *HomePageData) GetBillsData() {
	database_package.DatabaseConnect()

	rows, err := database_package.DB.Query("SELECT users.id, bills.id, services.name, services.price, bills.date FROM bills INNER JOIN services ON bills.service = services.id INNER JOIN users ON bills.user = users.id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		var billID int
		var serviceName string
		var servicesPrice float32
		var billDate string

		err := rows.Scan(&userID, &billID, &serviceName, &servicesPrice, &billDate)

		if err != nil {
			log.Fatal(err)
		}

		d.InsertBillIntoData(billID, serviceName, servicesPrice, billDate, userID)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}

func (d *HomePageData) InsertBillIntoData(id int, service string, price float32, date string, userId int) {
	row1 := []interface{}{id, service, price, date, userId}

	d.Bills = append(d.Bills, row1)
}

func (d *HomePageData) GetServicesData() {
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

func (d *HomePageData) InsertServiceIntoData(id int, name string, price float32) {
	row1 := []interface{}{id, name, price}

	d.Services = append(d.Services, row1)
}
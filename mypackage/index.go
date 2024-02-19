package mypackage

import (
    "log"
	"html/template"
    "net/http"
	"fmt"
	"strconv"
)

type IndexPageData struct {
	Title string
	Content string
	Services [][]interface{}
	Bills [][]interface{}
	Username string
	PermissionDeleteBill bool
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := cookieStore().Get(r, "session-name")

	authenticated := session.Values["auth"]
	username := session.Values["username"].(string)
	permissionDeleteBill := session.Values["permission_delete_bill"].(bool)

	if authenticated == false {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {

		data := IndexPageData {
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
	
		tmpl, err := template.ParseFiles("pages/index.html", "pages/navbar.html")
	
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

func (d *IndexPageData) GetBillsData() {
	DatabaseConnect()

	rows, err := DB.Query("SELECT bills.id, services.name, services.price, bills.date FROM bills INNER JOIN services ON bills.service = services.id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var billID int
		var serviceName string
		var servicesPrice float32
		var billDate string

		err := rows.Scan(&billID, &serviceName, &servicesPrice, &billDate)

		if err != nil {
			log.Fatal(err)
		}

		d.InsertBillIntoData(billID, serviceName, servicesPrice, billDate)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	DatabaseDisconnect()
}

func (d *IndexPageData) GetServicesData() {
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

func AddBillHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}


	DatabaseConnect()

	sID := r.FormValue("service-list")

	DB.QueryRow("INSERT INTO bills (service) VALUES (?)", sID)

	fmt.Println("Bill added")

	DatabaseDisconnect()


	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteBillHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := cookieStore().Get(r, "session-name")
	
	permissionDeleteBill := session.Values["permission_delete_bill"].(bool)

	if permissionDeleteBill {
		DatabaseConnect()

		idStr := r.URL.Path[len("/deletebill/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
	
		_, err = DB.Exec("DELETE FROM bills WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Failed to delete item", http.StatusInternalServerError)
			return
		}
	
		// fmt.Fprintf(w, "Item ID %d deleted successfully", id)
	
		DatabaseDisconnect()
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (d *IndexPageData) InsertBillIntoData(id int, service string, price float32, date string) {
	row1 := []interface{}{id, service, price, date}

	d.Bills = append(d.Bills, row1)
}

func (d *IndexPageData) InsertServiceIntoData(id int, name string, price float32) {
	row1 := []interface{}{id, name, price}

	d.Services = append(d.Services, row1)
}
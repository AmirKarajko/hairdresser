package main

import (
    "database/sql"
    "fmt"
    "log"
	"html/template"
    "net/http"
	"strconv"

    "github.com/go-sql-driver/mysql"
)

type PageData struct {
	Title string
	Content string
	Services [][]interface{}
	Bills [][]interface{}
}

var db *sql.DB

func databaseConnect() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "",
		Net:    "tcp",
		Addr:   "localhost",
		DBName: "hairdresser",
		AllowNativePasswords: true,
	}
	
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// fmt.Println("Database connected")
}

func databaseDisconnect() {
	defer db.Close()

	// fmt.Println("Database disconnected")
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    http.HandleFunc("/", indexHandler)

	http.HandleFunc("/addbill", addBillHandler)
	http.HandleFunc("/deletebill/", deleteBillHandler)

	http.HandleFunc("/service", serviceHandler)
	http.HandleFunc("/addservice", addServiceHandler)
	http.HandleFunc("/deleteservice/", deleteServiceHandler)

	http.HandleFunc("/calculator", calculatorHandler)

    http.ListenAndServe(":8080", nil)
}

// Index Page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData {
		Title: "Hairdresser",
		Content: "This is a hairdresser web application.",
		Services: [][]interface{}{
		},
		Bills: [][]interface{}{
		},
	}

	data.getServicesData()
	data.getBillsData()

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

func addBillHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}


	databaseConnect()

	sID := r.FormValue("service-list")
	billDate := r.FormValue("bill-date")

	query := "INSERT INTO bills (service, date) VALUES (?, ?)"
	db.QueryRow(query, sID, billDate)

	fmt.Println("Bill added")

	databaseDisconnect()


	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Service Page
func serviceHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData {
		Title: "Hairdresser | Service",
		Content: "This is a hairdresser web application.",
		Services: [][]interface{}{
		},
		Bills: [][]interface{}{
		},
	}

	data.getServicesData()

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

func addServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}


	databaseConnect()

	sName := r.FormValue("service-name")
	sPrice := r.FormValue("service-price")

	query := "INSERT INTO services (name, price) VALUES (?, ?)"
	db.QueryRow(query, sName, sPrice)

	fmt.Println("Service added")

	databaseDisconnect()


	http.Redirect(w, r, "/service", http.StatusSeeOther)
}

// Get services data
func (d *PageData) getServicesData() {
	databaseConnect()

	rows, err := db.Query("SELECT id, name, price FROM services")
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

		d.insertServiceIntoData(serviceID, serviceName, servicePrice)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	databaseDisconnect()
}

func (d *PageData) insertServiceIntoData(id int, name string, price float32) {
	row1 := []interface{}{id, name, price}

	d.Services = append(d.Services, row1)
}

// Get bills data
func (d *PageData) getBillsData() {
	databaseConnect()

	rows, err := db.Query("SELECT id, service, date FROM bills")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var billID int
		var billService int
		var billDate string

		err := rows.Scan(&billID, &billService, &billDate)

		if err != nil {
			log.Fatal(err)
		}

		d.insertBillIntoData(billID, billService, billDate)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	databaseDisconnect()
}

func (d *PageData) insertBillIntoData(id int, service int, date string) {
	row1 := []interface{}{id, service, date}

	d.Bills = append(d.Bills, row1)
}

func deleteBillHandler(w http.ResponseWriter, r *http.Request) {
	databaseConnect()

	idStr := r.URL.Path[len("/deletebill/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM bills WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Item ID %d deleted successfully", id)

	databaseDisconnect()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	databaseConnect()

	idStr := r.URL.Path[len("/deleteservice/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM services WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Item ID %d deleted successfully", id)

	databaseDisconnect()

	http.Redirect(w, r, "/service", http.StatusSeeOther)
}

type CalculatorPageData struct {
	Title string
	Content string
	Result float32
}

// Calculator Page
func calculatorHandler(w http.ResponseWriter, r *http.Request) {
	data := CalculatorPageData {
		Title: "Hairdresser | Calculator",
		Content: "This is a hairdresser web application.",
		Result: 0,
	}

	if r.Method == http.MethodGet {
		billDateFrom := r.FormValue("bill-date-from")
		billDateTo := r.FormValue("bill-date-to")

		data.Result = getCalculatorResult(billDateFrom, billDateTo)
	}

	tmpl, err := template.ParseFiles("pages/calculator.html", "pages/navbar.html")

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

func getCalculatorResult(billDateFrom string, billDateTo string) float32 {
	var result float32

	databaseConnect()

	query := "SELECT SUM(services.price) AS result FROM services INNER JOIN bills ON bills.service = services.ID WHERE bills.date BETWEEN ? AND ?"
	
	db.QueryRow(query, billDateFrom, billDateTo).Scan(&result)

	databaseDisconnect()

	return result
}
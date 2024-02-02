package mypackage

import (
    "log"
	"html/template"
    "net/http"
)

type CalculatorPageData struct {
	Title string
	Content string
	Result float32
}

func CalculatorHandler(w http.ResponseWriter, r *http.Request) {
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

	DatabaseConnect()

	query := "SELECT SUM(services.price) AS result FROM services INNER JOIN bills ON bills.service = services.ID WHERE bills.date BETWEEN ? AND ?"
	
	DB.QueryRow(query, billDateFrom, billDateTo).Scan(&result)

	DatabaseDisconnect()

	return result
}
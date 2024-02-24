package packages

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
)

type CalculatorPageData struct {
	Title string
	Content string
	Result float32
	Username string
}

func CalculatorHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := database_package.CookieStore().Get(r, "session-name")

	authenticated := session.Values["auth"]
	username := session.Values["username"].(string)

	if authenticated == false {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		data := CalculatorPageData {
			Title: "Hairdresser | Calculator",
			Content: "This is a hairdresser web application.",
			Result: 0,
			Username: username,
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
}

func getCalculatorResult(billDateFrom string, billDateTo string) float32 {
	var result float32

	database_package.DatabaseConnect()

	query := "SELECT SUM(services.price) AS result FROM services INNER JOIN bills ON bills.service = services.ID WHERE bills.date BETWEEN ? AND ?"
	
	database_package.DB.QueryRow(query, billDateFrom, billDateTo).Scan(&result)

	database_package.DatabaseDisconnect()

	return result
}
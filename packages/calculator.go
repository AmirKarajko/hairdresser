package packages

import (
    "log"
	"html/template"
    "net/http"

	"hairdresser/packages/database_package"
)

type CalculatorPageData struct {
	PageTitle string
	Title string
	Content string
	Result float32

	IsAdmin bool
}

func CalculatorHandler(w http.ResponseWriter, r *http.Request) {
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

	data := CalculatorPageData {
		PageTitle: "Hairdresser | Calculator",
		Title: "Calculator",
		Content: "This page provides a tool to calculate your earnings.",
		Result: 0,

		IsAdmin: isAdmin,
	}

	if r.Method == http.MethodGet {
		billDateFrom := r.FormValue("billDateFrom")
		billDateTo := r.FormValue("billDateTo")

		data.Result = CalculateResult(billDateFrom, billDateTo)
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

func CalculateResult(billDateFrom string, billDateTo string) float32 {
	var result float32

	database_package.DatabaseConnect()

	database_package.DB.QueryRow("SELECT SUM(services.price) AS result FROM services INNER JOIN bills ON bills.service = services.ID WHERE bills.date BETWEEN ? AND ?", billDateFrom, billDateTo).Scan(&result)

	database_package.DatabaseDisconnect()

	return result
}
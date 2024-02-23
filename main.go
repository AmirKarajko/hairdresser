package main

import (
    "net/http"

	"hairdresser/mypackage"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    http.HandleFunc("/", mypackage.IndexHandler)

	http.HandleFunc("/admin", mypackage.AdminHandler)

	http.HandleFunc("/login", mypackage.LoginHandler)
	http.HandleFunc("/logout", mypackage.LogoutHandler)
	http.HandleFunc("/login_service", mypackage.LoginServiceHandler)

	http.HandleFunc("/addbill", mypackage.AddBillHandler)
	http.HandleFunc("/deletebill/", mypackage.DeleteBillHandler)

	http.HandleFunc("/service", mypackage.ServiceHandler)
	http.HandleFunc("/addservice", mypackage.AddServiceHandler)
	http.HandleFunc("/deleteservice/", mypackage.DeleteServiceHandler)

	http.HandleFunc("/calculator", mypackage.CalculatorHandler)

    http.ListenAndServe(":8080", nil)
}
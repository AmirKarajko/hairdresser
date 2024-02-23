package main

import (
    "net/http"

	"hairdresser/packages"
	"hairdresser/packages/login_package"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    http.HandleFunc("/", packages.IndexHandler)

	http.HandleFunc("/admin", packages.AdminHandler)

	http.HandleFunc("/login", login_package.LoginHandler)
	http.HandleFunc("/logout", login_package.LogoutHandler)
	http.HandleFunc("/login_service", login_package.LoginServiceHandler)

	http.HandleFunc("/addbill", packages.AddBillHandler)
	http.HandleFunc("/deletebill/", packages.DeleteBillHandler)

	http.HandleFunc("/service", packages.ServiceHandler)
	http.HandleFunc("/addservice", packages.AddServiceHandler)
	http.HandleFunc("/deleteservice/", packages.DeleteServiceHandler)

	http.HandleFunc("/calculator", packages.CalculatorHandler)

    http.ListenAndServe(":8080", nil)
}
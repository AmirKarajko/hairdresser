package main

import (
    "net/http"

	"hairdresser/packages"
	"hairdresser/packages/login_package"
	"hairdresser/packages/users_package"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", RedirectHandler)

    http.HandleFunc("/home", packages.HomeHandler)

	http.HandleFunc("/admin", packages.AdminHandler)
	
	http.HandleFunc("/adduser", users_package.AddUserHandler)
	http.HandleFunc("/edituser/", users_package.EditUserHandler)
	http.HandleFunc("/deleteuser/", users_package.DeleteUserHandler)

	http.HandleFunc("/login", login_package.LoginHandler)
	http.HandleFunc("/logout", login_package.LogoutHandler)

	http.HandleFunc("/addbill", packages.AddBillHandler)
	http.HandleFunc("/deletebill/", packages.DeleteBillHandler)

	http.HandleFunc("/service", packages.ServiceHandler)
	http.HandleFunc("/addservice", packages.AddServiceHandler)
	http.HandleFunc("/deleteservice/", packages.DeleteServiceHandler)

	http.HandleFunc("/calculator", packages.CalculatorHandler)

    http.ListenAndServe(":8080", nil)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
package main

import (
    "net/http"

	"hairdresser/packages"
	"hairdresser/packages/login_package"
	"hairdresser/packages/users_package"
	"hairdresser/packages/bills_package"
	"hairdresser/packages/services_package"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", RedirectHandler)

    http.HandleFunc("/dashboard", packages.DashboardHandler)

	http.HandleFunc("/users", users_package.UsersHandler)
	http.HandleFunc("/add_user", users_package.AddUserHandler)
	http.HandleFunc("/edit_user/", users_package.EditUserHandler)
	http.HandleFunc("/delete_user/", users_package.DeleteUserHandler)

	http.HandleFunc("/statistics", packages.StatisticsHandler)
	http.HandleFunc("/get_users_result_data", packages.GetUsersResultData)

	http.HandleFunc("/login", login_package.LoginHandler)
	http.HandleFunc("/logout", login_package.LogoutHandler)

	http.HandleFunc("/add_bill", bills_package.AddBillHandler)
	http.HandleFunc("/delete_bill/", bills_package.DeleteBillHandler)

	http.HandleFunc("/services", packages.ServicesHandler)
	http.HandleFunc("/add_service", services_package.AddServiceHandler)
	http.HandleFunc("/delete_service/", services_package.DeleteServiceHandler)

	http.HandleFunc("/calculator", packages.CalculatorHandler)

    http.ListenAndServe(":8080", nil)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
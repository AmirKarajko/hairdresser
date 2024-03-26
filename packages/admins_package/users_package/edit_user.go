package users_package

import (
	"log"
	"html/template"
    "net/http"
	"strconv"

	"hairdresser/packages/database_package"
	"hairdresser/packages/utils_package"
)

type EditUserPageData struct {
	PageTitle string
	Title string
	Content string
	
	Id int
	Username string
	Password string
	PermissionDeleteBill bool
	PermissionDeleteService bool

	IsAdmin bool
}

var (
	id int
	username string
	password string
	permissionDeleteBill bool
	permissionDeleteService bool
)

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
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

	if !isAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.URL.Path[len("/edit_user/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	LoadUserData(id)

	if r.Method == http.MethodPost {
		inputUsername := r.FormValue("username")
		inputPassword := r.FormValue("password")
		inputPermissionDeleteBill := r.Form.Get("permissionDeleteBill")
		inputPermissionDeleteService := r.Form.Get("permissionDeleteService")

		isPermissionDeleteBillChecked := inputPermissionDeleteBill == "on"
		isPermissionDeleteServiceChecked := inputPermissionDeleteService == "on"
		
		UpdateUser(id, inputUsername, inputPassword, isPermissionDeleteBillChecked, isPermissionDeleteServiceChecked)

		LoadUserData(id)
	}	

	data := EditUserPageData {
		PageTitle: "Hairdresser | Edit User",
		Title: "Edit User",
		Content: "This is a hairdresser web application.",
		
		Id: id,
		Username: username,
		Password: password,
		PermissionDeleteBill: permissionDeleteBill,
		PermissionDeleteService: permissionDeleteService,

		IsAdmin: isAdmin,
	}

	tmpl, err := template.ParseFiles("pages/admin/users/edit_user.html", "pages/navbar.html")

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

func LoadUserData(id int) {
	database_package.DatabaseConnect()

	rows, err := database_package.DB.Query("SELECT username, password, permission_delete_bill, permission_delete_service FROM users WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&username, &password, &permissionDeleteBill, &permissionDeleteService)

		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}

func UpdateUser(id int, username string, password string, permissionDeleteBill bool, permissionDeleteService bool) {
	hashedPassword := utils_package.HashPassword(password)

	database_package.DatabaseConnect()
	_, err := database_package.DB.Exec("UPDATE users SET username = ?, password = ?, permission_delete_bill = ?, permission_delete_service = ? WHERE id = ?", username, hashedPassword, permissionDeleteBill, permissionDeleteService, id)
	if err != nil {
		log.Fatal(err)
	}
	database_package.DatabaseDisconnect()
}
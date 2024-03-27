package bills_package

import (
    "log"

	"hairdresser/packages/database_package"
)

type BillsData struct {
	USER string
	ID int
	SERVICE_NAME string
	SERVICE_PRICE float32
	DATE string
}

var Bills []BillsData

func LoadBillsData() {
	Bills = nil

	database_package.DatabaseConnect()

	rows, err := database_package.DB.Query("SELECT users.username, bills.id, services.name, services.price, bills.date FROM bills INNER JOIN services ON bills.service = services.id INNER JOIN users ON bills.user = users.id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			user string
			billID int
			serviceName string
			servicePrice float32
			billDate string
		)

		err := rows.Scan(&user, &billID, &serviceName, &servicePrice, &billDate)
		if err != nil {
			log.Fatal(err)
		}

		Bills = append(Bills, BillsData{user, billID, serviceName, servicePrice, billDate})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}
package services_package

import (
    "log"

	"hairdresser/packages/database_package"
)

type ServicesData struct {
	ID int
	NAME string
	PRICE float32
}

var Services []ServicesData

func LoadServicesData() {
	Services = nil

	database_package.DatabaseConnect()

	rows, err := database_package.DB.Query("SELECT id, name, price FROM services")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			serviceID int
			serviceName string
			servicePrice float32
		)

		err := rows.Scan(&serviceID, &serviceName, &servicePrice)
		if err != nil {
			log.Fatal(err)
		}

		Services = append(Services, ServicesData{serviceID, serviceName, servicePrice})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	database_package.DatabaseDisconnect()
}
package mypackage

import (
    "database/sql"
    "log"

    "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DatabaseConnect() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "",
		Net:    "tcp",
		Addr:   "localhost",
		DBName: "hairdresser",
		AllowNativePasswords: true,
	}
	
	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// fmt.Println("Database connected")
}

func DatabaseDisconnect() {
	defer DB.Close()

	// fmt.Println("Database disconnected")
}
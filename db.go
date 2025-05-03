package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	// This is your Railway public URL
	dsn := "root:rIMXsXaJfXnQHgsLRiGqpAvTsgTAgMzU@tcp(ballast.proxy.rlwy.net:49782)/railway?parseTime=true&tls=false"


	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	log.Println("Connected to Railway DB!")
}

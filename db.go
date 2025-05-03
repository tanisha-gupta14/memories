package main

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "root:tan2005@tcp(127.0.0.1:3306)/memorydb")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
}

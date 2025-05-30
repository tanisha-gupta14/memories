package main

import (
	"database/sql"
	"log"
	"os"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using Railway env vars if available")
	}

	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?parseTime=true"

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Failed to open DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Failed to ping DB:", err)
	}

	log.Println("✅ Connected to DB!")
}

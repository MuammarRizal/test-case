package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/test_case?parseTime=true"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}
	log.Println("Database connected successfully")
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return DB
}
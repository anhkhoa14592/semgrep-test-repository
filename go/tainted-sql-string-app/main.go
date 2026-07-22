package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Single entry point for this application.
func main() {
	db, dbErr := sql.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/dbname")
	if dbErr != nil {
		log.Fatalf("failed to open db: %v", dbErr)
	}
	defer db.Close()

	log.Println("listening on :8080")
	log.Fatal(StartServer(db))
}

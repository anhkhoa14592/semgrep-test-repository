package main

import (
	"database/sql"
	"log"
)

// Single entry point for this application (the only main() in this
// folder - unlike the flat go/ fixtures directory, this package is meant
// to actually build and run as one program).
func main() {
	db, err := sql.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/dbname")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	log.Println("listening on :8080")
	log.Fatal(StartServer(db))
}

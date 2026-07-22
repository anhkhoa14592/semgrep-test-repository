package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	kafka "github.com/segmentio/kafka-go"
)

// Single entry point for this application.
func main() {
	db, err := sql.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/dbname")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "order-events",
		GroupID: "order-status-updater",
	})
	defer reader.Close()

	log.Println("consuming order-events from Kafka")
	if err := RunWithRouter(context.Background(), reader, db); err != nil {
		log.Fatalf("consumer stopped: %v", err)
	}
}

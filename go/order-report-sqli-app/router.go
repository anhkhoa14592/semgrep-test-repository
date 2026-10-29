package main

// Router layer: map HTTP path -> controller.

import (
	"database/sql"
	"net/http"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	svc := NewOrderReportService(NewOrderRepository(db))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/orders/report", OrderReportController(svc))
	return mux
}

// StartServer starts the HTTP server. Called from main.go.
func StartServer(db *sql.DB) error {
	return http.ListenAndServe(":8080", NewRouter(db))
}

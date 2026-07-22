package main

// Router layer (optional - see mvc_simulated_request.go for the router-less
// variant of the same flow).

import (
	"database/sql"
	"net/http"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/table/delete", DeleteTableController(db))
	mux.HandleFunc("/table/delete-safe", DeleteTableSafeController(db))
	mux.HandleFunc("/table/select-sprintf", SelectTableSprintfController(db))
	mux.HandleFunc("/table/select-builder", SelectTableBuilderController(db))
	mux.HandleFunc("/table/select-concat", SelectTableConcatController(db))
	mux.HandleFunc("/table/select-safe", SelectTableSafeController(db))

	return mux
}

// StartServer starts the HTTP server. Called from main.go.
func StartServer(db *sql.DB) error {
	return http.ListenAndServe(":8080", NewRouter(db))
}

package main

// Router layer (optional - see mvc_simulated_request.go for the router-less
// variant of the same flow).
//
// Wires HTTP paths to the controllers in mvc_controller.go, and, for
// contrast, one path directly to an already-existing handler
// (redirectHandler, defined in tn_open-redirect.go) that needs no new
// controller at all because it already has the (w, r) shape of one.

import (
	"database/sql"
	"net/http"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/users/search", SearchUserByEmailController(db))
	mux.HandleFunc("/api/users/search-alt", SearchUserByEmailAltController(db))
	mux.HandleFunc("/api/users/by-age", AgeQueryController(db))
	mux.HandleFunc("/api/users/delete-by-age", DeleteQueryController(db))

	// existing handler wired straight in - no new controller needed.
	mux.HandleFunc("/redirect", redirectHandler)

	return mux
}

// StartServer starts the HTTP server. Called from main.go.
func StartServer(db *sql.DB) error {
	return http.ListenAndServe(":8080", NewRouter(db))
}

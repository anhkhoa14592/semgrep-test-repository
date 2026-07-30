package main

// Router-less variant: instead of going through NewRouter/ServeMux, this
// file builds a request exactly as it would arrive from the internet
// (method, path, query string) with net/http/httptest and dispatches it
// straight to a controller's ServeHTTP. Same source -> controller -> model
// -> sink path as mvc_router.go, minus the router hop.

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
)

func SimulateIncomingSearchRequest(db *sql.DB, email string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/api/users/search?email="+email, nil)
	rec := httptest.NewRecorder()

	controller := SearchUserByEmailController(db)
	controller.ServeHTTP(rec, req)

	return rec
}

func SimulateIncomingAgeRequest(db *sql.DB, age string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/api/users/by-age?age="+age, nil)
	rec := httptest.NewRecorder()

	controller := AgeQueryController(db)
	controller.ServeHTTP(rec, req)

	return rec
}

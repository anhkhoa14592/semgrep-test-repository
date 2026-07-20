package main

// Router-less variant: builds a request exactly as it would arrive from
// the internet (net/http/httptest) and dispatches it straight to a
// controller's ServeHTTP - no router hop.

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
)

func SimulateIncomingDeleteRequest(db *sql.DB, id string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/table/delete?del=del&Id="+id, nil)
	rec := httptest.NewRecorder()

	controller := DeleteTableController(db)
	controller.ServeHTTP(rec, req)

	return rec
}

func SimulateIncomingSelectRequest(db *sql.DB, id string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/table/select-sprintf?del=del&Id="+id, nil)
	rec := httptest.NewRecorder()

	controller := SelectTableSprintfController(db)
	controller.ServeHTTP(rec, req)

	return rec
}

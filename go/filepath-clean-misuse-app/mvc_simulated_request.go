package main

// Router-less variant: builds a request exactly as it would arrive from
// the internet (net/http/httptest) and dispatches it straight to the
// controller's ServeHTTP - no router hop.

import (
	"net/http"
	"net/http/httptest"
)

func SimulateIncomingFileReadRequest(file string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/files/read?file="+file, nil)
	rec := httptest.NewRecorder()

	controller := FileServeController()
	controller.ServeHTTP(rec, req)

	return rec
}

package main

import (
	"log"
	"net/http"
)

// RenderServerError is the one View responsibility this app needs: the
// existing Model functions call panic(err) on a DB error (see
// tp_tainted-sql-string.go), so the Controller recovers and renders a
// proper response instead of taking the whole process down.
func RenderServerError(w http.ResponseWriter, rec interface{}) {
	log.Printf("recovered panic: %v", rec)
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

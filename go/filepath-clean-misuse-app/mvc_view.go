package main

import "net/http"

// RenderMethodNotAllowed is the one View responsibility this app needs,
// used by the Controller before it ever reaches the Model.
func RenderMethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

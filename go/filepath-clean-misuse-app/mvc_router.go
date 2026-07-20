package main

// Router layer (optional - see mvc_simulated_request.go for the router-less
// variant of the same flow).

import "net/http"

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/files/read", FileServeController())
	return mux
}

// StartServer starts the HTTP server. Called from main.go.
func StartServer() error {
	return http.ListenAndServe(":8080", NewRouter())
}

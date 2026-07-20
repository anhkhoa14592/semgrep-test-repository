package main

import (
	"fmt"
	"net/http"
)

func isValid(token string) bool {
	return true
}

func vulnerableHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tok := r.FormValue("token")

	if !isValid(tok) {
		// FIX: escape user input before writing to response
		fmt.Fprintf(w, "Invalid token: %q", html.EscapeString(tok))
		return
	}
}
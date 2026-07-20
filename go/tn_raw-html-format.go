package main

import (
	"fmt"
	"net/http"
	"strings"
)

func profileHandler(w http.ResponseWriter, r *http.Request) {

	// SOURCE
	name := r.FormValue("name")

	/*
		Manual mitigation:
		- strict allowlist
		- only alphanumeric + space
		- no HTML metacharacters possible
	*/
	if !isSafeDisplayName(name) {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	/*
		Semgrep vẫn alert vì:
		- source -> HTML sink
		- KHÔNG dùng html.EscapeString(...)
		- rule chỉ recognize sanitizer đó
	*/
	html := "<h1>Welcome " + name + "</h1>"

	fmt.Fprintf(w, html)
}

func isSafeDisplayName(s string) bool {

	if len(s) == 0 || len(s) > 40 {
		return false
	}

	for _, c := range s {

		// allow:
		// a-z A-Z 0-9 space
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == ' ') {
			return false
		}
	}

	// extra hardening
	if strings.ContainsAny(s, "<>\"'&") {
		return false
	}

	return true
}

func main() {}
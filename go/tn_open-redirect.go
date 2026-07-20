package main

import (
	"net/http"
	"strings"
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {

	// SOURCE
	next := r.FormValue("next")

	/*
		Manual mitigation:
		- chỉ cho relative path
		- reject absolute URL
		- reject protocol-relative URL
	*/
	if !isSafeRedirect(next) {
		http.Error(w, "invalid redirect", http.StatusBadRequest)
		return
	}

	/*
		Semgrep vẫn alert vì:
		- taint source -> sink
		- không có sanitizer trong rule
		- CLEAN chỉ apply khi:
		      "https://..." + input
		  nên case này không match CLEAN
	*/
	http.Redirect(w, r, next, http.StatusFound)
}

func isSafeRedirect(s string) bool {

	// must start with single "/"
	if !strings.HasPrefix(s, "/") {
		return false
	}

	// block protocol-relative redirect
	if strings.HasPrefix(s, "//") {
		return false
	}

	// block CRLF
	if strings.ContainsAny(s, "\r\n") {
		return false
	}

	return true
}

func main() {}
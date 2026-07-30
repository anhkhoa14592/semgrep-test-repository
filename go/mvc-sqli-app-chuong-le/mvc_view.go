package main

// View layer: tiny response renderers used by the MVC controllers in this
// package. Kept intentionally minimal - the vulnerable Model functions
// (bad2/bad3/bad4/bad5 in tn_gosql-sqli.go) still write their own output,
// unchanged, since those files are not modified here.

import "net/http"

func RenderBadRequest(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}

func RenderAccepted(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(msg))
}

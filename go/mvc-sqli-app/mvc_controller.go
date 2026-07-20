package main

// Controller layer.
//
// Each controller here is the point where a request "from the internet"
// first touches this codebase. Two shapes are demonstrated on purpose:
//
//  1. SearchUserByEmailController / SearchUserByEmailAltController extract a
//     plain value (email) from the request and hand it to the Model layer
//     (UserModel, mvc_model.go) - a clean Controller -> Model boundary.
//  2. AgeQueryController / DeleteQueryController forward the whole
//     *http.Request straight into the existing bad2/bad4 functions
//     (tn_gosql-sqli.go), which already do their own extraction
//     (req.FormValue("age")) - a Controller -> pre-existing Model-ish
//     function boundary, matching how those functions were originally
//     written.
//
// In both shapes, the taint source is a value read off *http.Request, and
// the sink is the raw SQL built inside tn_gosql-sqli.go - now reachable
// end-to-end via this file's functions, which is the call-graph edge this
// sample exists to exercise.

import (
	"database/sql"
	"net/http"
)

func SearchUserByEmailController(db *sql.DB) http.HandlerFunc {
	model := NewUserModel(db)

	return func(w http.ResponseWriter, r *http.Request) {
		// SOURCE: value read straight off the incoming request.
		email := r.URL.Query().Get("email")
		if email == "" {
			RenderBadRequest(w, "missing email")
			return
		}

		// forwarded, unsanitized, into the Model layer -> bad3 sink.
		model.SearchByEmail(email)
		RenderAccepted(w, "search dispatched")
	}
}

func SearchUserByEmailAltController(db *sql.DB) http.HandlerFunc {
	model := NewUserModel(db)

	return func(w http.ResponseWriter, r *http.Request) {
		// SOURCE
		email := r.FormValue("email")
		if email == "" {
			RenderBadRequest(w, "missing email")
			return
		}

		// forwarded, unsanitized, into the Model layer -> bad5 sink.
		model.SearchByEmailAlt(email)
		RenderAccepted(w, "search dispatched")
	}
}

func AgeQueryController(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// whole request forwarded as-is; bad2 extracts "age" itself.
		bad2(db, r)
		RenderAccepted(w, "age query dispatched")
	}
}

func DeleteQueryController(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// whole request forwarded as-is; bad4 extracts "age" itself.
		bad4(db, r)
		RenderAccepted(w, "delete-by-age query dispatched")
	}
}

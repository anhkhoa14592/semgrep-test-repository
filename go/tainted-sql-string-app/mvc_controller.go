package main

// Controller layer.
//
// Each controller is where a request "from the internet" first touches
// this codebase. It recovers panics from the Model layer (see mvc_view.go)
// and otherwise just dispatches straight through to TableModel - the
// request itself (with its "del"/"Id" query params, still untouched) is
// what carries the taint from here into tp_tainted-sql-string.go.

import (
	"database/sql"
	"net/http"
)

func withRecover(w http.ResponseWriter, next func()) {
	defer func() {
		if rec := recover(); rec != nil {
			RenderServerError(w, rec)
		}
	}()
	next()
}

func DeleteTableController(db *sql.DB) http.HandlerFunc {
	model := NewTableModel(db)
	return func(w http.ResponseWriter, r *http.Request) {
		withRecover(w, func() { model.Delete(w, r) })
	}
}

func DeleteTableSafeController(db *sql.DB) http.HandlerFunc {
	model := NewTableModel(db)
	return func(w http.ResponseWriter, r *http.Request) {
		withRecover(w, func() { model.DeleteSafe(w, r) })
	}
}

func SelectTableSprintfController(db *sql.DB) http.HandlerFunc {
	model := NewTableModel(db)
	return func(w http.ResponseWriter, r *http.Request) {
		withRecover(w, func() { model.SelectViaSprintf(w, r) })
	}
}

func SelectTableBuilderController(db *sql.DB) http.HandlerFunc {
	model := NewTableModel(db)
	return func(w http.ResponseWriter, r *http.Request) {
		withRecover(w, func() { model.SelectViaBuilder(w, r) })
	}
}

func SelectTableConcatController(db *sql.DB) http.HandlerFunc {
	model := NewTableModel(db)
	return func(w http.ResponseWriter, r *http.Request) {
		withRecover(w, func() { model.SelectViaConcat(w, r) })
	}
}

func SelectTableSafeController(db *sql.DB) http.HandlerFunc {
	model := NewTableModel(db)
	return func(w http.ResponseWriter, r *http.Request) {
		withRecover(w, func() { model.SelectSafe(w, r) })
	}
}

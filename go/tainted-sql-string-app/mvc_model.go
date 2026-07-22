package main

// Model layer.
//
// TableModel wraps *sql.DB and gives the Controller layer a proper Model
// entry point. It delegates straight into the pre-existing
// handler-generator functions in tp_tainted-sql-string.go, which already
// contain the SQL building (vulnerable and safe variants alike) - nothing
// about the actual sink logic is changed here.

import (
	"database/sql"
	"net/http"
)

type TableModel struct {
	DB *sql.DB
}

func NewTableModel(db *sql.DB) *TableModel {
	return &TableModel{DB: db}
}

// Delete -> DeleteHandler: vulnerable, "..." + id string concatenation.
func (m *TableModel) Delete(w http.ResponseWriter, r *http.Request) {
	DeleteHandler(m.DB)(w, r)
}

// DeleteSafe -> DeleteHandlerOk: safe, bind parameter.
func (m *TableModel) DeleteSafe(w http.ResponseWriter, r *http.Request) {
	DeleteHandlerOk(m.DB)(w, r)
}

// SelectViaSprintf -> SelectHandler: vulnerable, fmt.Sprintf.
func (m *TableModel) SelectViaSprintf(w http.ResponseWriter, r *http.Request) {
	SelectHandler(m.DB)(w, r)
}

// SelectViaBuilder -> SelectHandler2: vulnerable, strings.Builder.
func (m *TableModel) SelectViaBuilder(w http.ResponseWriter, r *http.Request) {
	SelectHandler2(m.DB)(w, r)
}

// SelectViaConcat -> SelectHandler3: vulnerable, += concatenation.
func (m *TableModel) SelectViaConcat(w http.ResponseWriter, r *http.Request) {
	SelectHandler3(m.DB)(w, r)
}

// SelectSafe -> SelectHandlerOk: safe, bind parameter.
func (m *TableModel) SelectSafe(w http.ResponseWriter, r *http.Request) {
	SelectHandlerOk(m.DB)(w, r)
}

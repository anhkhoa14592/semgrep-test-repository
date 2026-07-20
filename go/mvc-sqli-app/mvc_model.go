package main

// Model layer.
//
// UserModel wraps *sql.DB and exposes plain, HTTP-agnostic methods. It does
// NOT read anything from *http.Request itself - callers (Controllers) must
// hand it already-extracted values. Internally it delegates to the
// already-existing vulnerable sink functions defined in tn_gosql-sqli.go
// (bad3/bad5), so the actual SQL injection lives exactly where it already
// did; this file only adds the call-graph edge from a proper Model type
// into those pre-existing functions.

import "database/sql"

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

// SearchByEmail delegates to bad3 (tn_gosql-sqli.go), which builds
// "SELECT * FROM users WHERE email='%s'" via fmt.Sprintf and runs it with
// db.Exec. If email is attacker-controlled, this is SQL injection.
func (m *UserModel) SearchByEmail(email string) {
	bad3(m.DB, email)
}

// SearchByEmailAlt delegates to bad5 (tn_gosql-sqli.go), same shape sink,
// exercised as a second call-graph path from the Model layer.
func (m *UserModel) SearchByEmailAlt(email string) {
	bad5(m.DB, email)
}

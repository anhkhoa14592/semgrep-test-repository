package main

import (
	"database/sql"
	"net/http"
)

func safeButStillAlert(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// user-controlled input
		username := r.FormValue("username")

		/*
			Manual escaping / allowlist
			=> reviewer thấy đã mitigate SQLi
		*/
		if !isSafeUsername(username) {
			http.Error(w, "invalid username", http.StatusBadRequest)
			return
		}

		/*
			Vẫn trigger rule vì:
			- source: r.FormValue(...)
			- sink: SQL string concat
			- KHÔNG match sanitizer của rule
		*/
		query := "SELECT * FROM users WHERE username = '" + username + "'"

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		w.Write([]byte("ok"))
	}
}

func isSafeUsername(s string) bool {
	if len(s) == 0 || len(s) > 32 {
		return false
	}

	for _, c := range s {

		// strict allowlist:
		// only a-z A-Z 0-9 _
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_') {
			return false
		}
	}

	return true
}

func main() {}
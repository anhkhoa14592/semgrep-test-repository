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
		*/

		// INLINE isSafeUsername logic
		isSafe := true

		if len(username) == 0 || len(username) > 32 {
			isSafe = false
		} else {
			for _, c := range username {
				// only a-z A-Z 0-9 _
				if !((c >= 'a' && c <= 'z') ||
					(c >= 'A' && c <= 'Z') ||
					(c >= '0' && c <= '9') ||
					c == '_') {
					isSafe = false
					break
				}
			}
		}

		if !isSafe {
			http.Error(w, "invalid username", http.StatusBadRequest)
			return
		}

		/*
			Still triggers rule:
			- tainted input used in SQL concat
			- no recognized sanitizer in Semgrep rule
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

func main() {}
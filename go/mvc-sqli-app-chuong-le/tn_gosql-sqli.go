package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// helper: sanitize email (very strict)
func bad1(req *http.Request) {
	db, err := sql.Open("mysql", "theUser:thePassword@/theDbName")
	if err != nil {
		panic(err)
	}

	ageStr := req.FormValue("age")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return
	}

	query := fmt.Sprintf("SELECT name FROM users WHERE age=%d", age)
	db.Query(query)
}

func bad2(db *sql.DB, req *http.Request) {
	ageStr := req.FormValue("age")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return
	}

	query := fmt.Sprintf("SELECT name FROM users WHERE age=%d", age)
	db.QueryRow(query)
}

func bad3(db *sql.DB, email string) {

	// INLINE isSafeEmail logic
	isSafe := true

	if len(email) < 5 || len(email) > 254 {
		isSafe = false
	} else {
		atCount := 0
		for _, c := range email {
			if c == '@' {
				atCount++
			}
			if !(c == '.' || c == '_' || c == '-' ||
				(c >= 'a' && c <= 'z') ||
				(c >= 'A' && c <= 'Z') ||
				(c >= '0' && c <= '9') ||
				c == '@') {
				isSafe = false
				break
			}
		}
		if atCount != 1 {
			isSafe = false
		}
	}

	if !isSafe {
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE email='%s'", email)
	db.Exec(query)
}

func bad4(db *sql.DB, req *http.Request) {
	ageStr := req.FormValue("age")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return
	}

	db.Exec(fmt.Sprintf("SELECT name FROM users WHERE age=%d", age))
}

func bad5(db *sql.DB, email string) {

	// INLINE isSafeEmail logic
	isSafe := true

	if len(email) < 5 || len(email) > 254 {
		isSafe = false
	} else {
		atCount := 0
		for _, c := range email {
			if c == '@' {
				atCount++
			}
			if !(c == '.' || c == '_' || c == '-' ||
				(c >= 'a' && c <= 'z') ||
				(c >= 'A' && c <= 'Z') ||
				(c >= '0' && c <= '9') ||
				c == '@') {
				isSafe = false
				break
			}
		}
		if atCount != 1 {
			isSafe = false
		}
	}

	if !isSafe {
		return
	}

	db.Exec(fmt.Sprintf("SELECT * FROM users WHERE email='%s'", email))
}

func ok1(db *sql.DB) {
	query := fmt.Sprintf("SELECT * FROM users WHERE email=hello;")
	// ok: gosql-sqli
	db.Exec(query)
}

func ok2(db *sql.DB) {
	query := "SELECT name FROM users WHERE age=" + "3"
	// ok: gosql-sqli
	db.Query(query)
}

func ok3(db *sql.DB) {
	query := "SELECT name FROM users WHERE age="
	query += "3"
	// ok: gosql-sqli
	db.Query(query)
}

func ok4(db *sql.DB) {
	// ok: gosql-sqli
	db.Exec("INSERT INTO users(name, email) VALUES($1, $2)",
		"Jon Calhoun", "jon@calhoun.io")
}

func ok5(db *sql.DB) {
	// ok: gosql-sqli
	db.Exec("SELECT name FROM users WHERE age=" + "3")
}

func ok6(db *sql.DB) {
	// ok: gosql-sqli
	db.Exec(fmt.Sprintf("SELECT * FROM users WHERE email=hello;"))
}

var _ = strings.TrimSpace

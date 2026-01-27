// Test 5
package main 
 
import (
	"database/sql"
	"net/http"
)

func DeleteHandlerTest(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		del := req.URL.Query().Get("del")
		id := req.URL.Query().Get("Id")
		if del == "del" {
			// ruleid: tainted-sql-string test 7
			
			_, err = db.Exec("DELETE FROM table WHERE Id = " + id) // test total scanned targets and findings
			if err != nil {
				panic(err)
			}
		}
	}
}

func DeleteHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		del := req.URL.Query().Get("del")
		id := req.URL.Query().Get("Id")
		if del == "del" {
			// ruleid: tainted-sql-string test 7
			
			_, err = db.Exec("DELETE FROM table WHERE Id = " + id) // test total scanned targets and findings
			if err != nil {
				panic(err)
			}
		}
	}
}

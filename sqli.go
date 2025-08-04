// # test total scanned targets and findings
// test update semgrep version
// test time: 1
package main 
 
import (
	"database/sql"
	"net/http"
)

func DeleteHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		del := req.URL.Query().Get("del")
		id := req.URL.Query().Get("Id")
		if del == "del" {
			// ruleid: tainted-sql-string
			_, err = db.Exec("DELETE FROM table WHERE Id = " + id)
			if err != nil {
				panic(err)
			}
		}
	}
}

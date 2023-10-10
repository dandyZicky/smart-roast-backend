package user

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM mock_user")
	if err != nil {
		panic("Query failed")
	}
	defer rows.Close()

	var result []User

	for rows.Next() {
		var each = User{}
		var err = rows.Scan(&each.id, &each.name, &each.email)

		if err != nil {
			panic(err.Error())
		}

		result = append(result, each)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}
  var names string
  
	for _, each := range result {
	}
}

func GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	fmt.Fprint(w, name)
}

type User struct {
	id    uint16
	name  string
	email string
}

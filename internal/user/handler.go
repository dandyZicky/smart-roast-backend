package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, db *sql.DB) {
	rows, err := db.Query("SELECT id, name, email FROM mock_user_v2")
	if err != nil {
		panic("Query failed")
	}
	defer rows.Close()

	var result []User

	for rows.Next() {
		var each = User{}
		var err = rows.Scan(&each.ID, &each.Name, &each.Email)

		if err != nil {
			panic(err.Error())
		}

		result = append(result, each)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "%s", string(jsonResult))
}

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params, db *sql.DB) {
	var u RegisteredUser

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	row := db.QueryRow("SELECT email FROM mock_user_v2 WHERE email = $1", u.Email)

	var email string

	err = row.Scan(&email)

	if err == nil {
		http.Error(w, "Email already exist", 400)
		return
	}

	// Temporary
	u.Password = "h986h8679#$#F@$vyrtuvV$"
	u.Salt = "C$^^V$7gy645y6f44#Y"

	jsonResult, err := json.Marshal(u)

	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	insertQuery := `
  INSERT INTO mock_user_v2 (name, email, password, salt) 
  VALUES ($1, $2, $3, $4)`

	_, e := db.Exec(insertQuery, u.Name, u.Email, u.Password, u.Salt)

	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}

	fmt.Fprintf(w, string(jsonResult))
}

type User struct {
	ID    uint16 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RegisteredUser struct {
	User
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

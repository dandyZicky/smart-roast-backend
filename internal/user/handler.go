package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

func GetUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, db *sql.DB) {
	rows, err := db.Query(`SELECT name, email FROM users`)
	if err != nil {
		panic("Query failed")
	}
	defer rows.Close()

	var result []UserSafe

	for rows.Next() {
		each := UserSafe{}
		err := rows.Scan(&each.Name, &each.Email)
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
	var u User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if !isEmailValid(&u.Email) {
		http.Error(w, "Invalid email form", http.StatusBadRequest)
		return
	}

	row := db.QueryRow(`SELECT email FROM users WHERE email = $1`, u.Email)

	err = row.Scan(&u.Email)

	if err == nil {
		http.Error(w, "Bad email", 400)
		return
	}

	// Temporary
	u.Salt = "C$^^V$7gy645y6f44#Y"
	u.Password += u.Salt

	insertQuery := `INSERT INTO users (name, email, password, salt) VALUES ($1, $2, $3, $4) RETURNING id`

	row = db.QueryRow(insertQuery, u.Name, u.Email, u.Password, u.Salt)

	err = row.Scan(&u.ID)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func isEmailValid(s *string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(*s)
}

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type UserSafe struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

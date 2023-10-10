package routes

import (
	"database/sql"
	"net/http"

	"github.com/dandyZicky/v2-project/internal/user"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(db *sql.DB) *httprouter.Router {
	r := httprouter.New()

	r.GET("/user", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		user.GetUsers(w, r, p, db)
	})

	r.GET("/user/:name", user.GetUser)
	return r
}

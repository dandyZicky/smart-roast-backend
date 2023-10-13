package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/dandyZicky/v2-project/internal/roast"
	"github.com/dandyZicky/v2-project/internal/user"
)

func NewRouter(db *sql.DB) *httprouter.Router {
	r := httprouter.New()

	r.GET("/user", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		user.GetUsers(w, r, p, db)
	})

	r.POST("/user", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		user.CreateUser(w, r, p, db)
	})

	r.GET("/roast/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		roast.NewRoastSession(p.ByName("id"), &w, db)
	})

	return r
}

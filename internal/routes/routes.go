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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")
		roast.NewRoastSession(p.ByName("id"), w, r, db)
	})

	return r
}

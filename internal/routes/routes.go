package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/smart-roast/backend/internal/roast"
	"github.com/smart-roast/backend/internal/user"
)

func NewRouter(db *sql.DB) *httprouter.Router {
	r := httprouter.New()

	r.GET("/user", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		user.GetUsers(w, r, p, db)
	})

	r.POST("/user", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		user.CreateUser(w, r, p, db)
	})

	r.GET("/roast", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		res, err := roast.GetRoastSessions(db)
		if err != nil {
			http.Error(w, "Query failed", 400)
		}
		fmt.Fprintln(w, res)
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

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

	r.GET(
		"/roast-measurements/:sessionId",
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			res, err := roast.GetMeasurements(db, p.ByName("sessionId"))
			if err != nil {
				http.Error(w, "Query failed", 400)
				return
			}
			fmt.Fprintln(w, res)
		},
	)

	r.GET("/roast/:id/sessions", func(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
		res, err := roast.GetRoastSessions(db, p.ByName("id"))
		if err != nil {
			http.Error(w, "Query failed", 400)
			return
		}
		fmt.Fprintln(w, res)
	})

	r.GET(
		"/roast/:id/new-session",
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			roasterId := r.URL.Query().Get("roaster_id")

			if roasterId == "" {
				http.Error(w, `{message: "invalid roaster id"}`, 400)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.Header().Set("Content-Type", "text/event-stream")
			roast.NewRoastSession(p.ByName("id"), roasterId, w, r, db)
		},
	)

	return r
}

package main

import (
	"fmt"
	"net/http"

	"github.com/smart-roast/backend/internal/db"
	"github.com/smart-roast/backend/internal/routes"
)

func main() {
	h := "127.0.0.1"
	p := 3000
	addr := fmt.Sprintf("%s:%d", h, p)

	cs := "user=postgres dbname=test password=;"
	db, err := db.Db(&cs)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	r := routes.NewRouter(db)

	fmt.Printf("Listening on address: %s\n", addr)
	e := http.ListenAndServe(addr, r)

	if e != nil {
		panic(err.Error())
	}
}

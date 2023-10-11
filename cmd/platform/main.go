package main

import (
	"fmt"
	"net/http"

	"github.com/dandyZicky/v2-project/internal/db"
	"github.com/dandyZicky/v2-project/internal/routes"
)

func main() {
	h := "127.0.0.1"
	p := 3000
	addr := fmt.Sprintf("%s:%d", h, p)

	cs := "user=postgres dbname=test password=;"

	db, err := db.Db(&cs)
	defer db.Close()

	if err != nil {
		panic(err.Error())
	}

	r := routes.NewRouter(db)

	fmt.Printf("Listening on address: %s\n", addr)
	e := http.ListenAndServe(addr, r)

	if e != nil {
		panic(fmt.Sprintf("Invalid address: %s", e))
	}
}

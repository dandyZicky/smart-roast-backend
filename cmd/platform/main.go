package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/dandyZicky/v2-project/internal/routes"
	_ "github.com/lib/pq"
)

func main() {
	h := "127.0.0.1"
	p := 3000
	addr := fmt.Sprintf("%s:%d", h, p)

	connStr := "user=postgres dbname=test password=;"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Panicf("DB ERR: %s", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	r := routes.NewRouter(db)

	fmt.Printf("Listening on address: %s\n", addr)
	e := http.ListenAndServe(addr, r)

	if e != nil {
		panic(fmt.Sprintf("Invalid address: %s", e))
	}
}

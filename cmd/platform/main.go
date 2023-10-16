package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/smart-roast/backend/internal/db"
	"github.com/smart-roast/backend/internal/routes"
)

func main() {
	h := "127.0.0.1"
	p := 3000
	addr := fmt.Sprintf("%s:%d", h, p)

	log.Println("Connecting to database...")
	cs := "user=postgres dbname=test password=;"
	db, err := db.Db(&cs)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	log.Println("DB connection is set")

	r := routes.NewRouter(db)

	log.Printf("Listening on address: %s\n", addr)
	e := http.ListenAndServe(addr, r)

	if e != nil {
		panic(err.Error())
	}
}

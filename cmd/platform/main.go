package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/smart-roast/backend/internal/db"
	"github.com/smart-roast/backend/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	h := os.Getenv("SERVER_HOST_URL")
	p := os.Getenv("SERVER_HOST_PORT")
	addr := fmt.Sprintf("%s:%s", h, p)

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	log.Println("Connecting to database...")
	var cs string

	if os.Getenv("PROD") == "true" {
		cs = os.Getenv("CONN_STR")
	} else {
		cs = fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s",
			dbHost,
			dbPort,
			dbUser,
			dbName,
			dbPass,
		)
	}

	db, err := db.Db(&cs, "postgres")
	if err != nil {
		log.Fatalf("DBERR => %s", err.Error())
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

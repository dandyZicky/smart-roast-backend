// Deprecated test

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/smart-roast/backend/internal/db"
	"github.com/smart-roast/backend/internal/routes"
)

func main() {
	h := os.Getenv("SERVER_HOST_URL")
	p := os.Getenv("SERVER_HOST_PORT")
	addr := fmt.Sprintf("%s:%s", h, p)

	if h == "" {
		log.Println("No environment variables found, trying to read .env file")
		err := godotenv.Load()
		if err != nil {
			log.Fatal(".env file not found")
		}
		h = os.Getenv("SERVER_HOST_URL")
		p = os.Getenv("SERVER_HOST_PORT")
		addr = fmt.Sprintf("%s:%s", h, p)
	}

	log.Println("Connecting to database...")
	var cs string

	if os.Getenv("PROD") == "true" {
		log.Println("DB: PROD")
		cs = os.Getenv("CONN_STR")
	} else {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")

		log.Println("DB: DEV")
		cs = fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s",
			dbHost,
			dbPort,
			dbUser,
			dbName,
			dbPass,
		)
	}

	redisUrl := os.Getenv("REDIS_URL")

	db, err := db.Db(&cs, "postgres")
	if err != nil {
		log.Fatalf("DBERR => %s", err.Error())
	}
	defer db.Close()

	log.Println("DB connection is set")

	ctx := context.Background()
	opt, _ := redis.ParseURL(
		redisUrl,
	)
	client := redis.NewClient(opt)
	defer client.Close()
	log.Println("Redis connected")

	r := routes.NewRouter(db, client, &ctx)

	log.Printf("Listening on address: %s\n", addr)
	e := http.ListenAndServe(addr, r)

	if e != nil {
		panic(err.Error())
	}
}

package main

import (
	"database/sql"
	"github.com/halfdb/herro-world/internal/app/server"
	"log"
	"os"
)

func main() {
	errLogger := log.New(os.Stderr, "", 0)

	databaseUrl := os.Getenv("DB_STRING")
	db, err := sql.Open("mysql", databaseUrl)
	if err != nil {
		errLogger.Fatal(err)
	}
	port := os.Getenv("PORT")

	errLogger.Fatal(server.New(":"+port, db))
}

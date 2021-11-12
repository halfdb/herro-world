package main

import (
	"context"
	"database/sql"
	"github.com/halfdb/herro-world/internal/app/server"
	"log"
	"os"
)

func main() {
	errLogger := log.New(os.Stderr, "", 0)

	serverCtx := context.Background()
	databaseUrl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("mysql", databaseUrl)
	if err != nil {
		errLogger.Fatal(err)
	}
	port := os.Getenv("PORT")

	errLogger.Fatal(server.New(":"+port, serverCtx, db))
}

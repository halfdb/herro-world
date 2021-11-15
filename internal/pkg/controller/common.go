package controller

import (
	"database/sql"
)

var (
	db *sql.DB
)

func InitializeController(d *sql.DB) {
	db = d
}

func checkInitialized() {
	if db == nil {
		panic("controllers not initialized")
	}
}

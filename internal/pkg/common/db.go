package common

import (
	"context"
	"database/sql"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func GetDB() *sql.DB {
	return boil.GetDB().(*sql.DB)
}

func BeginTx() (*sql.Tx, error) {
	return GetDB().BeginTx(context.Background(), nil)
}

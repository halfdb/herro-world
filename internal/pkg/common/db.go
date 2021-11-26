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

type TxFn func(*sql.Tx) error

func DoInTx(fn TxFn) error {
	tx, err := BeginTx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return fn(tx)
}

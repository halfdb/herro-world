package dao

import (
	"database/sql"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/labstack/echo/v4"
)

func UpdateUser(uid int, updates models.M) error {
	tx, err := common.BeginTx()
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
	rowsAff, err := models.Users(models.UserWhere.UID.EQ(uid)).UpdateAll(tx, updates)
	if rowsAff == 0 {
		return sql.ErrNoRows
	} else if rowsAff != 1 || err != nil {
		return echo.ErrInternalServerError
	}
	return nil
}

func FetchUser(uid int) (*models.User, error) {
	return models.Users(models.UserWhere.UID.EQ(uid)).One(common.GetDB())
}

func FetchUsers(uids ...int) (models.UserSlice, error) {
	return models.Users(models.UserWhere.UID.IN(uids)).All(common.GetDB())
}

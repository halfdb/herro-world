package dao

import (
	"database/sql"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/labstack/echo/v4"
)

func UpdateUser(uid int, updates models.M) error {
	return common.DoInTx(func(tx *sql.Tx) error {
		rowsAff, err := models.Users(models.UserWhere.UID.EQ(uid)).UpdateAll(tx, updates)
		if rowsAff == 0 {
			return sql.ErrNoRows
		} else if rowsAff != 1 || err != nil {
			return echo.ErrInternalServerError
		}
		return nil
	})
}

func FetchUser(uid int) (*models.User, error) {
	return models.Users(models.UserWhere.UID.EQ(uid)).One(common.GetDB())
}

func FetchUsers(uids ...int) (models.UserSlice, error) {
	return models.Users(models.UserWhere.UID.IN(uids)).All(common.GetDB())
}

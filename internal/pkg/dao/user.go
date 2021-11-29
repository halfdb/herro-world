package dao

import (
	"database/sql"
	"errors"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strconv"
)

func UpdateUser(user *models.User) error {
	return common.DoInTx(func(tx *sql.Tx) error {
		rowsAff, err := user.Update(tx, boil.Infer())
		if rowsAff == 0 {
			return sql.ErrNoRows
		} else if err != nil {
			return err
		} else if rowsAff != 1 {
			return errors.New("unexpected: rowsAff = " + strconv.FormatInt(rowsAff, 10))
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

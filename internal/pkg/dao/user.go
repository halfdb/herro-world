package dao

import (
	"database/sql"
	"errors"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strconv"
)

// LookupUser finds the user with specified login name and password, raises error if more than 1 is found
func LookupUser(loginName, password string) (*models.User, error) {
	users, err := models.Users(models.UserWhere.LoginName.EQ(loginName), models.UserWhere.Password.EQ(password)).All(common.GetDB())
	if err != nil {
		return nil, err
	} else if len(users) > 1 {
		return nil, errors.New("more than 1 user found")
	} else {
		return users[0], nil
	}
}

func CreateUser(tx *sql.Tx, user *models.User) (*models.User, error) {
	err := user.Insert(tx, boil.Infer())
	return user, err
}

func ExistUser(loginName string) (bool, error) {
	count, err := models.Users(models.UserWhere.LoginName.EQ(loginName)).Count(common.GetDB())
	return count > 0, err
}

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

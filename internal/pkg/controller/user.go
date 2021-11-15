package controller

import (
	"context"
	"github.com/halfdb/herro-world/internal/pkg/auth"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type UserBind struct {
	Uid           int    `boil:"uid" json:"uid"`
	LoginName     string `boil:"login_name" json:"login_name,omitempty"`
	Nickname      string `boil:"nickname" json:"nickname"`
	ShowLoginName bool   `boil:"show_login_name" json:"show_login_name"`
}

func extractUid(c echo.Context, key string) (int, error) {
	uidString := c.Param(key)
	return strconv.Atoi(uidString)
}

func FetchUser(uid int, bind interface{}) error {
	checkInitialized()
	return models.Users(models.UserWhere.UID.EQ(uid)).Bind(nil, db, bind)
}

func UpdateUser(uid int, updates models.M) error {
	checkInitialized()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	rowsAffected, err := models.Users(models.UserWhere.UID.EQ(uid)).UpdateAll(tx, updates)
	if rowsAffected > 1 || err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return echo.ErrInternalServerError
	}
	return tx.Commit()
}

// GetUserInfo is open to public
func GetUserInfo(c echo.Context) error {
	var user UserBind
	uid, err := extractUid(c, "uid")
	if err != nil {
		c.Logger().Error("failed to extract uid")
		return err
	}
	if err := FetchUser(uid, &user); err != nil {
		c.Logger().Error("failed to fetch user")
		return err
	}
	if auth.GetUid(c) != user.Uid && !user.ShowLoginName {
		user.LoginName = ""
	}
	return c.JSON(http.StatusOK, user)
}

// PatchUserInfo asserts user is authorized, so the uid in token is same with that in query params
func PatchUserInfo(c echo.Context) error {
	uid := auth.GetUid(c)

	values := c.QueryParams()
	updates := make(models.M)
	for _, param := range []string{"nickname", "show_login_name", "password"} {
		value := values.Get(param)
		if value != "" {
			updates[param] = value
		}
	}
	if len(updates) == 0 {
		c.Logger().Error("no parameters provided")
		return echo.ErrBadRequest
	}
	if values.Get("show_login_name") != "" { // implies that updates["show_login_name"] must be set
		var err error
		updates["show_login_name"], err = strconv.ParseBool(updates["show_login_name"].(string))
		if err != nil {
			c.Logger().Error("invalid bool string")
			return echo.ErrBadRequest
		}
	}

	if err := UpdateUser(uid, updates); err != nil {
		c.Logger().Error("failed to update user")
		return err
	}

	var user UserBind
	if err := FetchUser(uid, &user); err != nil {
		c.Logger().Error("failed to fetch user")
		return err
	}
	return c.JSON(http.StatusOK, user)
}

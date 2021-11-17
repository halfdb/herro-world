package controller

import (
	"github.com/halfdb/herro-world/internal/pkg/auth"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/halfdb/herro-world/pkg/dto"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	keyUid           = "uid"
	keyLoginName     = "login_name"
	keyNickname      = "nickname"
	keyShowLoginName = "show_login_name"
	keyPassword      = "password"
)

func extractUid(c echo.Context, key string) (int, error) {
	uidString := c.Param(key)
	return strconv.Atoi(uidString)
}

func GetUserInfo(c echo.Context) error {
	var user dto.User
	uid, err := extractUid(c, "uid")
	if err != nil {
		c.Logger().Error("failed to extract uid")
		return err
	}
	if err := dao.FetchUser(uid, &user); err != nil {
		c.Logger().Error("failed to fetch user")
		return err
	}
	// hide login name
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
	for _, param := range []string{keyNickname, keyShowLoginName, keyPassword} {
		value := values.Get(param)
		if value != "" {
			updates[param] = value
		}
	}
	if len(updates) == 0 {
		c.Logger().Error("no parameters provided")
		return echo.ErrBadRequest
	}
	if values.Get(keyShowLoginName) != "" { // implies that updates[keyShowLoginName] must be set
		var err error
		updates[keyShowLoginName], err = strconv.ParseBool(updates[keyShowLoginName].(string))
		if err != nil {
			c.Logger().Error("invalid bool string")
			return echo.ErrBadRequest
		}
	}

	if err := dao.UpdateUser(uid, updates); err != nil {
		c.Logger().Error("failed to update user")
		return err
	}

	var user dto.User
	if err := dao.FetchUser(uid, &user); err != nil {
		c.Logger().Error("failed to fetch user")
		return err
	}
	return c.JSON(http.StatusOK, user)
}

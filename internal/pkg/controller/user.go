package controller

import (
	"database/sql"
	"github.com/halfdb/herro-world/internal/pkg/authorization"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/halfdb/herro-world/pkg/dto"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	keyUid           = "uid"
	keyNickname      = "nickname"
	keyShowLoginName = "show_login_name"
	keyPassword      = "password"
)

// TODO refactor

func convertUser(user *models.User) *dto.User {
	result := &dto.User{
		Uid:           user.UID,
		LoginName:     user.LoginName,
		ShowLoginName: user.ShowLoginName,
	}
	if user.Nickname.Valid {
		result.Nickname = user.Nickname.String
	}
	return result
}

func GetUserInfo(c echo.Context) error {
	uid, err := parsePathInt(c, "uid")
	if err != nil {
		c.Logger().Error("failed to extract uid")
		return err
	}

	user, err := dao.FetchUser(uid)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	} else if err != nil {
		c.Logger().Error("failed to fetch user")
		return err
	}
	// hide login name
	if authorization.GetUid(c) != user.UID && !user.ShowLoginName {
		user.LoginName = ""
	}
	return c.JSON(http.StatusOK, convertUser(user))
}

// PatchUserInfo asserts user is authorized, so the uid in token is same with that in query params
func PatchUserInfo(c echo.Context) error {
	uid := authorization.GetUid(c)

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
		if err != sql.ErrNoRows {
			c.Logger().Error("failed to update user")
			return err
		}
	}

	user, err := dao.FetchUser(uid)
	if err != nil {
		c.Logger().Error("failed to fetch user")
		return err
	}
	return c.JSON(http.StatusOK, convertUser(user))
}

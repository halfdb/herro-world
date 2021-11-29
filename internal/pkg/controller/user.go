package controller

import (
	"database/sql"
	"github.com/halfdb/herro-world/internal/pkg/authorization"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/halfdb/herro-world/pkg/dto"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
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

	user, err := dao.FetchUser(uid)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	} else if err != nil {
		c.Logger().Error("failed to fetch user")
		return err
	}

	values := c.QueryParams()
	if len(values) == 0 {
		c.Logger().Error("no parameters provided")
		return echo.ErrBadRequest
	}
	updates := 0
	for key, strings := range values {
		value := strings[0]
		if value == "" {
			continue
		}
		switch key {
		case keyNickname:
			if user.Nickname.Valid && user.Nickname.String != value {
				user.Nickname = null.StringFrom(value)
				updates += 1
			}
		case keyShowLoginName:
			result, err := strconv.ParseBool(value)
			if err != nil {
				return echo.ErrBadRequest
			}
			if user.ShowLoginName != result {
				user.ShowLoginName = result
				updates += 1
			}
		case keyPassword:
			if user.Password != value {
				user.Password = value
				updates += 1
			}
		default:
			c.Logger().Error("unrecognised query param: " + key)
			return echo.ErrBadRequest
		}
	}

	if updates > 0 {
		if err := dao.UpdateUser(user); err != nil {
			c.Logger().Error("failed to update user")
			return err
		}
	}
	return c.JSON(http.StatusOK, convertUser(user))
}

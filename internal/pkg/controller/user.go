package controller

import (
	"bytes"
	"encoding/base64"
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
	keyPublicKey     = "public_key"
	keyQuery         = "query"
	keyByNickname    = "by_nickname"
	keyByLoginName   = "by_login_name"
)

func convertUser(user *models.User, respectShowLoginName bool) *dto.User {
	result := &dto.User{
		Uid:           user.UID,
		ShowLoginName: user.ShowLoginName,
	}
	if user.ShowLoginName || !respectShowLoginName {
		result.LoginName = user.LoginName
	}
	if user.Nickname.Valid {
		result.Nickname = user.Nickname.String
	}
	if user.PublicKey.Valid {
		result.PublicKey = base64.StdEncoding.EncodeToString(user.PublicKey.Bytes)
	}
	return result
}

func GetUserInfo(c echo.Context) error {
	var uid int
	err := echo.PathParamsBinder(c).Int(keyUid, &uid).BindError()
	if err != nil {
		c.Logger().Error("failed to extract uid")
		return err
	}

	user, err := dao.FetchUser(uid)
	if err != nil {
		c.Logger().Error("failed to fetch user")
		return err
	}
	if user == nil {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, convertUser(user, authorization.GetUid(c) != user.UID))
}

// PatchUserInfo asserts user is authorized, so the uid in token is same with that in query params
func PatchUserInfo(c echo.Context) error {
	uid := authorization.GetUid(c)

	user, err := dao.FetchUser(uid)
	if err != nil {
		c.Logger().Error("failed to fetch user")
		return err
	}
	if user == nil {
		return echo.ErrNotFound
	}

	values := c.QueryParams()
	if len(values) == 0 {
		c.Logger().Error("no parameters provided")
		return echo.ErrBadRequest
	}
	update := false
	for key, strings := range values {
		// /api?key=value1&key=value2
		// strings = [value1, value2]
		value := strings[0]
		if value == "" {
			continue
		}
		switch key {
		case keyNickname:
			if !user.Nickname.Valid || user.Nickname.String != value {
				user.Nickname = null.StringFrom(value)
				update = true
			}
		case keyShowLoginName:
			result, err := strconv.ParseBool(value)
			if err != nil {
				return echo.ErrBadRequest
			}
			if user.ShowLoginName != result {
				user.ShowLoginName = result
				update = true
			}
		case keyPassword:
			if user.Password != value {
				user.Password = value
				update = true
			}
		case keyPublicKey:
			key, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				return err
			}
			if len(key) > 300 {
				return echo.ErrBadRequest
			}
			if !user.PublicKey.Valid || bytes.Compare(user.PublicKey.Bytes, key) != 0 {
				user.PublicKey = null.BytesFrom(key)
				update = true
			}
		default:
			c.Logger().Error("unrecognised query param: " + key)
			return echo.ErrBadRequest
		}
	}

	if update {
		if err := dao.UpdateUser(user); err != nil {
			c.Logger().Error("failed to update user")
			return err
		}
	}
	return c.JSON(http.StatusOK, convertUser(user, false))
}

func SearchUser(c echo.Context) error {
	uid := authorization.GetUid(c)
	query := ""
	byNickname := true
	byLoginName := true
	err := echo.QueryParamsBinder(c).
		String(keyQuery, &query).
		Bool(keyByNickname, &byNickname).
		Bool(keyByLoginName, &byLoginName).
		BindError()
	c.Logger().Debug(keyQuery + query)
	c.Logger().Debug(keyByNickname + strconv.FormatBool(byNickname))
	c.Logger().Debug(keyByLoginName + strconv.FormatBool(byLoginName))
	if err != nil {
		return err
	}
	results := make([]*dto.User, 0)
	if byLoginName {
		user, err := dao.LookupUserLoginName(query, false)
		if err != nil {
			return err
		}
		if user != nil {
			results = append(results, convertUser(user, uid != user.UID))
		}
	}
	if byNickname {
		users, err := dao.LookupUserNickname(query)
		if err != nil {
			return err
		}
		dtoUsers := make([]*dto.User, len(users))
		for i, user := range users {
			dtoUsers[i] = convertUser(user, uid != user.UID)
		}
		results = append(results, dtoUsers...)
	}

	return c.JSON(http.StatusOK, results)
}

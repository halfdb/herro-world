package authorization

import (
	"github.com/golang-jwt/jwt"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/labstack/echo/v4"
	"strconv"
)

const (
	keyAuthenticated = "authenticated"
	keyUid           = "uid"
)

func Skipper(e echo.Context) bool {
	return (e.Path() == "/login" || e.Path() == "/users") && e.Request().Method == "POST"
}

func SetAuthedContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Debug("start setting authed context")
		if Skipper(c) {
			c.Logger().Debug("skipped")
			return next(c)
		}
		if token := c.Get("user"); token != nil {
			c.Logger().Debug("authed")
			c.Set(keyAuthenticated, true)
			claims := token.(*jwt.Token).Claims.(*common.Claims)
			c.Set(keyUid, claims.Uid)
		} else {
			c.Logger().Debug("not authed")
			c.Set(keyAuthenticated, false)
		}
		return next(c)
	}
}

// GetUid asserts that the request is authed, returns the uid related to the request
func GetUid(c echo.Context) int {
	if authed := c.Get(keyAuthenticated); authed == nil || !authed.(bool) {
		panic("request not authed")
	}
	return c.Get(keyUid).(int)
}

// AuthorizeSelf asserts the request is authenticated, and checks if the user is authorized for the request
func AuthorizeSelf(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Debug("check if user is allowed for the path")
		queryUid := c.Param(keyUid)
		tokenUid := strconv.Itoa(GetUid(c))
		if queryUid != tokenUid {
			c.Logger().Debug("no, not allowed")
			return echo.ErrForbidden
		}
		return next(c)
	}
}

package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const (
	keyAuthenticated = "authenticated"
	keyUid           = "uid"
)

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
			claims := token.(*jwt.Token).Claims.(*Claims)
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

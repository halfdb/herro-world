package controller

import (
	"github.com/halfdb/herro-world/internal/pkg/authorization"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func Herro(c echo.Context) error {
	return c.String(http.StatusOK, "Herro, "+strconv.Itoa(authorization.GetUid(c)))
}

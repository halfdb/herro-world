package controller

import (
	"github.com/halfdb/herro-world/internal/pkg/auth"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func Herro(c echo.Context) error {
	return c.String(http.StatusOK, "Herro, "+strconv.Itoa(auth.GetUid(c)))
}

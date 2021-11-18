package controller

import "github.com/labstack/echo/v4"

func parsePathInt(c echo.Context, key string) (int, error) {
	i := 0
	err := echo.PathParamsBinder(c).Int(key, &i).BindError()
	if err != nil || i == 0 {
		return 0, echo.ErrBadRequest
	}
	return i, nil
}

package auth

import (
	"context"
	"database/sql"
	"github.com/golang-jwt/jwt"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"os"
	"time"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func GetJWTSecret() string {
	return jwtSecret
}

type Claims struct {
	Uid int
	jwt.StandardClaims
}

func Validator(ctx context.Context, db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Debug("start validating")
		loginName := c.QueryParam("login_name")
		password := c.QueryParam("password")
		c.Logger().Debug(loginName, password)
		users, err := models.Users(qm.Where("login_name=? and password=?", loginName, password)).All(ctx, db)
		if err != nil {
			return err
		}
		if len(users) != 1 {
			return echo.ErrUnauthorized
		}

		user := users[0]
		claims := Claims{
			Uid: user.UID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": signedToken,
		})
	}
}

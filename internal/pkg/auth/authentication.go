package auth

import (
	"database/sql"
	"github.com/golang-jwt/jwt"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

func signUser(user *models.User) (string, error) {
	claims := Claims{
		Uid: user.UID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(GetJWTSecret()))
}

func Validator(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Debug("start validating")
		loginName := c.QueryParam("login_name")
		password := c.QueryParam("password")
		c.Logger().Debug(loginName, password)
		users, err := models.Users(qm.Where("login_name=? and password=?", loginName, password)).All(db)
		if err != nil {
			return err
		}
		if len(users) != 1 {
			c.Logger().Debug("len(users) != 1")
			return echo.ErrUnauthorized
		}

		signedToken, err := signUser(users[0])
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": signedToken,
			"uid":   users[0].UID,
		})
	}
}

func Register(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := &models.User{
			LoginName: c.QueryParam("login_name"),
			Password:  c.QueryParam("password"),
		}
		if nickname := c.QueryParam("nickname"); nickname != "" {
			user.Nickname = null.StringFrom(nickname)
		}
		if err := user.Insert(db, boil.Infer()); err != nil {
			return err
		}
		signedToken, err := signUser(user)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": signedToken,
			"uid":   user.UID,
		})
	}
}

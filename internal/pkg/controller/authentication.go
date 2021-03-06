package controller

import (
	"database/sql"
	"github.com/golang-jwt/jwt"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"net/http"
	"os"
	"time"
)

const (
	keyLoginName = "login_name"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func GetJWTSecret() string {
	return jwtSecret
}

func signUser(user *models.User) (string, error) {
	claims := common.Claims{
		Uid: user.UID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(GetJWTSecret()))
}

func validateLoginName(loginName string) bool {
	return loginName != ""
}

func validatePassword(password string) bool {
	return password != ""
}

func Login(c echo.Context) error {
	loginName := c.FormValue(keyLoginName)
	password := c.FormValue(keyPassword)
	if !validateLoginName(loginName) || !validatePassword(password) {
		return echo.ErrBadRequest
	}
	user, err := dao.LookupUser(loginName, password)
	if err != nil {
		return err
	}
	if user == nil {
		return echo.ErrUnauthorized
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

func Register(c echo.Context) error {
	loginName := c.FormValue(keyLoginName)
	password := c.FormValue(keyPassword)
	nickname := c.FormValue(keyNickname)
	if !validateLoginName(loginName) || !validatePassword(password) {
		return echo.ErrBadRequest
	}

	exists, err := dao.ExistUserLoginName(loginName)
	if err != nil {
		return err
	}
	if exists {
		c.Logger().Info("login name already exists")
		return c.String(http.StatusConflict, "Login name conflict")
	}

	user := &models.User{
		LoginName: loginName,
		Password:  password,
		Nickname:  null.NewString(nickname, nickname != ""),
	}

	err = common.DoInTx(func(tx *sql.Tx) error {
		user, err = dao.CreateUser(tx, user)
		return err
	})
	if err != nil {
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

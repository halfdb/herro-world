package server

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/halfdb/herro-world/internal/pkg/auth"
	"github.com/halfdb/herro-world/internal/pkg/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func New(address string, db *sql.DB) error {
	controller.InitializeController(db)

	// Echo instance
	e := echo.New()
	e.Debug = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper:    auth.Skipper,
		SigningKey: []byte(auth.GetJWTSecret()),
		Claims:     &auth.Claims{},
	}))
	e.Use(auth.SetAuthedContext)

	privateGroup := e.Group("/users/:uid/*")
	privateGroup.Use(auth.AuthorizeSelf(middleware.DefaultSkipper))

	// Routes
	e.GET("/", controller.Herro)
	e.POST("/users", auth.Register(db))
	e.POST("/login", auth.Validator(db))
	e.GET("/users/:uid", controller.GetUserInfo)
	e.Add(http.MethodPatch, "/users/:uid", controller.PatchUserInfo, auth.AuthorizeSelf(middleware.DefaultSkipper))

	// Start server
	return e.Start(address)
}

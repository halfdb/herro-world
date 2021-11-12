package server

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/halfdb/herro-world/internal/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func New(address string, ctx context.Context, db *sql.DB) error {
	// Echo instance
	e := echo.New()
	e.Debug = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper: func(e echo.Context) bool {
			return (e.Path() == "/login" || e.Path() == "/users") && e.Request().Method == "POST"
		},
		SigningKey: []byte(auth.GetJWTSecret()),
	}))

	// Routes
	e.GET("/", hello)
	e.POST("/login", auth.Validator(ctx, db))

	// Start server
	return e.Start(address)
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Herro, World!")
}

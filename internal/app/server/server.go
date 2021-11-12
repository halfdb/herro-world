package server

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/halfdb/herro-world/internal/pkg/auth"
	"github.com/halfdb/herro-world/internal/pkg/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(address string, ctx context.Context, db *sql.DB) error {
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

	// Routes
	e.GET("/", controller.Herro)
	e.POST("/users", auth.Register(ctx, db))
	e.POST("/login", auth.Validator(ctx, db))

	// Start server
	return e.Start(address)
}

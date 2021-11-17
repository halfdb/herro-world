package server

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/halfdb/herro-world/internal/pkg/auth"
	"github.com/halfdb/herro-world/internal/pkg/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func New(address string, db *sql.DB) error {
	boil.SetDB(db)
	boil.DebugMode = true

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
	e.POST("/users", auth.Register(db))
	e.POST("/login", auth.Validator(db))
	e.GET("/users/:uid", controller.GetUserInfo)
	e.PATCH("/users/:uid", controller.PatchUserInfo, auth.AuthorizeSelf)
	e.POST("/users/:uid/contacts", controller.PostContacts, auth.AuthorizeSelf)
	e.GET("/users/:uid/contacts", controller.GetContacts, auth.AuthorizeSelf)
	e.PATCH("/users/:uid/contacts/:uid_other", controller.PatchContact, auth.AuthorizeSelf)
	e.DELETE("/users/:uid/contacts/:uid_other", controller.DeleteContact, auth.AuthorizeSelf)
	e.GET("/users/:uid/chats", controller.GetChats, auth.AuthorizeSelf)
	// TODO API to get chat info
	e.GET("/chats/:cid/messages", controller.GetMessages, auth.AuthorizeChatMember)
	e.POST("/chats/:cid/messages", controller.PostMessage, auth.AuthorizeChatMember)

	// Start server
	return e.Start(address)
}

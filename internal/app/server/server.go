package server

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/halfdb/herro-world/internal/pkg/authorization"
	"github.com/halfdb/herro-world/internal/pkg/common"
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
		Skipper:    authorization.Skipper,
		SigningKey: []byte(controller.GetJWTSecret()),
		Claims:     &common.Claims{},
	}))
	e.Use(authorization.SetAuthedContext)

	// Routes
	// auth
	e.GET("/", controller.Herro)
	e.POST("/users", controller.Register)
	e.POST("/login", controller.Login)
	// user info
	e.GET("/users/:uid", controller.GetUserInfo)
	e.PATCH("/users/:uid", controller.PatchUserInfo, authorization.AuthorizeSelf)
	e.POST("/users/search", controller.SearchUser)
	// contact
	e.POST("/users/:uid/contacts", controller.PostContacts, authorization.AuthorizeSelf)
	e.GET("/users/:uid/contacts", controller.GetContacts, authorization.AuthorizeSelf)
	e.PATCH("/users/:uid/contacts/:uid_other", controller.PatchContact, authorization.AuthorizeSelf)
	e.DELETE("/users/:uid/contacts/:uid_other", controller.DeleteContact, authorization.AuthorizeSelf)
	// chat & message
	e.GET("/users/:uid/chats", controller.GetChats, authorization.AuthorizeSelf)
	e.GET("/chats/:cid/messages", controller.GetMessages, authorization.AuthorizeChatMember)
	e.POST("/chats/:cid/messages", controller.PostMessage, authorization.AuthorizeChatMember)
	// group chat
	e.POST("/chats", controller.PostChats)
	e.GET("/chats/:cid/members", controller.GetChatMembers, authorization.AuthorizeChatMember)
	e.POST("/chats/:cid/members", controller.PostChatMembers, authorization.AuthorizeChatMember)
	e.DELETE("/chats/:cid/members/:uid", controller.DeleteChatMember, authorization.AuthorizeSelf, authorization.AuthorizeChatMember)

	// Start server
	return e.Start(address)
}

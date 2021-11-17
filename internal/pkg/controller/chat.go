package controller

import (
	"github.com/halfdb/herro-world/internal/pkg/auth"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/halfdb/herro-world/pkg/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

func convertChats(chats models.ChatSlice) dto.Chats {
	result := make(dto.Chats, len(chats))
	for i, chat := range chats {
		result[i] = chat.Cid
	}
	return result
}

func GetChats(c echo.Context) error {
	uid := auth.GetUid(c)
	chats, err := dao.FetchAllChats(uid, false)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, convertChats(chats))
}

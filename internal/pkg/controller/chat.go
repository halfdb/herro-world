package controller

import (
	"errors"
	"github.com/halfdb/herro-world/internal/pkg/authorization"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/halfdb/herro-world/pkg/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

func makeChats(chats models.ChatSlice) ([]*dto.Chat, error) {
	cids := make([]int, len(chats))
	for i, chat := range chats {
		cids[i] = chat.Cid
	}
	uidsCh := make(chan map[int][]int, 1)
	go func() {
		uids, err := dao.GetUids(cids...)
		if err != nil {
			close(uidsCh)
			return
		}
		uidsCh <- uids
	}()

	result := make([]*dto.Chat, len(chats))
	for i, chat := range chats {
		result[i] = &dto.Chat{
			Cid:    chat.Cid,
			Direct: chat.Direct,
		}
		if chat.Name.Valid {
			result[i].Name = chat.Name.String
		}
	}
	uids, ok := <-uidsCh
	if !ok {
		return nil, errors.New("failed to get uids of chats")
	}
	for _, chat := range result {
		chat.Uids = uids[chat.Cid]
	}
	return result, nil
}

func GetChats(c echo.Context) error {
	uid := authorization.GetUid(c)
	chats, err := dao.FetchVisibleChats(uid)
	if err != nil {
		return err
	}

	chatsDtos, err := makeChats(chats)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, chatsDtos)
}

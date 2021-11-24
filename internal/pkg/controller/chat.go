package controller

import (
	"errors"
	"github.com/halfdb/herro-world/internal/pkg/authorization"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/halfdb/herro-world/pkg/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	keyUids = "uids"
	keyName = "name"
)

func makeChats(chats models.ChatSlice) ([]*dto.Chat, error) {
	cids := make([]int, len(chats))
	for i, chat := range chats {
		cids[i] = chat.Cid
	}
	uidsCh := make(chan map[int][]int, 1)
	go func() {
		uids, err := dao.GetUids(false, cids...)
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

func convertChat(chat *models.Chat, uids []int) *dto.Chat {
	result := &dto.Chat{
		Cid:    chat.Cid,
		Direct: chat.Direct,
		Uids:   uids,
	}
	if chat.Name.Valid {
		result.Name = chat.Name.String
	}
	return result
}

func PostChats(c echo.Context) error {
	uid := authorization.GetUid(c)
	uids := make([]int, 0)
	name := ""
	namePtr := &name
	echo.QueryParamsBinder(c).Ints(keyUids, &uids).String(keyName, namePtr)

	// sanity check
	uids = common.UniqueInt(uids)
	if len(uids) < 3 {
		return echo.ErrBadRequest
	}
	// check uids in contact
	hasSelf := false
	for _, uidOther := range uids {
		if uidOther == uid {
			hasSelf = true
			continue
		}
		exists, err := dao.ContactExists(uid, uidOther, false)
		if err != nil {
			return err
		} else if !exists {
			return echo.ErrForbidden
		}
	}
	if !hasSelf {
		return echo.ErrBadRequest
	}

	if *namePtr == "" {
		namePtr = nil
	}

	// begin tx
	tx, err := common.BeginTx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	chat, err := dao.CreateChat(tx, namePtr, false, uids...)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, convertChat(chat, uids))
}

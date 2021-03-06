package controller

import (
	"database/sql"
	"errors"
	"github.com/halfdb/herro-world/internal/pkg/authorization"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/halfdb/herro-world/pkg/dto"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
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
		directCids := make([]int, 0, len(chats))
		groupCids := make([]int, 0, len(chats))
		for _, chat := range chats {
			if chat.Direct {
				directCids = append(directCids, chat.Cid)
			} else {
				groupCids = append(groupCids, chat.Cid)
			}
		}
		directUids, err := dao.LookupMemberUids(true, directCids...)
		if err != nil {
			close(uidsCh)
			return
		}
		groupUids, err := dao.LookupMemberUids(false, groupCids...)
		// reuse directUids
		for cid, uids := range groupUids {
			directUids[cid] = uids
		}
		uidsCh <- directUids
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
	chats, err := dao.LookupAllChats(uid, false)
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
	err := echo.QueryParamsBinder(c).Ints(keyUids, &uids).String(keyName, &name).BindError()
	if err != nil {
		return err
	}

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
		exists, err := dao.ExistContact(uid, uidOther, false)
		if err != nil {
			return err
		} else if !exists {
			return echo.ErrForbidden
		}
	}
	if !hasSelf {
		return echo.ErrBadRequest
	}

	chat := &models.Chat{
		Direct: false,
		Name:   null.NewString(name, name != ""),
	}
	err = common.DoInTx(func(tx *sql.Tx) error {
		chat, err = dao.CreateChat(tx, chat, uids...)
		return err
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, convertChat(chat, uids))
}

func GetChatMembers(c echo.Context) error {
	uid := authorization.GetUid(c)
	cid := authorization.GetCid(c)
	uidsMap, err := dao.LookupMemberUids(false, cid)
	if err != nil {
		return err
	}
	uids := uidsMap[cid]

	userSlice, err := dao.FetchUsers(uids...)
	users := make([]*dto.User, len(userSlice))
	for i, user := range userSlice {
		users[i] = convertUser(user, user.UID != uid)
	}

	return c.JSON(http.StatusOK, users)
}

func PostChatMembers(c echo.Context) error {
	uid := authorization.GetUid(c)
	cid := authorization.GetCid(c)
	chat := authorization.GetChat(c)
	// only adding to groups is allowed
	if chat.Direct {
		return echo.ErrBadRequest
	}

	uids := make([]int, 0)
	err := echo.QueryParamsBinder(c).Ints(keyUids, &uids).BindError()
	if err != nil {
		return err
	}
	uids = common.UniqueInt(uids)

	for _, uidOther := range uids {
		if uidOther == uid {
			continue
		}
		exists, err := dao.ExistContact(uid, uidOther, false)
		if err != nil {
			return err
		} else if !exists {
			return echo.ErrForbidden
		}
	}

	err = common.DoInTx(func(tx *sql.Tx) error {
		for _, uidOther := range uids {
			userChat, err := dao.FetchUserChat(uidOther, cid, true)
			if err != nil {
				return err
			}
			if userChat == nil { // user_chat does not exist, create it
				userChat = &models.UserChat{
					UID: uid,
					Cid: cid,
				}
				userChat, err = dao.CreateUserChat(tx, userChat)
			} else if userChat.DeletedAt.Valid { // user_chat deleted, restore it
				err = dao.RestoreUserChat(tx, &models.UserChat{
					UID: uidOther,
					Cid: cid,
				})
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return GetChatMembers(c)
}

func DeleteChatMember(c echo.Context) error {
	uid := authorization.GetUid(c)
	cid := authorization.GetCid(c)

	err := common.DoInTx(func(tx *sql.Tx) error {
		return dao.DeleteUserChat(tx, &models.UserChat{
			UID: uid,
			Cid: cid,
		})
	})
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "")
}

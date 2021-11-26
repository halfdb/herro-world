package controller

import (
	"database/sql"
	"encoding/base64"
	"github.com/halfdb/herro-world/internal/pkg/authorization"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/halfdb/herro-world/pkg/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	keyMime             = "mime"
	keyContent          = "content"
	defaultMessageLimit = 100
	contentLengthLimit  = 200
)

func convertMessage(message *models.Message) *dto.Message {
	return &dto.Message{
		Mid:     message.Mid,
		Cid:     message.Cid,
		Uid:     message.UID,
		Mime:    message.Mime,
		Content: base64.StdEncoding.EncodeToString(message.Content),
	}
}

func convertMessages(messages models.MessageSlice) []*dto.Message {
	result := make([]*dto.Message, len(messages))
	for i, message := range messages {
		result[i] = convertMessage(message)
	}
	return result
}

func GetMessages(c echo.Context) error {
	cid := authorization.GetCid(c)

	messages, err := dao.FetchAllMessages(cid, defaultMessageLimit)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, convertMessages(messages))
}

func PostMessage(c echo.Context) error {
	// params
	uid := authorization.GetUid(c)
	cid := authorization.GetCid(c)
	mime := "text/plain"
	content := ""
	err := echo.QueryParamsBinder(c).String(keyMime, &mime).String(keyContent, &content).BindError()
	if err != nil {
		return echo.ErrBadRequest
	}
	contentBytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return echo.ErrBadRequest
	} else if len(contentBytes) > contentLengthLimit {
		return echo.ErrStatusRequestEntityTooLarge
	}

	// build message
	message := &models.Message{
		Cid:     cid,
		UID:     uid,
		Mime:    mime,
		Content: contentBytes,
	}

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

	// add reverse contact and check if blocked
	chat := authorization.GetChat(c)
	if chat.Direct { // only handle direct chats
		uidsMap, err := dao.GetMemberUids(true, cid)
		if err != nil {
			return err
		}
		uids := uidsMap[cid]
		var uidOther int
		if uid == uids[0] {
			uidOther = uids[1]
		} else {
			uidOther = uids[0]
		}

		reverseContact, err := dao.FetchContact(uidOther, uid, true)
		switch {
		case err == sql.ErrNoRows: // create new contact
			reverseContact = &models.Contact{
				UIDSelf:  uidOther,
				UIDOther: uid,
				Cid:      cid,
			}
			_, err = dao.CreateContact(tx, reverseContact, false)
			if err != nil {
				c.Logger().Error("failed to create reverse contact")
				return err
			}
		case err != nil: // unknown error
			return err
		case reverseContact.BlockedAt.Valid: // blocked by receiver, 403
			return echo.ErrForbidden
		case reverseContact.DeletedAt.Valid: // not blocked, restore
			_, err = dao.RestoreContact(tx, reverseContact)
			if err != nil {
				return err
			}
		}
	}

	// create message
	message, err = dao.CreateMessage(tx, message)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, convertMessage(message))
}

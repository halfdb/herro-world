package controller

import (
	"encoding/base64"
	"github.com/halfdb/herro-world/internal/pkg/auth"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/halfdb/herro-world/pkg/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	keyMime    = "mime"
	keyContent = "content"
)

func parseCid(c echo.Context) (int, error) {
	cid := 0
	err := echo.PathParamsBinder(c).Int(keyCid, &cid).BindError()
	if err != nil || cid == 0 {
		return 0, echo.ErrBadRequest
	}
	return cid, nil
}

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
	cid, err := parseCid(c)
	if err != nil {
		return err
	}

	messages, err := dao.FetchAllMessages(cid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, convertMessages(messages))
}

func PostMessage(c echo.Context) error {
	// params
	uid := auth.GetUid(c)
	cid, err := parseCid(c)
	if err != nil {
		return err
	}
	mime := "text/plain"
	content := ""
	err = echo.QueryParamsBinder(c).String(keyMime, &mime).String(keyContent, &content).BindError()
	if err != nil {
		return echo.ErrBadRequest
	}
	contentBytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return echo.ErrBadRequest
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

	// TODO reverse-adding contact
	// create message
	message, err = dao.CreateMessage(tx, message)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, convertMessage(message))
}

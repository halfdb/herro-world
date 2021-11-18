package controller

import (
	"database/sql"
	"github.com/halfdb/herro-world/internal/pkg/auth"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/dao"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/halfdb/herro-world/pkg/dto"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"net/http"
	"strconv"
	"time"
)

const (
	keyUidSelf     = "uid_self"
	keyUidOther    = "uid_other"
	keyDisplayName = "display_name"
	keyCid         = "cid"
	keyDeletedAt   = "deleted_at"
	keyBlocked     = "blocked"
	keyBlockedAt   = "blocked_at"
)

func PostContacts(c echo.Context) error {
	uid := auth.GetUid(c)
	// parse param
	uidOther := 0
	displayName := ""
	err := echo.QueryParamsBinder(c).
		Int(keyUid, &uidOther). // use `uid` to fetch `uid_other` as designed
		String(keyDisplayName, &displayName).
		BindError()
	if err != nil {
		c.Logger().Error("invalid parameter")
		return echo.ErrBadRequest
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

	// check if contact already exists
	// TODO try refactoring using upsert
	contact, err := dao.FetchContact(uid, uidOther, true)
	switch {
	case err == sql.ErrNoRows: // does not exist
		// create contact
		contact = &models.Contact{
			UIDSelf:     uid,
			UIDOther:    uidOther,
			DisplayName: null.NewString(displayName, displayName != ""),
		}

		contact, err = dao.CreateContact(tx, contact, true)
		if err != nil {
			c.Logger().Error("failed to create contact")
			c.Logger().Error(err)
			return echo.ErrInternalServerError
		}
	case err == nil && contact.DeletedAt.Valid: // contact has been deleted before
		// restore user_chat
		c.Logger().Debug("restoring user chat")
		err := dao.RestoreUserChat(tx, contact.UIDSelf, contact.Cid)
		if err != nil {
			return err
		}
		// restore contact
		c.Logger().Debug("restoring contact")
		contact.DisplayName = null.NewString(displayName, displayName != "")
		contact, err = dao.RestoreContact(tx, contact)
		if err != nil {
			return err
		}
	case err == nil && !contact.DeletedAt.Valid:
		return c.String(http.StatusConflict, "Already in contact.")
	default: // error
		c.Logger().Error("error while checking contact existence")
		return echo.ErrInternalServerError
	}
	// commit to prevent later errors from rolling back the tx
	err = tx.Commit()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, convertContact(contact))
}

func convertContact(original *models.Contact) *dto.Contact {
	converted := &dto.Contact{}
	converted.Uid = original.UIDOther
	if !original.DisplayName.IsZero() {
		converted.DisplayName = original.DisplayName.String
	}
	converted.Cid = original.Cid
	return converted
}

func GetContacts(c echo.Context) error {
	uid := auth.GetUid(c)
	boilContacts, err := dao.FetchAllContacts(uid, false, false)
	if err != nil {
		return err
	}
	contacts := make([]*dto.Contact, len(boilContacts))
	for i, contact := range boilContacts {
		contacts[i] = convertContact(contact)
	}

	return c.JSON(http.StatusOK, contacts)
}

func PatchContact(c echo.Context) error {
	// parse uids
	uidSelf := auth.GetUid(c)
	uidOther := 0
	err := echo.PathParamsBinder(c).Int(keyUidOther, &uidOther).BindError()
	if err != nil || uidOther == 0 {
		c.Logger().Error("invalid parameter")
		return echo.ErrBadRequest
	}

	// parse query params
	updates := models.M{}
	params := c.QueryParams()
	if len(params) == 0 {
		c.Logger().Error("invalid parameter")
		return echo.ErrBadRequest
	}
	for key, value := range params {
		switch key {
		case keyDisplayName:
			updates[keyDisplayName] = value[0]
		case keyBlocked:
			blocked, err := strconv.ParseBool(value[0])
			if err != nil || !blocked {
				c.Logger().Error("invalid parameter")
				return echo.ErrBadRequest
			}
			now := time.Now()
			updates[keyBlockedAt] = now
			updates[keyDeletedAt] = now
		}
	}

	// update
	err = dao.UpdateContact(uidSelf, uidOther, updates)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	} else if err != nil {
		c.Logger().Error("failed to update contact")
		return err
	}

	// fetch
	contact, err := dao.FetchContact(uidSelf, uidOther, true)
	if err != nil {
		c.Logger().Error("failed to fetch contact")
		return echo.ErrInternalServerError
	}
	if contact.DeletedAt.Valid {
		return c.String(http.StatusOK, "")
	}
	return c.JSON(http.StatusOK, convertContact(contact))
}

func DeleteContact(c echo.Context) error {
	// parse uids
	uidSelf := auth.GetUid(c)
	uidOther := 0
	err := echo.PathParamsBinder(c).Int(keyUidOther, &uidOther).BindError()
	if err != nil || uidOther == 0 {
		c.Logger().Error("invalid parameter")
		return echo.ErrBadRequest
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

	// find out cid first
	cid, err := dao.LookupDirectChat(uidSelf, uidOther, false)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	} else if err != nil {
		return err
	}
	// delete user_chat
	err = dao.DeleteUserChat(tx, uidSelf, cid)
	if err == sql.ErrNoRows {
		c.Logger().Error("unexpected: failed to delete user_chat")
		return echo.ErrInternalServerError
	} else if err != nil {
		return err
	}

	// delete contact
	err = dao.DeleteContact(tx, uidSelf, uidOther)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	} else if err != nil {
		c.Logger().Error("failed to delete contact")
		return err
	}

	// commit to prevent later errors from rolling back the tx
	err = tx.Commit()
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "")
}

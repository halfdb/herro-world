package controller

import (
	"database/sql"
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
	keyUidOther    = "uid_other"
	keyDisplayName = "display_name"
	keyBlocked     = "blocked"
)

func PostContacts(c echo.Context) error {
	uid := authorization.GetUid(c)
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

	// check if contact already exists
	contact, err := dao.FetchContact(uid, uidOther, true)
	if err != nil { // error
		c.Logger().Error("error while checking contact existence")
		return err
	}
	if contact != nil && !contact.DeletedAt.Valid { // already exists
		return c.String(http.StatusConflict, "Already in contact.")
	}

	err = common.DoInTx(func(tx *sql.Tx) error {
		if contact == nil { // never existed before, create
			c.Logger().Debug("not added before, create chat and contact")
			// get user nickname
			user, err := dao.FetchUser(uidOther)
			if err != nil {
				return err
			}
			if user == nil {
				c.Logger().Info("target user does not exist")
				return echo.ErrNotFound
			}
			if displayName == "" && user.Nickname.Valid {
				displayName = user.Nickname.String
			}
			// create contact
			contact = &models.Contact{
				UIDSelf:     uid,
				UIDOther:    uidOther,
				DisplayName: null.NewString(displayName, true),
			}
			contact.Cid, err = dao.LookupDirectChat(uidOther, uid, true)
			if err != nil {
				return err
			}
			contact, err = dao.CreateContact(tx, contact, contact.Cid == 0)

			if err != nil {
				c.Logger().Error("failed to create contact")
			}
			return err
		} else /* if contact.DeletedAt.Valid */ { // deleted before; restore
			// restore user_chat
			c.Logger().Debug("restoring user chat")
			err := dao.RestoreUserChat(tx, &models.UserChat{
				UID: contact.UIDSelf,
				Cid: contact.Cid,
			})
			if err != nil {
				c.Logger().Error("failed to restore user chat")
				c.Logger().Error(err)
				return err
			}
			// restore contact
			c.Logger().Debug("restoring contact")
			contact.DisplayName = null.NewString(displayName, displayName != "")
			contact, err = dao.RestoreContact(tx, contact)
			if err != nil {
				c.Logger().Error("failed to restore contact")
			}
			return err
		}
	})
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
	uid := authorization.GetUid(c)
	boilContacts, err := dao.LookupAllContacts(uid, false, false)
	c.Logger().Debug(boilContacts)
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
	uidSelf := authorization.GetUid(c)
	uidOther := 0
	err := echo.PathParamsBinder(c).Int(keyUidOther, &uidOther).BindError()
	if err != nil || uidOther == 0 {
		c.Logger().Error("invalid parameter")
		return echo.ErrBadRequest
	}

	// parse query params
	params := c.QueryParams()
	if !params.Has(keyDisplayName) {
		c.Logger().Error("invalid parameter")
		return echo.ErrBadRequest
	}
	displayName := params.Get(keyDisplayName)

	// fetch
	contact, err := dao.FetchContact(uidSelf, uidOther, false)
	if err != nil {
		c.Logger().Error("failed to fetch contact")
		return err
	}
	if contact == nil {
		return echo.ErrNotFound
	}

	if contact.DisplayName.Valid && contact.DisplayName.String != displayName {
		contact.DisplayName = null.StringFrom(displayName)
		err = common.DoInTx(func(tx *sql.Tx) error {
			// update
			err := dao.UpdateContact(tx, contact)
			if err != nil {
				c.Logger().Error("failed to update contact")
			}
			return err
		})
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, convertContact(contact))
}

func DeleteContact(c echo.Context) error {
	// parse uids
	uidSelf := authorization.GetUid(c)
	uidOther := 0
	err := echo.PathParamsBinder(c).Int(keyUidOther, &uidOther).BindError()
	if err != nil || uidOther == 0 {
		c.Logger().Error("invalid parameter")
		return echo.ErrBadRequest
	}
	// parse blocked
	blocked := false
	err = echo.QueryParamsBinder(c).Bool(keyBlocked, &blocked).BindError()
	if err != nil {
		c.Logger().Error("invalid parameter")
		return echo.ErrBadRequest
	}

	// fetch
	contact, err := dao.FetchContact(uidSelf, uidOther, false)
	if err != nil {
		c.Logger().Error("failed to fetch contact")
		return err
	}
	if contact == nil {
		return echo.ErrNotFound
	}

	err = common.DoInTx(func(tx *sql.Tx) error {
		// delete contact
		err := dao.DeleteContact(tx, contact, blocked)
		if err != nil {
			c.Logger().Error("failed to delete contact")
		}
		return err
	})

	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "")
}

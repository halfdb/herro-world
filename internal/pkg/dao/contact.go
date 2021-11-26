package dao

import (
	"database/sql"
	"errors"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strconv"
	"time"
)

func FetchContact(uidSelf, uidOther int, withDeleted bool) (*models.Contact, error) {
	mods := append(make([]qm.QueryMod, 0),
		models.ContactWhere.UIDSelf.EQ(uidSelf),
		models.ContactWhere.UIDOther.EQ(uidOther),
	)
	if withDeleted {
		mods = append(mods, qm.WithDeleted())
	}
	return models.Contacts(mods...).One(common.GetDB())
}

func ContactExists(uidSelf, uidOther int, withDeleted bool) (bool, error) {
	_, err := FetchContact(uidSelf, uidOther, withDeleted)
	if err == sql.ErrNoRows { // no row
		return false, nil
	} else if err != nil { // error
		return false, err
	} else { // exists
		return true, nil
	}
}

func FetchAllContacts(uid int, withDeleted bool, withBlocked bool) (models.ContactSlice, error) {
	mods := append(make([]qm.QueryMod, 0),
		qm.Select(models.ContactColumns.UIDOther, models.ContactColumns.DisplayName, models.ContactColumns.Cid),
		models.ContactWhere.UIDSelf.EQ(uid),
	)
	if withDeleted {
		mods = append(mods, qm.WithDeleted())
	}
	if !withBlocked {
		mods = append(mods, models.ContactWhere.BlockedAt.IsNull())
	}
	return models.Contacts(mods...).All(common.GetDB())
}

func UpdateContact(executor boil.Executor, uidSelf, uidOther int, updates models.M) error {
	rowsAff, err := models.Contacts(
		models.ContactWhere.UIDSelf.EQ(uidSelf),
		models.ContactWhere.UIDOther.EQ(uidOther),
	).UpdateAll(executor, updates)
	if rowsAff == 0 {
		return sql.ErrNoRows
	} else if err != nil {
		return err
	} else if rowsAff != 1 {
		return errors.New("unexpected: rowsAff = " + strconv.FormatInt(rowsAff, 10))
	}
	return nil
}

func CreateContact(tx *sql.Tx, contact *models.Contact, createChat bool) (*models.Contact, error) {
	if createChat {
		// create chat in advance
		chat, err := CreateChat(tx, nil, true, contact.UIDSelf, contact.UIDOther) // direct chat does not have a name
		if err != nil {
			return nil, err
		}
		contact.Cid = chat.Cid
	}
	// insert the new contact
	if err := contact.Insert(tx, boil.Infer()); err != nil {
		return nil, err
	}
	return contact, nil
}

func DeleteContact(executor boil.Executor, uidSelf, uidOther int, block bool) error {
	// fetch contact
	contact, err := models.FindContact(executor, uidSelf, uidOther)
	if err != nil {
		return err
	}

	// delete user_chat
	err = DeleteUserChat(executor, uidSelf, contact.Cid)
	if err != nil {
		return err
	}

	// block
	if block {
		contact.BlockedAt = null.NewTime(time.Now(), true)
		rowsAff, err := contact.Update(executor, boil.Infer())
		if err != nil {
			return err
		} else if rowsAff != 1 {
			return errors.New("unexpected: rowsAff = " + strconv.FormatInt(rowsAff, 10))
		}
	}
	// delete
	rowsAff, err := contact.Delete(executor, false)
	if rowsAff == 0 {
		return sql.ErrNoRows
	} else if err != nil {
		return echo.ErrInternalServerError
	} else if rowsAff != 1 {
		return errors.New("unexpected: rowsAff = " + strconv.FormatInt(rowsAff, 10))
	}
	return nil
}

func RestoreContact(executor boil.Executor, contact *models.Contact) (*models.Contact, error) {
	contact.DeletedAt.Valid = false
	contact.BlockedAt.Valid = false
	rowsAff, err := contact.Update(executor, boil.Infer())
	if rowsAff == 0 {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, echo.ErrInternalServerError
	} else if rowsAff != 1 {
		return nil, errors.New("unexpected: rowsAff = " + strconv.FormatInt(rowsAff, 10))
	}
	return contact, nil
}

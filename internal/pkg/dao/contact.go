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

func ExistContact(uidSelf, uidOther int, withDeleted bool) (bool, error) {
	_, err := FetchContact(uidSelf, uidOther, withDeleted)
	if err == sql.ErrNoRows { // no row
		return false, nil
	} else if err != nil { // error
		return false, err
	} else { // exists
		return true, nil
	}
}

func LookupAllContacts(uid int, withDeleted bool, withBlocked bool) (models.ContactSlice, error) {
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

func UpdateContact(tx *sql.Tx, contact *models.Contact) error {
	rowsAff, err := contact.Update(tx, boil.Infer())
	if err != nil {
		return err
	}
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
		chat := &models.Chat{
			Direct: true,
		}
		chat, err := CreateChat(tx, chat, contact.UIDSelf, contact.UIDOther) // direct chat does not have a name
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

func DeleteContact(tx *sql.Tx, contact *models.Contact, block bool) error {
	// delete user_chat
	err := DeleteUserChat(tx, &models.UserChat{
		UID: contact.UIDSelf,
		Cid: contact.Cid,
	})
	if err != nil {
		return err
	}

	// block
	if block {
		contact.BlockedAt = null.NewTime(time.Now(), true)
		rowsAff, err := contact.Update(tx, boil.Infer())
		if err != nil {
			return err
		} else if rowsAff != 1 {
			return errors.New("unexpected: rowsAff = " + strconv.FormatInt(rowsAff, 10))
		}
	}
	// delete
	rowsAff, err := contact.Delete(tx, false)
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

package dao

import (
	"database/sql"
	"errors"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strconv"
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

func FetchAllContacts(uid int) (models.ContactSlice, error) {
	return models.Contacts(
		qm.Select(models.ContactColumns.UIDOther, models.ContactColumns.DisplayName, models.ContactColumns.Cid),
		models.ContactWhere.UIDSelf.EQ(uid),
		models.ContactWhere.BlockedAt.IsNull(),
	).All(common.GetDB())
}

func UpdateContact(uidSelf, uidOther int, updates models.M) error {
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

	rowsAff, err := models.Contacts(
		models.ContactWhere.UIDSelf.EQ(uidSelf),
		models.ContactWhere.UIDOther.EQ(uidOther),
	).UpdateAll(tx, updates)
	if rowsAff == 0 {
		return sql.ErrNoRows
	} else if err != nil {
		return echo.ErrInternalServerError
	} else if rowsAff != 1 {
		return errors.New("unexpected: rowsAff = " + strconv.FormatInt(rowsAff, 10))
	}
	return nil
}

func CreateContact(executor boil.Executor, contact *models.Contact) (*models.Contact, error) {
	// create chat in advance
	chat, err := CreateChat(executor, nil, true, contact.UIDSelf, contact.UIDOther) // direct chat does not have a name
	if err != nil {
		return nil, err
	}
	// insert the new contact
	contact.Cid = chat.Cid
	if err := contact.Insert(executor, boil.Infer()); err != nil {
		return nil, err
	}
	return contact, nil
}

func DeleteContact(executor boil.Executor, uidSelf, uidOther int) error {
	contact, err := models.FindContact(executor, uidSelf, uidOther)
	if err != nil {
		return err
	}
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
	rowsAff, err := contact.Update(executor, boil.Whitelist(
		models.ContactColumns.DeletedAt,
		models.ContactColumns.BlockedAt,
		models.ContactColumns.DisplayName,
	))
	if rowsAff == 0 {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, echo.ErrInternalServerError
	} else if rowsAff != 1 {
		return nil, errors.New("unexpected: rowsAff = " + strconv.FormatInt(rowsAff, 10))
	}
	return contact, nil
}

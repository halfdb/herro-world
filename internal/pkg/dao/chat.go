package dao

import (
	"database/sql"
	"errors"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func CreateChat(tx *sql.Tx, chat *models.Chat, uids ...int) (*models.Chat, error) {
	// sanity check
	if len(uids) < 2 || (chat.Direct && len(uids) != 2) || (chat.Direct && chat.Name.Valid) {
		return nil, errors.New("invalid parameter while creating chat")
	}

	if err := chat.Insert(tx, boil.Greylist(models.ChatColumns.Direct)); err != nil {
		return nil, err
	}
	cid := chat.Cid
	for _, uid := range uids {
		userChat := &models.UserChat{
			UID: uid,
			Cid: cid,
		}
		_, err := CreateUserChat(tx, userChat)
		if err != nil {
			return nil, err
		}
	}

	return chat, nil
}

func LookupDirectChat(uid1, uid2 int, withDeleted bool) (int, error) {
	mods := append(make([]qm.QueryMod, 0),
		models.ContactWhere.UIDSelf.EQ(uid1),
		models.ContactWhere.UIDOther.EQ(uid2),
	)
	if withDeleted {
		mods = append(mods, qm.WithDeleted())
	}
	contact, err := models.Contacts(mods...).One(common.GetDB())
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return contact.Cid, nil
}

func FetchChat(cid int, withDeleted bool) (*models.Chat, error) {
	mods := append(make([]qm.QueryMod, 0),
		models.ChatWhere.Cid.EQ(cid),
	)
	if withDeleted {
		mods = append(mods, qm.WithDeleted())
	}
	return models.Chats(mods...).One(common.GetDB())
}

func LookupAllChats(uid int, withDeletedUserChat bool) (models.ChatSlice, error) {
	mods := append(make([]qm.QueryMod, 0),
		qm.Select(
			models.ChatTableColumns.Cid+" as "+models.ChatColumns.Cid,
			models.ChatTableColumns.Direct+" as "+models.ChatColumns.Direct,
			models.ChatTableColumns.Name+" as "+models.ChatColumns.Name,
		),
		qm.InnerJoin(
			models.TableNames.UserChat+" on "+models.UserChatTableColumns.Cid+" = "+models.ChatTableColumns.Cid,
		),
		models.UserChatWhere.UID.EQ(uid),
	)
	if !withDeletedUserChat {
		mods = append(mods, models.UserChatWhere.DeletedAt.IsNull())
	}
	return models.Chats(mods...).All(common.GetDB())
}

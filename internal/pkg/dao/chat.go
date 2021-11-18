package dao

import (
	"errors"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func CreateChat(executor boil.Executor, name *string, direct bool, uids ...int) (*models.Chat, error) {
	// sanity check
	if len(uids) < 2 || (direct && len(uids) != 2) || (direct && name != nil) {
		return nil, errors.New("invalid parameter while creating chat")
	}

	chat := &models.Chat{
		Direct: direct,
	}
	if name != nil {
		chat.Name = null.StringFrom(*name)
	}

	if err := chat.Insert(executor, boil.Infer()); err != nil {
		return nil, err
	}
	cid := chat.Cid
	for _, uid := range uids {
		_, err := CreateUserChat(executor, uid, cid)
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
	if err != nil {
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

func FetchAllChats(uid int, withDeleted bool) (models.ChatSlice, error) {
	mods := append(make([]qm.QueryMod, 0),
		qm.Select(models.ChatTableColumns.Cid+" as `cid`"),
		qm.InnerJoin(
			models.TableNames.UserChat+" on "+models.UserChatTableColumns.Cid+" = "+models.ChatTableColumns.Cid,
		),
		models.UserChatWhere.UID.EQ(uid),
	)
	if withDeleted {
		mods = append(mods, qm.WithDeleted())
	}
	return models.Chats(mods...).All(common.GetDB())
}

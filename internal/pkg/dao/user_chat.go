package dao

import (
	"database/sql"
	"errors"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strconv"
	"time"
)

func FetchUserChat(uid, cid int, withDeleted bool) (*models.UserChat, error) {
	mods := append(make([]qm.QueryMod, 0),
		models.UserChatWhere.UID.EQ(uid),
		models.UserChatWhere.Cid.EQ(cid),
	)
	if withDeleted {
		mods = append(mods, qm.WithDeleted())
	}
	userChat, err := models.UserChats(mods...).One(common.GetDB())
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return userChat, err
}

func CreateUserChat(tx *sql.Tx, userChat *models.UserChat) (*models.UserChat, error) {
	err := userChat.Insert(tx, boil.Infer())
	if err != nil {
		return nil, err
	}
	return userChat, nil
}

func DeleteUserChat(tx *sql.Tx, userChat *models.UserChat) error {
	userChat.DeletedAt = null.NewTime(time.Now(), true)
	return updateUserChat(tx, userChat)
}

func RestoreUserChat(tx *sql.Tx, userChat *models.UserChat) error {
	userChat.DeletedAt.Valid = false
	return updateUserChat(tx, userChat)
}

func updateUserChat(tx *sql.Tx, userChat *models.UserChat) error {
	rowsAff, err := userChat.Update(tx, boil.Infer())
	if err != nil {
		return err
	}
	if rowsAff != 1 {
		return errors.New("unexpected: rowsAff = " + strconv.FormatInt(rowsAff, 10))
	}
	return nil
}

func LookupMemberUids(withDeleted bool, cids ...int) (map[int][]int, error) {
	mods := append(make([]qm.QueryMod, 0),
		models.UserChatWhere.Cid.IN(cids),
	)
	if withDeleted {
		mods = append(mods, qm.WithDeleted())
	}
	userChats, err := models.UserChats(mods...).All(common.GetDB())
	if err != nil {
		return nil, err
	}
	result := make(map[int][]int)
	for _, userChat := range userChats {
		cid := userChat.Cid
		uid := userChat.UID
		result[cid] = append(result[cid], uid)
	}
	return result, nil
}

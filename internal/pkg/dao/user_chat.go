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
	return models.UserChats(mods...).One(common.GetDB())
}

func CreateUserChat(executor boil.Executor, uid, cid int) (*models.UserChat, error) {
	userChat := &models.UserChat{
		UID: uid,
		Cid: cid,
	}
	err := userChat.Insert(executor, boil.Infer())
	if err != nil {
		return nil, err
	}
	return userChat, nil
}

func DeleteUserChat(executor boil.Executor, uid, cid int) error {
	userChat := &models.UserChat{
		UID:       uid,
		Cid:       cid,
		DeletedAt: null.NewTime(time.Now(), true),
	}
	return updateUserChat(executor, userChat)
}

func RestoreUserChat(executor boil.Executor, uid, cid int) error {
	userChat := &models.UserChat{
		UID:       uid,
		Cid:       cid,
		DeletedAt: null.Time{Valid: false},
	}
	return updateUserChat(executor, userChat)
}

func updateUserChat(executor boil.Executor, userChat *models.UserChat) error {
	rowsAff, err := userChat.Update(executor, boil.Infer())
	if rowsAff == 0 {
		return sql.ErrNoRows
	} else if err != nil {
		return err
	} else if rowsAff != 1 {
		return errors.New("unexpected: rowsAff = " + strconv.FormatInt(rowsAff, 10))
	}
	return nil
}

func GetUids(withDeleted bool, cids ...int) (map[int][]int, error) {
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

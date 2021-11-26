package dao

import (
	"database/sql"
	"github.com/halfdb/herro-world/internal/pkg/common"
	"github.com/halfdb/herro-world/internal/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func FetchAllMessages(cid, limit int) (models.MessageSlice, error) {
	return models.Messages(models.MessageWhere.Cid.EQ(cid), qm.Limit(limit)).All(common.GetDB())
}

func CreateMessage(tx *sql.Tx, message *models.Message) (*models.Message, error) {
	err := message.Insert(tx, boil.Infer())
	if err != nil {
		return nil, err
	}
	return message, nil
}

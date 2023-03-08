package repo

import (
	"errors"
	"learn-im/internal/logic/domain/message/model"
	"learn-im/pkg/db"
	"learn-im/pkg/gerrors"

	"gorm.io/gorm"
)

type MessageRepo struct {
}

func NewMessageRepo() *MessageRepo {
	return &MessageRepo{}
}

func (m *MessageRepo) SaveMsg(msg *model.Message) error {
	return db.DB.Model(&model.Message{}).Create(msg).Error
}

func (m *MessageRepo) GetMsgBySeq(userId, seq int, limit int64) ([]*model.Message, bool, error) {
	db := db.DB.Model(&model.Message{}).Where("user_id = ? and seq > ?", userId, seq)
	var count int64
	db.Count(&count)
	if count == 0 {
		return nil, false, nil
	}
	messages := make([]*model.Message, 0)
	if err := db.Limit(int(limit)).Find(&messages).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, false, gerrors.WarpError(err)
	}
	return messages, count > limit, nil
}

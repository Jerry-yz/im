package repo

import (
	"database/sql"
	"learn-im/pkg/db"
	"learn-im/pkg/gerrors"
)

type SeqRepo struct {
}

func NewSeqRepo() *SeqRepo {
	return &SeqRepo{}
}

func (s *SeqRepo) Incr(objectType int, objectId int64) (int64, error) {
	tx := db.DB.Begin()
	defer tx.Rollback()
	var seq int64
	err := db.DB.Raw("select seq from seq where object_type = ? and object_id = ? for update", objectType, objectId).Row().Scan(&seq)
	if err != nil && err != sql.ErrNoRows {
		return 0, gerrors.WarpError(err)
	}
	if err == sql.ErrNoRows {
		if err := db.DB.Exec("insert into seq (object_type, object_id, seq) values(?,?,?)", objectType, objectId, seq+1).Error; err != nil {
			return 0, gerrors.WarpError(err)
		}
	} else {
		if err := db.DB.Exec("update set seq = seq + 1 where object_type = ? and object_id = ?", objectType, objectId).Error; err != nil {
			return 0, gerrors.WarpError(err)
		}
	}
	tx.Commit()
	return seq + 1, nil
}

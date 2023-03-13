package service

import (
	"context"
	"learn-im/internal/logic/domain/message/repo"
)

type SeqService struct {
	*repo.SeqRepo
}

func NewSeqService() *SeqService {
	return &SeqService{
		repo.NewSeqRepo(),
	}
}

func (s *SeqService) GetNextSeq(ctx context.Context, userId int) (int64, error) {
	return s.SeqRepo.Incr(repo.SeqObjectTypeUser, int64(userId))
}

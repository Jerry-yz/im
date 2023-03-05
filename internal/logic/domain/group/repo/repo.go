package repo

import (
	"learn-im/internal/logic/domain/group/entity"
	"learn-im/pkg/gerrors"
)

type GroupRepo struct {
	*GroupCache
	*GroupDao
	*GroupUserRepo
}

func NewGroupRepo() *GroupRepo {
	return &GroupRepo{
		NewGroupCache(),
		NewGroupDao(),
		NewGroupUserRepo(),
	}
}

func (g *GroupRepo) GetGroup(groupId int) (*entity.Group, error) {
	group, err := g.GroupCache.GetGroupCache(groupId)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	if group != nil {
		return group, nil
	}
	group, err = g.GroupDao.GetGroup(groupId)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	users, err := g.GroupUserRepo.ListGroupUsers(groupId)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	for _, user := range users {
		group.Members = append(group.Members, *user)
	}
	if err := g.GroupCache.SetGroupCache(group); err != nil {
		return nil, gerrors.WarpError(err)
	}
	return group, nil
}

func (g *GroupRepo) SaveGroup(groupReq *entity.Group) error {
	if err := g.GroupDao.Save(groupReq); err != nil {
		return gerrors.WarpError(err)
	}
	for k, v := range groupReq.Members {
		if v.UpdateType == entity.UpdateTypeUpdate {
			if err := g.GroupUserRepo.SaveGroupUser(&groupReq.Members[k]); err != nil {
				return gerrors.WarpError(err)
			}
		}
		if v.UpdateType == entity.UpdateTypeDelete {
			if err := g.GroupUserRepo.DeleteGroupUser(groupReq.Members[k].GroupId, groupReq.Members[k].UserId); err != nil {
				return gerrors.WarpError(err)
			}
		}
	}
	if groupReq.ID != 0 {
		if err := g.GroupCache.DelGroupCache(int(groupReq.ID)); err != nil {
			return gerrors.WarpError(err)
		}
	}
	return nil
}

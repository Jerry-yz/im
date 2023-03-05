package group

import (
	"context"
	"learn-im/internal/logic/domain/group/entity"
	"learn-im/internal/logic/domain/group/repo"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/protocol/pb"
)

type GroupApp struct {
	*repo.GroupRepo
}

func NewGroupApp() *GroupApp {
	return &GroupApp{
		repo.NewGroupRepo(),
	}
}

func (g *GroupApp) CreateGroup(ctx context.Context, userId int, req *pb.CreateGroupReq) (int, error) {
	group := entity.CreateGroup(int64(userId), req)
	if err := g.GroupRepo.SaveGroup(group); err != nil {
		return 0, gerrors.WarpError(err)
	}
	return int(group.ID), nil
}

func (g *GroupApp) GetGroup(ctx context.Context, req *pb.GetGroupReq) (*pb.Group, error) {
	group, err := g.GroupRepo.GetGroup(int(req.GroupId))
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	return group.ToProto(), nil
}

func (g *GroupApp) GetUserGroups(ctx context.Context, userId int) ([]*pb.Group, error) {
	groups, err := g.GroupUserRepo.GetGroupsByUserId(userId)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	res := make([]*pb.Group, 0)
	for _, group := range groups {
		res = append(res, group.ToProto())
	}
	return res, nil
}

func (g *GroupApp) UpdateGroup(ctx context.Context, userId int, req *pb.UpdateGroupReq) error {
	group, err := g.GroupRepo.GetGroup(int(req.GroupId))
	if err != nil {
		return gerrors.WarpError(err)
	}
	if err := group.UpdateGroup(ctx, req); err != nil {
		return gerrors.WarpError(err)
	}
	if err := g.GroupRepo.Save(group); err != nil {
		return gerrors.WarpError(err)
	}
	return group.PushUpdate(ctx, userId)
}

func (g *GroupApp) AddMembers(ctx context.Context, userId, groupId int, userIds []int) ([]int, error) {
	group, err := g.GroupRepo.GetGroup(groupId)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	existIds, addIds, err := group.AddMembers(ctx, userIds)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	if err := g.GroupRepo.Save(group); err != nil {
		return nil, gerrors.WarpError(err)
	}
	if err := group.PushAddMember(ctx, userId, addIds); err != nil {
		return nil, gerrors.WarpError(err)
	}
	return existIds, nil
}

func (g *GroupApp) UpdateGroupMember(ctx context.Context, req *pb.UpdateGroupMemberReq) error {
	group, err := g.GroupRepo.GetGroup(int(req.GroupId))
	if err != nil {
		return gerrors.WarpError(err)
	}
	if err := group.UpdateMember(ctx, req); err != nil {
		return gerrors.WarpError(err)
	}
	if err := g.GroupRepo.Save(group); err != nil {
		return gerrors.WarpError(err)
	}
	return nil
}

func (g *GroupApp) DeleteGroupMember(ctx context.Context, groupId, userId, optId int) error {
	group, err := g.GroupRepo.GetGroup(groupId)
	if err != nil {
		return gerrors.WarpError(err)
	}
	if err := group.RemoveMember(ctx, userId); err != nil {
		return gerrors.WarpError(err)
	}
	if err := g.GroupRepo.DeleteGroupUser(groupId, userId); err != nil {
		return gerrors.WarpError(err)
	}
	return group.PushDeleteMember(ctx, optId, userId)
}

func (g *GroupApp) GetGroupMember(ctx context.Context, req *pb.GetGroupMembersReq) ([]*pb.GroupMember, error) {
	group, err := g.GroupRepo.GetGroup(int(req.GroupId))
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	return group.GetMembers(ctx)
}

func (g *GroupApp) SendMessage(ctx context.Context, fromDeviceId, fromUserId int, req *pb.SendMessageReq) (int64, error) {
	group, err := g.GroupRepo.GetGroup(int(req.ReceiverId))
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	return group.SendMessage(ctx, fromDeviceId, fromUserId, req)
}

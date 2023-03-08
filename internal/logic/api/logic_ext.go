package api

import (
	"context"
	"learn-im/internal/logic/domain/device"
	"learn-im/internal/logic/domain/friend"
	"learn-im/internal/logic/domain/group"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/grpclib"
	"learn-im/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type LogicExtServer struct {
	pb.UnsafeLogicExtServer
}

var deviceApp = device.NewDeviceApp()
var friendApp = friend.NewFriendApp()
var groupApp = group.NewGroupApp()

func (l *LogicExtServer) RegisterDevice(ctx context.Context, req *pb.RegisterDeviceReq) (*pb.RegisterDeviceResp, error) {
	deviceId, err := deviceApp.Register(ctx, req)
	return &pb.RegisterDeviceResp{DeviceId: deviceId}, err
}

func (l *LogicExtServer) SendMessageToFriend(ctx context.Context, req *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	userId, deviceId, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return &pb.SendMessageResp{}, gerrors.WarpError(err)
	}
	seq, err := friendApp.SendToFriend(ctx, int(deviceId), int(userId), req)
	if err != nil {
		return &pb.SendMessageResp{}, gerrors.WarpError(err)
	}
	return &pb.SendMessageResp{Seq: seq}, nil
}

func (l *LogicExtServer) AddFriend(ctx context.Context, req *pb.AddFriendReq) (*emptypb.Empty, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	if err := friendApp.AddFriend(ctx, int(userId), int(req.FriendId), req.Remarks, req.Description); err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &emptypb.Empty{}, nil
}

func (l *LogicExtServer) AgreeAddFriend(ctx context.Context, req *pb.AgreeAddFriendReq) (*emptypb.Empty, error) {
	friendId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	if err := friendApp.AgreeAddFriend(ctx, int(req.UserId), int(friendId)); err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &emptypb.Empty{}, nil
}

func (l *LogicExtServer) SetFriend(ctx context.Context, req *pb.SetFriendReq) (*pb.SetFriendResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	if err := friendApp.SetFriend(ctx, int(userId), req); err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &pb.SetFriendResp{FriendId: userId}, nil
}

func (l *LogicExtServer) GetFriend(ctx context.Context, req *emptypb.Empty) (*pb.GetFriendsResp, error) {
	userId, deviceId, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	friends, err := friendApp.List(ctx, int(userId), int(deviceId))
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	res := new(pb.GetFriendsResp)
	for _, friend := range friends {
		res.Friends = append(res.Friends, &pb.Friend{
			UserId:    friend.UserId,
			Phone:     friend.Phone,
			NickName:  friend.NickName,
			Sex:       friend.Sex,
			AvatorUrl: friend.AvatorUrl,
			UserExtra: friend.UserExtra,
			Remarks:   friend.Remarks,
			Extra:     friend.Extra,
		})
	}
	return res, nil
}

func (l *LogicExtServer) SendMessageToGroup(ctx context.Context, req *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	userId, deviceId, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	seq, err := groupApp.SendMessage(ctx, int(deviceId), int(userId), req)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &pb.SendMessageResp{Seq: seq}, nil
}

func (l *LogicExtServer) CreateGroup(ctx context.Context, req *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	groupId, err := groupApp.CreateGroup(ctx, int(userId), req)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &pb.CreateGroupResp{GorupId: int64(groupId)}, nil
}

func (l *LogicExtServer) UpdateGroup(ctx context.Context, req *pb.UpdateGroupReq) (*emptypb.Empty, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	if err := groupApp.UpdateGroup(ctx, int(userId), req); err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &emptypb.Empty{}, nil
}

func (l *LogicExtServer) GetGroup(ctx context.Context, req *pb.GetGroupReq) (*pb.GetGroupResp, error) {
	group, err := groupApp.GetGroup(ctx, req)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &pb.GetGroupResp{Group: group}, nil
}

func (l *LogicExtServer) GetGroups(ctx context.Context, req *emptypb.Empty) (*pb.GetGroupsResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	groups, err := groupApp.GetUserGroups(ctx, int(userId))
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &pb.GetGroupsResp{Groups: groups}, nil
}

func (l *LogicExtServer) AddGroupMembers(ctx context.Context, req *pb.AddGroupMembersReq) (*pb.AddGroupMembersResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	userIds := make([]int, 0)
	for _, id := range req.UserIds {
		userIds = append(userIds, int(id))
	}
	res, err := groupApp.AddMembers(ctx, int(userId), int(req.GroupId), userIds)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	ids := make([]int64, 0)
	for _, id := range res {
		ids = append(ids, int64(id))
	}
	return &pb.AddGroupMembersResp{UserIds: ids}, nil
}

func (l *LogicExtServer) UpdateGroupMember(ctx context.Context, req *pb.UpdateGroupMemberReq) (*emptypb.Empty, error) {
	if err := groupApp.UpdateGroupMember(ctx, req); err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &emptypb.Empty{}, nil
}

func (l *LogicExtServer) DeleteGroupMember(ctx context.Context, req *pb.DeleteGroupMemberReq) (*emptypb.Empty, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	if err := groupApp.DeleteGroupMember(ctx, int(req.GroupId), int(req.UserId), int(userId)); err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &emptypb.Empty{}, nil
}

func (l *LogicExtServer) GetGroupMembers(ctx context.Context, req *pb.GetGroupMembersReq) (*pb.GetGroupMembersResp, error) {
	res, err := groupApp.GetGroupMember(ctx, req)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &pb.GetGroupMembersResp{Members: res}, nil
}

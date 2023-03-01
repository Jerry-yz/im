package group

import (
	"context"
	"errors"
	"learn-im/internal/logic/proxy"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/grpclib"
	"learn-im/pkg/protocol/pb"
	"learn-im/pkg/rpc"
	"learn-im/pkg/util"
	"time"

	// "github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/proto"
	// "go.starlark.net/lib/proto"
	"gorm.io/gorm"
)

const (
	UpdateTypeUpdate = 1
	UpdateTypeDelete = 2
)

type Group struct {
	// Id           int64
	Name         string
	AvatarUrl    string
	Introduction string
	Extra        string
	UserNum      int32
	Members      []GroupUser `gorm:"-:群成员"`
	gorm.Model
}

type GroupUser struct {
	// gorm.Model
	UserId     int `gorm:"unique_index:user_group_id"`
	GroupId    int `gorm:"unique_index:user_group_id"`
	Status     int
	MemberType int
	Remarks    string
	Extra      string
	UpdateType int
}

func (g *Group) ToProto() *pb.Group {
	if g == nil {
		return nil
	}
	return &pb.Group{
		GroupId:      int64(g.ID),
		Name:         g.Name,
		AvatorUrl:    g.AvatarUrl,
		Introduction: g.Introduction,
		UserNum:      int64(g.UserNum),
		Extra:        g.Extra,
	}
}

func CreateGroup(userId int64, in *pb.CreateGroupReq) *Group {
	group := &Group{
		Name:         in.Name,
		AvatarUrl:    in.AvatorUrl,
		Introduction: in.Introduction,
		Extra:        in.Extra,
	}
	group.Members = append(group.Members,
		GroupUser{
			UserId:     int(userId),
			GroupId:    int(group.ID),
			MemberType: int(pb.MemberType_GMT_ADMIN),
			UpdateType: UpdateTypeUpdate,
		},
	)
	for _, memberId := range in.MemberIds {
		group.Members = append(group.Members, GroupUser{
			UserId:     int(memberId),
			GroupId:    int(group.ID),
			MemberType: int(pb.MemberType_GMT_MEMBER),
			UpdateType: UpdateTypeUpdate,
		})
	}
	return group
}

func (g *Group) UpdateGroup(ctx context.Context, req *pb.UpdateGroupReq) error {
	g.Name = req.Name
	g.AvatarUrl = req.AvatorUrl
	g.Introduction = req.Introduction
	g.Extra = req.Extra
	return nil
}

func (g *Group) PushUpdate(ctx context.Context, userId int) error {
	resp, err := rpc.GetBusinessIntClient().GetUser(ctx, &pb.GetUserReq{UserId: int64(userId)})
	if err != nil {
		return gerrors.WarpError(err)
	}
	return g.PushMessage(ctx, pb.PushCode_PC_GROUP_MESSAGE, pb.UpdateGroupPush{
		OptId:        int64(userId),
		OptName:      resp.User.Nickname,
		Name:         g.Name,
		AvatarUrl:    g.AvatarUrl,
		Introduction: g.Introduction,
		Extra:        g.Extra,
	}, true)
}

func (g *Group) SendMessage(ctx context.Context, fromDeviceId, fromUserId int, req *pb.SendMessageReq) (int64, error) {
	if !g.IsMember(fromUserId) {
		return 0, gerrors.WarpError(errors.New("不是群組用戶"))
	}
	sender, err := rpc.GetSender(fromDeviceId, fromUserId)
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	push := &pb.UserMessagePush{
		Sender:     sender,
		ReceiverId: req.ReceiverId,
		Content:    req.Content,
	}
	byt, err := proto.Marshal(push)
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	msg := &pb.Message{
		Code:     int32(pb.PushCode_PC_GROUP_MESSAGE),
		Content:  byt,
		SendTime: time.Now().Unix(),
	}
	seq, err := proxy.MessageProxy.SendToUser(ctx, fromDeviceId, int64(fromUserId), msg, true)
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	go func() {
		defer util.Recover()
		for _, member := range g.Members {
			_, err := proxy.MessageProxy.SendToUser(ctx, fromDeviceId, int64(member.ID), msg, true)
			if err != nil {
				return
			}
		}
	}()
	return seq, nil
}

func (g *Group) IsMember(userId int) bool {
	for _, member := range g.Members {
		if uint(member.UserId) == uint(userId) {
			return true
		}
	}
	return false
}

func (g *Group) GetMembers(ctx context.Context) ([]*pb.GroupMember, error) {
	userMap := make(map[int64]int32, len(g.Members))
	for i := range g.Members {
		userMap[int64(g.Members[i].UserId)] = 0
	}
	resp, err := rpc.GetBusinessIntClient().GetUsers(ctx, &pb.GetUsersReq{UserIds: userMap})
	if err != nil {
		return []*pb.GroupMember{}, gerrors.WarpError(err)
	}
	groupMembers := make([]*pb.GroupMember, len(resp.Users))
	for i := range g.Members {
		member := &pb.GroupMember{
			UserId:     int64(g.Members[i].UserId),
			MemberType: pb.MemberType(g.Members[i].MemberType),
			Remarks:    g.Members[i].Remarks,
			Extra:      g.Members[i].Extra,
		}
		if user, ok := resp.Users[member.UserId]; ok {
			member.NickName = user.Nickname
			member.AvatorUrl = user.AvatarUrl
			member.Sex = int64(user.Sex)
			member.Extra = user.Extra
		}
	}
	return groupMembers, nil
}

func (g *Group) AddMembers(ctx context.Context, userIds []int) ([]int, []int, error) {
	existIds := make([]int, 0)
	addIds := make([]int, 0)
	for i := range userIds {
		if g.IsMember(userIds[i]) {
			existIds = append(existIds, userIds[i])
			continue
		}
		g.Members = append(g.Members, GroupUser{
			UserId:     userIds[i],
			GroupId:    int(g.ID),
			MemberType: int(pb.MemberType_GMT_MEMBER),
			Extra:      g.Extra,
			UpdateType: UpdateTypeUpdate,
		})
		addIds = append(addIds, userIds[i])
	}
	return existIds, addIds, nil
}

func (g *Group) PushAddMember(ctx context.Context, optId int, addUserIds []int) error {
	addIdMap := make(map[int64]int32, len(addUserIds))
	for i := range addUserIds {
		addIdMap[int64(addUserIds[i])] = 0
	}
	resp, err := rpc.GetBusinessIntClient().GetUsers(ctx, &pb.GetUsersReq{UserIds: addIdMap})
	if err != nil {
		return gerrors.WarpError(err)
	}
	members := make([]*pb.GroupMember, len(addIdMap))
	for i := range addIdMap {
		user, ok := resp.Users[i]
		if !ok {
			continue
		}
		members = append(members, &pb.GroupMember{
			UserId:    user.UserId,
			NickName:  user.Nickname,
			AvatorUrl: user.AvatarUrl,
			Sex:       int64(user.Sex),
			Remarks:   "",
			Extra:     "",
		})
	}
	optUser := resp.Users[int64(optId)]
	return g.PushMessage(ctx, pb.PushCode_PC_ADD_GROUP_MEMBERS, &pb.AddGroupMembersPush{
		OptId:   optUser.UserId,
		OptName: optUser.Nickname,
		Members: members,
	}, true)
}

func (g *Group) PushMessage(ctx context.Context, code pb.PushCode, msg proto.Message, isPersist bool) error {
	go func() {
		defer util.Recover()
		for _, member := range g.Members {
			_, err := proxy.PushToUser(grpclib.NewAndCopyRequestId(ctx), int64(member.UserId), code, msg, isPersist)
			if err != nil {
				return
			}
		}
	}()
	return nil
}

func (g *Group) GetMember(ctx context.Context, userId int) *GroupUser {
	for i := range g.Members {
		if g.Members[i].UserId == userId {
			return &g.Members[i]
		}
	}
	return nil
}

func (g *Group) UpdateMember(ctx context.Context, req *pb.UpdateGroupMemberReq) error {
	member := g.GetMember(ctx, int(req.UserId))
	if member == nil {
		return nil
	}
	member.MemberType = int(req.MemberType)
	member.Extra = req.Extra
	member.GroupId = int(req.GroupId)
	member.Remarks = req.Remarks
	member.UpdateType = UpdateTypeUpdate
	return nil
}

func (g *Group) RemoveMember(ctx context.Context, userId int) error {
	member := g.GetMember(ctx, userId)
	if member == nil {
		return nil
	}
	member.MemberType = UpdateTypeDelete
	return nil
}

func (g *Group) PushDeleteMember(ctx context.Context, optId, userId int) error {
	resp, err := rpc.GetBusinessIntClient().GetUser(ctx, &pb.GetUserReq{UserId: int64(optId)})
	if err != nil {
		return gerrors.WarpError(err)
	}
	return g.PushMessage(ctx, pb.PushCode_PC_REMOVE_GROUP_MEMBER, &pb.RemoveGroupMemberPush{
		OptId:         int64(optId),
		OptName:       resp.User.Nickname,
		DeletedUserId: int64(userId),
	}, true)
}

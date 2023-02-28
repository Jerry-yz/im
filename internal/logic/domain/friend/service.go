package friend

import (
	"context"
	"learn-im/internal/logic/proxy"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/protocol/pb"
	"learn-im/pkg/rpc"

	"github.com/golang/protobuf/proto"
)

type FriendService struct {
	*FriendRepo
}

func NewFriendService() *FriendService {
	return &FriendService{
		NewFriendRepo(),
	}
}

func (f *FriendService) ListFriend(ctx context.Context, userId, status int) ([]*pb.Friend, error) {
	friends := make([]*Friend, 0)
	friends, err := f.FriendRepo.ListFriend(userId, status)
	if err != nil {
		return []*pb.Friend{}, gerrors.WarpError(err)
	}
	userIdsMap := make(map[int64]int32, 0)
	for _, friend := range friends {
		userIdsMap[friend.FriendId] = 0
	}
	users, err := rpc.GetBusinessIntClient().GetUsers(ctx, &pb.GetUsersReq{
		UserIds: userIdsMap,
	})
	if err != nil {
		return []*pb.Friend{}, gerrors.WarpError(err)
	}
	infos := make([]*pb.Friend, 0)
	for k, friend := range friends {
		user, ok := users.Users[friend.FriendId]
		if !ok {
			continue
		}
		infos = append(infos, &pb.Friend{
			UserId:    friends[k].FriendId,
			Phone:     friends[k].Extra,
			NickName:  user.Nickname,
			AvatorUrl: user.AvatarUrl,
			UserExtra: user.Extra,
			Remarks:   friends[k].Remarks,
			Sex:       int64(user.Sex),
		})
	}
	return infos, nil
}

func (f *FriendService) AddFriend(ctx context.Context, userId, friendId int, remarks, description string) error {
	friend, err := f.FriendRepo.GetFriend(userId, friendId)
	if err != nil {
		return gerrors.WarpError(err)
	}
	if friend != nil {
		if friend.Status == FriendStatusAgree {
			return gerrors.ErrAreadyIsFriend
		}
		if friend.Status == FriendStatusApply {
			return nil
		}
	}
	if err := f.FriendRepo.Save(&Friend{
		UserId:   int64(userId),
		FriendId: int64(friendId),
		Remarks:  friend.Remarks,
		Extra:    friend.Extra,
		Status:   FriendStatusApply,
	}); err != nil {
		return gerrors.WarpError(err)
	}
	_, err = rpc.GetBusinessIntClient().GetUser(ctx, &pb.GetUserReq{
		UserId: friend.FriendId,
	})
	if err != nil {
		return gerrors.WarpError(err)
	}
	// TODO pushToUser
	return nil
}

func (f *FriendService) AgreeAddFriend(ctx context.Context, userId, friendId int) error {
	friend, err := f.FriendRepo.GetFriend(userId, friendId)
	if err != nil {
		return gerrors.WarpError(err)
	}
	if friend == nil {
		return gerrors.ErrNotFound
	}
	if friend.Status == FriendStatusAgree {
		return nil
	}
	if err := f.FriendRepo.Save(friend); err != nil {
		return gerrors.WarpError(err)
	}
	if err := f.FriendRepo.Save(&Friend{
		UserId:   int64(userId),
		FriendId: int64(friendId),
		Remarks:  friend.Remarks,
		Extra:    friend.Extra,
		Status:   FriendStatusAgree,
	}); err != nil {
		return gerrors.WarpError(err)
	}
	resp, err := rpc.GetBusinessIntClient().GetUser(ctx, &pb.GetUserReq{UserId: int64(userId)})
	if err != nil {
		return gerrors.WarpError(err)
	}
	_, err = proxy.PushToUser(ctx, friend.Id, pb.PushCode_PC_ADD_FRIEND, &pb.AddFriendPush{
		FriendId:  int64(userId),
		Nickname:  resp.User.Nickname,
		AvatarUrl: resp.User.AvatarUrl,
	}, true)
	if err != nil {
		return gerrors.WarpError(err)
	}
	return nil
}

func (f *FriendService) SendToFriend(ctx context.Context, deviceId, fromId int64, req *pb.SendMessageReq) (int64, error) {
	sender, err := rpc.GetSender(int(deviceId), int(fromId))
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	pushMsg := &pb.UserMessagePush{
		Sender:     sender,
		ReceiverId: req.ReceiverId,
		Content:    req.Content,
	}
	byt, err := proto.Marshal(pushMsg)
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	msg := &pb.Message{
		Code:     int32(pb.PushCode_PC_USER_MESSAGE),
		Content:  byt,
		SendTime: req.SendTime,
	}
	seq, err := proxy.MessageProxy.SendToUser(ctx, deviceId, fromId, msg, true)
	if err != nil {
		return 0, err
	}
	_, err = proxy.MessageProxy.SendToUser(ctx, deviceId, req.ReceiverId, msg, true)
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	return seq, nil
}

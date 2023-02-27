package friend

import (
	"context"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/protocol/pb"
	"learn-im/pkg/rpc"
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
	_, err = rpc.GetBusinessIntClient().GetUser(ctx, &pb.GetUserReq{UserId: int64(userId)})
	if err != nil {
		return gerrors.WarpError(err)
	}
	// TODO pushToUser
	return nil
}

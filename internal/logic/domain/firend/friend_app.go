package friend

import (
	"context"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/protocol/pb"
)

type FriendApp struct {
	*FriendRepo
	*FriendService
}

func NewFriendApp() *FriendApp {
	return &FriendApp{
		NewFriendRepo(),
		NewFriendService(),
	}
}

func (f *FriendApp) GetFriend(ctx context.Context, userId, friendId int) (*Friend, error) {
	friend, err := f.FriendRepo.GetFriend(userId, friendId)
	if err != nil {
		return &Friend{}, gerrors.WarpError(err)
	}
	return friend, nil
}

func (f *FriendApp) List(ctx context.Context, userId, status int) ([]*pb.Friend, error) {
	return f.FriendService.ListFriend(ctx, userId, status)
}

func (f *FriendApp) AddFriend(ctx context.Context, userId, friendId int, remarks, description string) error {
	return f.FriendService.AddFriend(ctx, userId, friendId, remarks, description)
}

func (f *FriendApp) AgreeAddFriend(ctx context.Context, userId, friendId int) error {
	return f.FriendService.AgreeAddFriend(ctx, userId, friendId)
}

func (f *FriendApp) SetFriend(ctx context.Context, userId int, req *pb.SetFriendReq) error {
	friend, err := f.FriendRepo.GetFriend(userId, int(req.FriendId))
	if err != nil {
		return gerrors.WarpError(err)
	}
	return f.FriendRepo.Save(&Friend{
		UserId:   int64(userId),
		FriendId: friend.Id,
		Remarks:  friend.Remarks,
		Extra:    friend.Extra,
		Status:   friend.Status,
	})
}

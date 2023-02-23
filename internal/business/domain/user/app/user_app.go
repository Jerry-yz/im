package app

import (
	"context"
	"learn-im/internal/business/domain/user/repo"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/protocol/pb"
)

type UserApp struct {
	*repo.UserRepo
}

func NewUserApp() *UserApp {
	return &UserApp{
		repo.NewUserRepo(),
	}
}

func (u *UserApp) GetUserApp(ctx context.Context, userId int) (*pb.User, error) {
	user, err := u.UserRepo.Get(userId)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	return user.ToProto(), nil
}

func (u *UserApp) Update(ctx context.Context, userId int, req *pb.UpdateUserReq) error {
	user, err := u.UserRepo.Get(userId)
	if err != nil {
		return gerrors.WarpError(err)
	}
	user.Avator = req.AvatarUrl
	user.Extra = req.Extra
	user.NickName = req.Nickname
	user.Sex = int(req.Sex)
	return u.UserRepo.Save(user)
}

func (u *UserApp) GetByIds(ctx context.Context, userIds []int) (map[int64]*pb.User, error) {
	userList, err := u.UserRepo.GetByIds(userIds)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	res := make(map[int64]*pb.User, len(userList))
	for userId := range userList {
		res[int64(userId)] = userList[userId].ToProto()
	}
	return res, nil
}

func (u *UserApp) Search(ctx context.Context, keyword string) ([]*pb.User, error) {
	userList, err := u.UserRepo.Search(keyword)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	res := make([]*pb.User, len(userList))
	for userId := range userList {
		res = append(res, userList[userId].ToProto())
	}
	return res, nil
}

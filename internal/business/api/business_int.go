package business

import (
	"context"
	"learn-im/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServerInt struct {
	pb.UnsafeBusinessIntServer
}

func NewGrpcServerInt() *GrpcServerInt {
	return &GrpcServerInt{}
}

// var authApp = app.NewAuthApp()
// var authApp = app.NewUserApp()

func (g *GrpcServerInt) Auth(ctx context.Context, req *pb.AuthReq) (*emptypb.Empty, error) {
	err := authApp.Auth(ctx, int(req.UserId), int(req.DeviceId), req.Token)
	return &emptypb.Empty{}, err
}

func (g *GrpcServerInt) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error) {
	user, err := userApp.GetUserApp(ctx, int(req.UserId))
	return &pb.GetUserResp{User: user}, err
}

func (g *GrpcServerInt) GetUsers(ctx context.Context, req *pb.GetUsersReq) (*pb.GetUsersResp, error) {
	// slices.TypeConversion(req.UserIds)
	userIds := make([]int, len(req.UserIds))
	for userId := range req.UserIds {
		userIds = append(userIds, int(userId))
	}
	users, err := userApp.GetByIds(ctx, userIds)
	return &pb.GetUsersResp{Users: users}, err
}

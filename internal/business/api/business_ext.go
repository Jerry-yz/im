package business

import (
	"context"
	"learn-im/internal/business/domain/user/app"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/grpclib"
	"learn-im/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct {
	pb.UnsafeBusinessExtServer
}

func NewServer() *GrpcServer {
	return &GrpcServer{}
}

var authApp = app.NewAuthApp()
var userApp = app.NewUserApp()

func (g *GrpcServer) SignIn(ctx context.Context, in *pb.SignInReq) (*pb.SignInResp, error) {
	if err := authApp.SignIn(ctx, in.PhoneNumber, int(in.DeviceId)); err != nil {
		return &pb.SignInResp{}, gerrors.WarpError(err)
	}
	return &pb.SignInResp{}, nil
}

func (g *GrpcServer) GetUser(ctx context.Context, in *pb.GetUserReq) (*pb.GetUserResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}

	user, err := userApp.GetUserApp(ctx, int(userId))
	return &pb.GetUserResp{User: user}, err
}

func (g *GrpcServer) UpdateUser(ctx context.Context, in *pb.UpdateUserReq) (*emptypb.Empty, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	err = userApp.Update(ctx, int(userId), in)
	return new(emptypb.Empty), err
}

func (g *GrpcServer) SearchUser(ctx context.Context, in *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	users, err := userApp.Search(ctx, in.Key)
	return &pb.SearchUserResp{Users: users}, err
}

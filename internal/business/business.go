package business

import (
	"context"
	"learn-im/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct {
	pb.UnsafeBusinessExtServer
}

func NewServer() *GrpcServer {
	return &GrpcServer{}
}

func (g *GrpcServer) SignIn(ctx context.Context, in *pb.SignInReq) (*pb.SignInResp, error) {
	return nil, nil
}

func (g *GrpcServer) GetUser(ctx context.Context, in *pb.GetUserReq) (*pb.GetUserResp, error) {
	return nil, nil
}

func (g *GrpcServer) UpdateUser(ctx context.Context, in *pb.UpdateUserReq) (*emptypb.Empty, error) {
	return nil, nil
}

func (g *GrpcServer) SearchUser(ctx context.Context, in *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	return nil, nil
}

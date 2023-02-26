package config

import (
	"context"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/protocol/pb"
	"os"

	"google.golang.org/grpc"
)

var builder = map[string]Builder{
	// "default": defaultBuilder,
	// "k8s":     k8sBuilder,
}

var Conf = &Configuration{}

type Builder interface {
	Build() *Configuration
}

type Configuration struct {
	Mysql         string
	RedisHost     string
	RedisPassword string

	PushRoomSubscribeNum int
	PushAllSubscribeNum  int

	ConnectLocalAddr     string
	ConnectRPCListenAddr string
	ConnectTCPListenAddr string
	ConnectWSListenAddr  string

	LogicRPCListenAddr    string
	BusinessRPCListenAddr string
	FileHTTPListenAddr    string

	ConnectIntClientBuilder  func() pb.ConnectIntClient
	LogicIntClientBuilder    func() pb.LogicIntClient
	BusinessIntClientBuilder func() pb.BusinessIntClient
}

func init() {
	env := os.Getenv("im_env")
	build, ok := builder[env]
	if !ok {
		builder[env] = &defaultBuilder{}
	}
	Conf = build.Build()
}

func interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	return gerrors.WarpError(err)
}

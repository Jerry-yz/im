package rpc

import (
	"context"
	"learn-im/config"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/protocol/pb"
)

var (
	connectIntClient  pb.ConnectIntClient
	logicIntClient    pb.LogicIntClient
	businessIntClient pb.BusinessIntClient
)

func GetLogicIntClient() pb.LogicIntClient {
	if logicIntClient == nil {
		logicIntClient = config.Conf.LogicIntClientBuilder()
	}
	return logicIntClient
}

func GetBusinessIntClient() pb.BusinessIntClient {
	if businessIntClient == nil {
		businessIntClient = config.Conf.BusinessIntClientBuilder()
	}
	return businessIntClient
}

func GetConnectIntClient() pb.ConnectIntClient {
	if connectIntClient == nil {
		connectIntClient = config.Conf.ConnectIntClientBuilder()
	}
	return connectIntClient
}

func GetSender(deviceId, userId int) (*pb.Sender, error) {
	resp, err := GetBusinessIntClient().GetUser(context.Background(), &pb.GetUserReq{UserId: int64(userId)})
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	return &pb.Sender{
		UserId:    int64(userId),
		DeviceId:  int64(deviceId),
		AvatarUrl: resp.User.AvatarUrl,
		Nickname:  resp.User.Nickname,
		Extra:     resp.User.Extra,
	}, nil
}

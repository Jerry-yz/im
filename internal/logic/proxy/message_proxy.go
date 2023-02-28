package proxy

import (
	"context"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/protocol/pb"
	"time"

	"google.golang.org/protobuf/proto"
)

// func NewMessageProxy() MessageProxy {
// 	return MessageProxy
// }

var messageProxy MessageProxy

type MessageProxy interface {
	SendToUser(ctx context.Context, fromDeviceId, toUserId int64, message *pb.Message, isPersist bool) (int64, error)
}

func PushToUserBytes(ctx context.Context, toUserId int64, code int32, bytes []byte, isPersist bool) (int64, error) {
	message := &pb.Message{
		Code:     code,
		Content:  bytes,
		SendTime: time.Now().Unix(),
	}
	seq, err := messageProxy.SendToUser(ctx, 0, toUserId, message, isPersist)
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	return seq, nil
}

func PushToUser(ctx context.Context, toUserId int64, code pb.PushCode, msg proto.Message, isPersist bool) (int64, error) {
	byt, err := proto.Marshal(msg)
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	return PushToUserBytes(ctx, toUserId, int32(code), byt, isPersist)
}

package service

import (
	"context"
	"learn-im/internal/logic/domain/message/model"
	"learn-im/internal/logic/domain/message/repo"
	"learn-im/internal/logic/proxy"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/grpclib"
	"learn-im/pkg/grpclib/picker"
	"learn-im/pkg/protocol/pb"
	"learn-im/pkg/rpc"
	"time"

	"github.com/golang/protobuf/proto"
)

const MaxMsgBufLen = 65546
const MsgLimit = 50

type MessageService struct {
	*DeviceServiceACK
	*repo.MessageRepo
	// proxy.DeviceProxy
}

func NewMessageService() *MessageService {
	return &MessageService{
		NewDeviceACK(),
		repo.NewMessageRepo(),
	}
}

func (m *MessageService) Sync(ctx context.Context, userId, seq int) (*pb.SyncResp, error) {
	messages, err := m.ListByUserIdAndSeq(ctx, userId, seq)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	msg := model.MessagesToPB(messages)
	resp := &pb.SyncResp{Messages: msg, HasMore: false}
	msgLength := len(msg)
	byt, err := proto.Marshal(resp)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	for len(byt) > MaxMsgBufLen {
		msgLength = msgLength * 2 / 3
		resp = &pb.SyncResp{Messages: msg[:msgLength], HasMore: true}
		byt, err = proto.Marshal(resp)
		if err != nil {
			return nil, gerrors.WarpError(err)
		}
	}
	return resp, nil
}

// ListByUserIdAndSeq 查询消息
func (m *MessageService) ListByUserIdAndSeq(ctx context.Context, userId, seq int) ([]*model.Message, error) {
	var err error
	if seq == 0 {
		seq, err = m.DeviceServiceACK.GetMaxByUserId(ctx, userId)
		if err != nil {
			return nil, gerrors.WarpError(err)
		}
	}
	messages, _, err := m.MessageRepo.GetMsgBySeq(userId, seq, MsgLimit)
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	return messages, nil
}

func (m *MessageService) SendUser(ctx context.Context, fromDeviceId, toUserId int, msg *pb.Message, isPersist bool) (int64, error) {
	//消息是否需要持久化
	if isPersist {
		message := model.Message{
			UserId:    toUserId,
			RequestId: grpclib.GetCtxRequestId(ctx),
			Code:      int(msg.Code),
			Content:   msg.Content,
			Seq:       int(msg.Seq),
			SendTime:  time.Unix(msg.SendTime, 0),
			Status:    int(msg.Status),
		}
		if err := m.MessageRepo.SaveMsg(&message); err != nil {
			return 0, gerrors.WarpError(err)
		}
	}
	//查询所有在线设备
	devices, err := proxy.DevProxy.ListOnlineByUserId(ctx, int64(toUserId))
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	//循环发送消息给每个设备
	for _, device := range devices {
		if device.DeviceId == int64(fromDeviceId) {
			continue
		}
		if err := m.SendToDevice(ctx, device, msg); err != nil {
			return 0, gerrors.WarpError(err)
		}
	}
	return msg.Seq, nil
}

func (m *MessageService) SendToDevice(ctx context.Context, device *pb.Device, msg *pb.Message) error {
	_, err := rpc.GetConnectIntClient().DeliverMessage(picker.ContextWithAddr(ctx, device.ConnAddr), &pb.DeliverMessageReq{
		DeviceId: device.DeviceId,
		Message:  msg,
	})
	if err != nil {
		return gerrors.WarpError(err)
	}
	return nil
}

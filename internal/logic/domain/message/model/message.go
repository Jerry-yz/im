package model

import (
	"learn-im/pkg/protocol/pb"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	// Id        int        `json:"id"`
	UserId    int        `json:"user_id"`
	RequestId int        `json:"request_id"`
	Code      int        `json:"code"`
	Content   []byte     `json:"content"`
	Seq       int        `json:"seq"`
	SendTime  *time.Time `json:"send_time"`
	Status    int        `json:"status"`
}

func (m *Message) MessageToPB() *pb.Message {
	return &pb.Message{
		Code:     int32(m.Code),
		Content:  m.Content,
		Seq:      int64(m.Seq),
		SendTime: m.SendTime.Unix(),
		Status:   pb.MessageStatus(m.Status),
	}
}

func MessagesToPB(msg []*Message) []*pb.Message {
	res := make([]*pb.Message, 0)
	for _, message := range msg {
		res = append(res, &pb.Message{
			Code:     int32(message.Code),
			Content:  message.Content,
			Seq:      int64(message.Seq),
			SendTime: message.SendTime.Unix(),
			Status:   pb.MessageStatus(message.Status),
		})
	}
	return res
}

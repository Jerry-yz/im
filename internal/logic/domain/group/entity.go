package group

import (
	"learn-im/pkg/protocol/pb"

	"gorm.io/gorm"
)

const (
	UpdateTypeUpdate = 1
	UpdateTypeDelete = 2
)

type Group struct {
	// Id           int64
	Name         string
	AvatarUrl    string
	Introduction string
	Extra        string
	UserNum      int32
	Members      []GroupUser `gorm:"-:群成员"`
	gorm.Model
}

type GroupUser struct {
	gorm.Model
	UserId     int `gorm:"unique_index:user_group_id"`
	GroupId    int `gorm:"unique_index:user_group_id"`
	Status     int
	MemberType int
	Remarks    string
	Extra      string
	UpdateType int
}

func (g *Group) ToProto() *pb.Group {
	if g == nil {
		return nil
	}
	return &pb.Group{
		GroupId:      int64(g.ID),
		Name:         g.Name,
		AvatorUrl:    g.AvatarUrl,
		Introduction: g.Introduction,
		UserNum:      int64(g.UserNum),
		Extra:        g.Extra,
	}
}

func CreateGroup(userId int64, in *pb.CreateGroupReq) *Group {

}

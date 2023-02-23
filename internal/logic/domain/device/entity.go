package device

import "learn-im/pkg/protocol/pb"

const (
	DeviceOnLine = 1
	DeViceOffLie = 0
)

type Device struct {
	// gorm.Model
	Id            int
	UserId        int
	Type          int
	Brand         string
	Model         string
	SystemVersion string
	SdkVersion    string
	Status        int
	ConnAddr      string
	ClientAddr    string
}

func (d *Device) ToProto() *pb.Device {
	return &pb.Device{
		DeviceId:      int64(d.Id),
		UserId:        int64(d.UserId),
		Type:          int32(d.Type),
		Brand:         d.Brand,
		Model:         d.Model,
		SystemVersion: d.SystemVersion,
		SdkVersion:    d.SdkVersion,
		Status:        int32(d.Status),
		ConnAddr:      d.ConnAddr,
		ClientAddr:    d.ClientAddr,
	}
}

func (d *Device) IsLegal() bool {
	return d.Type == 0 || d.Brand == "" || d.Model == "" ||
		d.SystemVersion == "" || d.SdkVersion == ""
}

func (d *Device) OnLien(userId int, connAddr, clientAddr string) {
	d.UserId = userId
	d.ConnAddr = connAddr
	d.ClientAddr = clientAddr
	d.Status = DeviceOnLine
}

func (d *Device) OffLine(userId int, connAddr, clientAddr string) {
	d.UserId = userId
	d.ConnAddr = connAddr
	d.ClientAddr = clientAddr
	d.Status = DeViceOffLie
}

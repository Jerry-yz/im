package device

type DeviceRepo struct {
	*Dao
}

func NewDeviceRepo() *DeviceRepo {
	return &DeviceRepo{
		NewDao(),
	}
}

func (d *DeviceRepo) GetDevice(deviceId int) (*Device, error) {
	return d.Dao.GetDevice(deviceId)
}

func (d *DeviceRepo) Save(device *Device) error {
	return d.Dao.Save(device)
}

func (d *DeviceRepo) ListOnlineByUserId(userId int) ([]*Device, error) {
	return d.Dao.ListOnlineByUserId(userId)
}

func (d *DeviceRepo) ListOnlineByConnAddr(connAddr string) ([]*Device, error) {
	return d.Dao.ListOnlineByConnAddr(connAddr)
}

func (d *DeviceRepo) UpdateDeviceStatus(deviceId, status int, connAddr string) error {
	return d.Dao.UpdateStatus(deviceId, status, connAddr)
}

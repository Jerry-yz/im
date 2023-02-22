package model

type Device struct {
	Type   string `json:"type"` //设备类型
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

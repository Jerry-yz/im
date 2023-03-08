package model

type Sender struct {
	UserId    int    `json:"user_id"`
	DeviceId  int    `json:"device_id"`
	NickName  string `json:"nick_name"`
	AvatarUrl string `json:"avatar_url"`
	Extra     string `json:"extra"`
}

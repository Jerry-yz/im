package config

import "os"

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

	// pb.ConnectIntClient/
}

func init() {
	env := os.Getenv("im_env")
	build, ok := builder[env]
	if !ok {
		// builder[env] = defaultBuilder
	}
	Conf = build.Build()
}

package connect

import (
	"learn-im/logger"
	"time"

	"github.com/alberliu/gn"
	"github.com/alberliu/gn/codec"
	"go.uber.org/zap"
)

var server *gn.Server

type handler struct {
}

func StartTCPServer(addr string) {
	gn.SetLogger(logger.Sugar)
	var err error
	srv, err := gn.NewServer(addr, &handler{}, gn.WithDecoder(codec.NewUvarintDecoder()),
		gn.WithEncoder(codec.NewUvarintEncoder(1024)),
		gn.WithReadBufferLen(256),
		gn.WithTimeout(11*time.Minute),
		gn.WithAcceptGNum(10),
		gn.WithIOGNum(100),
	)
	if err != nil {
		return
	}
	srv.Run()
}

func (h *handler) OnConnect(c *gn.Conn) {
	conn := &Conn{
		CoonType: ConnTypeTCP,
		TCP:      c,
	}
	c.SetData(conn)
	logger.Logger.Debug("connect:", zap.Int32("fd", c.GetFd()), zap.String("addr", c.GetAddr()))
}

func (h *handler) OnMessage(c *gn.Conn, bytes []byte) {
	// conn := c.GetData().(*Conn) TODO
}

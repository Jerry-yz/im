package connect

import (
	"container/list"
	"sync"

	"github.com/alberliu/gn"
	"github.com/gorilla/websocket"
)

const (
	ConnTypeTCP int8 = 1
	ConnTypeWS  int8 = 2
)

type Conn struct {
	CoonType int8            // 连接类型
	TCP      *gn.Conn        // tcp连接
	WSMutex  sync.Mutex      // WS写锁
	WS       *websocket.Conn // websocket连接
	UserId   int64           // 用户ID
	DeviceId int64           // 设备ID
	RoomId   int64           // 订阅的房间ID
	Element  *list.Element   // 链表节点
}

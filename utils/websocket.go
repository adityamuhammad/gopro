package utils

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ConnectionRegistry struct {
	connections map[uint]*websocket.Conn
	lock        sync.RWMutex
}

func (cr *ConnectionRegistry) Add(userID uint, conn *websocket.Conn) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	cr.connections[userID] = conn
}

func (cr *ConnectionRegistry) Remove(userID uint) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	delete(cr.connections, userID)
}

func (cr *ConnectionRegistry) Get(userID uint) (*websocket.Conn, bool) {
	cr.lock.RLock()
	defer cr.lock.RUnlock()
	conn, exists := cr.connections[userID]
	return conn, exists
}

var Registry = ConnectionRegistry{
	connections: make(map[uint]*websocket.Conn),
}

type Message struct {
	ToUserID uint   `json:"to_user_id"`
	Content  string `json:"content"`
}

package ws

import (
	"golang.org/x/net/websocket"
	"sync"
)

type Manager struct {
	mux sync.RWMutex

	conns map[string]*websocket.Conn
}

func NewManager() *Manager {
	return &Manager{
		conns: make(map[string]*websocket.Conn),
	}
}

func (s *Manager) AcceptConn(conn *websocket.Conn) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.conns[conn.RemoteAddr().String()] = conn
}

package ws

import (
	"sync"
)

type Manager struct {
	mux sync.RWMutex

	conns map[*Client]bool
}

func NewManager() *Manager {
	return &Manager{
		conns: make(map[*Client]bool),
	}
}

func (s *Manager) AcceptConn(client *Client) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.conns[client] = true
}

func (s *Manager) RemoveConn(client *Client) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.conns, client)
}

func (s *Manager) GetClients() map[*Client]bool {
	clientCopy := make(map[*Client]bool)
	s.mux.RLock()
	for c := range s.conns {
		clientCopy[c] = true
	}
	s.mux.RUnlock()
	return clientCopy
}

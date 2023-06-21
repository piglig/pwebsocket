package ws

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"sync"
)

type Manager struct {
	clientMux sync.RWMutex

	conns map[*Client]bool
	Event chan Event

	handlerMux sync.Mutex
	handlers   map[EventType]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		conns:    make(map[*Client]bool),
		Event:    make(chan Event),
		handlers: make(map[EventType]EventHandler),
	}
	m.initEvents()
	return m
}

func (s *Manager) initEvents() {
	s.RegisterEvent(SingleChatEvent, nil)
	s.RegisterEvent(GroupChatEvent, nil)
	s.RegisterEvent(BroadChatEvent, nil)
}

func (s *Manager) RegisterEvent(eventType EventType, handler EventHandler) {
	if !isValidEventType(eventType) {
		log.Fatalf("invalid event type %s", eventType)
	}

	s.handlerMux.Lock()
	defer s.handlerMux.Unlock()
	_, ok := s.handlers[eventType]
	if ok {
		log.Fatalf("event type %s has been exist", eventType)
	}

	s.handlers[eventType] = handler
}

func (s *Manager) AcceptConn(client *Client) {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	s.conns[client] = true
}

func (s *Manager) RemoveConn(client *Client) {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	delete(s.conns, client)
}

func (s *Manager) GetClients() map[*Client]bool {
	clientCopy := make(map[*Client]bool)
	s.clientMux.RLock()
	for c := range s.conns {
		clientCopy[c] = true
	}
	s.clientMux.RUnlock()
	return clientCopy
}

func (s *Manager) Do() {
	for {
		select {
		case event := <-s.Event:
			fmt.Println(event)
		}
	}
}

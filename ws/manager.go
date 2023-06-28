package ws

import (
	"context"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"nhooyr.io/websocket"
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
	s.RegisterEvent(SingleChatEvent, s.singleChatEvent)
	s.RegisterEvent(ChangeRoomEvent, s.changeRoomEvent)
	s.RegisterEvent(GroupChatEvent, s.groupChatEvent)

	s.RegisterEvent(BroadChatEvent, s.broadcastChatEvent)
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

func (s *Manager) Start() {
	log.Info("start ws manager")
	defer func() {
		if err := recover(); err != nil {
			log.Info(err)
		}
	}()

	for {
		select {
		case event := <-s.Event:
			handler, ok := s.handlers[event.Type]
			if !ok {
				log.Infof("manager received invalid event type[%s]", event.Type)
				continue
			}
			handler(event.Client, event.Data, event.filter)
		}
	}
}

func (s *Manager) singleChatEvent(client *Client, message json.RawMessage, filter FilterFunc) {
	d := struct {
		UserID uint64
		Msg    string
	}{}

	err := json.Unmarshal(message, &d)
	if err != nil {
		client.Write(context.Background(), websocket.MessageText, []byte(err.Error()))
		return
	}

	for c := range s.conns {
		if c.UserID == d.UserID {
			err = c.Write(context.Background(), websocket.MessageText, []byte(d.Msg))
			if err != nil {
				log.Printf("singleChatEvent write err %v", err)
				continue
			}
			log.Info("singleChatEvent send msg", d.Msg)
		}
	}
}

func (s *Manager) groupChatEvent(client *Client, message json.RawMessage, filter FilterFunc) {
	d := struct {
		RoomID uint64
		Msg    string
	}{}

	err := json.Unmarshal(message, &d)
	if err != nil {
		client.Write(context.Background(), websocket.MessageText, []byte(err.Error()))
		return
	}

	for c := range s.conns {
		if c.RoomID != d.RoomID {
			continue
		}

		if filter != nil {
			if filter(c) {
				continue
			} else {
				err = c.Write(context.Background(), websocket.MessageText, []byte(d.Msg))
				if err != nil {
					log.Printf("groupChatEvent write err %v", err)
					continue
				}
			}
		} else {
			err = c.Write(context.Background(), websocket.MessageText, []byte(d.Msg))
			if err != nil {
				log.Printf("groupChatEvent write err %v", err)
				continue
			}
		}
	}

	log.Infof("groupChatEvent room_id %d send msg %s", d.RoomID, d.Msg)
}

func (s *Manager) changeRoomEvent(client *Client, message json.RawMessage, filter FilterFunc) {
	d := struct {
		RoomID uint64
	}{}

	err := json.Unmarshal(message, &d)
	if err != nil {
		client.Write(context.Background(), websocket.MessageText, []byte(err.Error()))
		return
	}

	client.RoomID = d.RoomID
	log.Infof("client %p user_id %d room_id %d", client, client.UserID, client.RoomID)
}

func (s *Manager) broadcastChatEvent(client *Client, message json.RawMessage, filter FilterFunc) {
	d := struct {
		Msg string
	}{}

	err := json.Unmarshal(message, &d)
	if err != nil {
		client.Write(context.Background(), websocket.MessageText, []byte(err.Error()))
		return
	}

	for c := range s.conns {
		if filter != nil {
			if filter(c) {
				continue
			} else {
				err = c.Write(context.Background(), websocket.MessageText, []byte(d.Msg))
				if err != nil {
					log.Printf("broadcastChatEvent write err %v", err)
					continue
				}
			}
		} else {
			err = c.Write(context.Background(), websocket.MessageText, []byte(d.Msg))
			if err != nil {
				log.Printf("broadcastChatEvent write err %v", err)
				continue
			}
		}
	}

	log.Infof("broadcastChatEvent send msg %s", d.Msg)
}

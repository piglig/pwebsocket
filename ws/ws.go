package ws

import (
	"log"
	"math/rand"
	"net/http"
	"nhooyr.io/websocket"
)

type WSHandler struct {
	manager *Manager
}

func NewWsHandler() *WSHandler {
	return &WSHandler{manager: NewManager()}
}

func (h *WSHandler) GetManager() *Manager {
	return h.manager
}

func (h *WSHandler) Upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: false, OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Println(err)
		return
	}

	conn.SetReadLimit(1024)

	client := &Client{Conn: conn, UserID: rand.Uint64()}
	log.Println("accept new client", client)
	h.manager.AcceptConn(client)
	client.Heartbeat(h.manager)
	go client.Do(h.manager.Event)
}

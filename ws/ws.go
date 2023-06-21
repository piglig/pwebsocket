package ws

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"nhooyr.io/websocket"
)

type WSHandler struct {
	manager *Manager
}

func NewWsHandler() *WSHandler {
	return &WSHandler{manager: NewManager()}
}

func (h *WSHandler) Hello(c echo.Context) error {
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: false, OriginPatterns: []string{"*"}})
		if err != nil {
			c.Logger().Error(err)
			return
		}

		client := &Client{Conn: conn}
		h.manager.AcceptConn(client)
		client.Heartbeat(h.manager)

		for {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				return
			}
			if err != nil {
				c.Logger().Error(err)
				return
			}
		}

	})(c.Response(), c.Request())
	return nil
}

func (h *WSHandler) Start() {
	log.Println("start ws manager")
	h.manager.Do()
}

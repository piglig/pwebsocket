package ws

import (
	"context"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"nhooyr.io/websocket"
	"pwebsocket/ws"
	"time"
)

type WSHandler struct {
}

func NewWsHandler() *WSHandler {
	return &WSHandler{}
}

var (
	wsManager = ws.NewManager()
)

func (h *WSHandler) Hello(c echo.Context) error {
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: false, OriginPatterns: []string{"*"}})
		if err != nil {
			c.Logger().Error(err)
			return
		}

		client := &ws.Client{Conn: conn}
		wsManager.AcceptConn(client)
		client.Heartbeat(wsManager)

		for {
			err = echoHandler(client.Conn)
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

func echoHandler(conn *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	typ, r, err := conn.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := conn.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	err = w.Close()
	return err
}

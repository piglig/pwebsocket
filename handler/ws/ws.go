package ws

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"pwebsocket/ws"
)

type HelloWsHandler struct {
}

func NewHelloWsHandler() *HelloWsHandler {
	return &HelloWsHandler{}
}

var (
	wsManager = ws.NewManager()
)

func (h *HelloWsHandler) Hello(c echo.Context) error {
	websocket.Handler(func(conn *websocket.Conn) {
		wsManager.AcceptConn(conn)
		defer func() {
			conn.Close()
		}()
		for {
			err := websocket.Message.Send(conn, "hello")
			if err != nil {
				c.Logger().Error(err)
				return
			}

			msg := ""
			err = websocket.Message.Receive(conn, &msg)
			if err != nil {
				c.Logger().Error(err)
				return
			}
			c.Logger().Info("websocket hello", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

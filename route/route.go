package route

import (
	"github.com/labstack/echo/v4"
	"pwebsocket/ws"
)

func InitWs(e *echo.Echo) {
	helloHandler := ws.NewWsHandler()
	go helloHandler.Start()
	e.GET("ws", helloHandler.Hello)
}

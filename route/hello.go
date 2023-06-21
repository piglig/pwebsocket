package route

import (
	"github.com/labstack/echo/v4"
	"pwebsocket/ws"
)

func initHello(e *echo.Echo) {
	helloHandler := ws.NewWsHandler()
	go helloHandler.Start()
	e.GET("ws", helloHandler.Hello)
}

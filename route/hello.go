package route

import (
	"github.com/labstack/echo/v4"
	"pwebsocket/handler/ws"
)

func initHello(e *echo.Echo) {
	helloHandler := ws.NewWsHandler()
	e.GET("ws", helloHandler.Hello)
}

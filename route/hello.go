package route

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"pwebsocket/handler/ws"
)

func initHello(e *echo.Echo) {
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	helloHandler := ws.NewHelloWsHandler()
	e.GET("ws", helloHandler.Hello)
}

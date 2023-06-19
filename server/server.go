package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"pwebsocket/route"
)

type Server struct {
	addr string
	*echo.Echo
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
		Echo: echo.New(),
	}
}

func (s *Server) Listen() {
	route.InitRoute(s.Echo)

	err := s.Echo.Start(s.addr)
	if err != nil {
		if err != http.ErrServerClosed {
			s.Echo.Logger.Errorf("Listen", "err", err)
			return
		}
	}
}

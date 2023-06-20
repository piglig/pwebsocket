package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"pwebsocket/route"
	"syscall"
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
	s.Echo.Logger.SetLevel(log.DEBUG)

	route.InitRoute(s.Echo)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-sigs
		s.Logger.Info("received signal", sig.String())
		defer cancel()
		err := s.Echo.Shutdown(ctx)
		if err != nil {
			s.Logger.Fatal(err)
		}
	}()

	err := s.Echo.Start(s.addr)
	if err != nil {
		if err != http.ErrServerClosed {
			s.Echo.Logger.Errorf("Listen", "err", err)
			return
		}
	}
}

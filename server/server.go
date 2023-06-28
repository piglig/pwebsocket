package server

import (
	"context"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"pwebsocket/route"
	"syscall"
	"time"
)

type Server struct {
	addr string
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Listen() {
	srv := &http.Server{
		Addr:              s.addr,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           route.InitWs(),
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-sigs
		log.Info("received signal", sig.String())
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Infof("listen on %s", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Errorf("Listen", "err", err)
			return
		}
	}
}

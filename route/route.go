package route

import (
	"net/http"
	"pwebsocket/ws"
)

func InitWs() *http.ServeMux {
	handler := ws.NewWsHandler()
	manager := handler.GetManager()
	go manager.Start()

	mux := http.NewServeMux()
	mux.HandleFunc("ws", handler.Upgrade)

	return mux
}

package app

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/peefourtee/labrack"
)

func init() {
}

func Mux(t <-chan labrack.Telemetry) *http.ServeMux {
	mux := http.NewServeMux()

	wsEndpoint := &wsEndpoint{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		server: &wsServer{conns: make(map[string]*websocket.Conn)},
	}

	go wsEndpoint.server.HandleTelemetry(t)

	mux.Handle("/ws", wsEndpoint)
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	return mux
}

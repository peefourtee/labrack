package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/peefourtee/labrack"
)

var ErrClientExists = errors.New("client already exists")

type wsMessage struct {
	Type string
	Data interface{}
}

type wsServer struct {
	mu    sync.Mutex
	conns map[string]*websocket.Conn
}

func (s *wsServer) Add(c *websocket.Conn) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.conns[c.RemoteAddr().String()]; ok {
		return ErrClientExists
	}
	s.conns[c.RemoteAddr().String()] = c

	// drain all incoming messages from the connection
	go func() {
		for {
			if _, _, err := c.NextReader(); err != nil {
				log.Printf("error reading from ws connection %s: %s", c.RemoteAddr(), err)
				s.Remove(c)
				break
			}
		}
	}()
	return nil
}

func (s *wsServer) HandleTelemetry(t <-chan labrack.Telemetry) {
	for d := range t {
		log.Printf("received telemetry: %+v", d)
		s.WriteAll(wsMessage{Type: "telemetry", Data: d})
	}
}

func (s *wsServer) Remove(c *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.remove(c)
}
func (s *wsServer) remove(c *websocket.Conn) {
	delete(s.conns, c.RemoteAddr().String())
}

func (s wsServer) WriteAll(m wsMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	errs := make([]*websocket.Conn, 0)

	for _, c := range s.conns {
		if err := c.WriteJSON(m); err != nil {
			log.Printf("failed to write %q message to %s: %s", m.Type, c.RemoteAddr(), err)
			errs = append(errs, c)
		}
	}
	for _, c := range errs {
		s.remove(c)
	}
}

type wsEndpoint struct {
	upgrader websocket.Upgrader
	server   *wsServer
}

func (e *wsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := e.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error with connecting websocket: ", err)
		fmt.Fprintf(w, err.Error())
		return
	}

	log.Print("got new client: ", conn.RemoteAddr())
	if err := e.server.Add(conn); err != nil {
		log.Print("failed to add websocket conn: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		conn.Close()
	}
}

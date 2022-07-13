package webinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebInfoServer struct {
	Info    OutInfo
	Server  *http.Server
	Clients map[*Client]int
	lock    sync.Mutex
}

type Client struct {
	conn  *websocket.Conn
	Data  chan []byte
	Close chan byte
}

func NewWebInfoServer(port int) *WebInfoServer {
	server := &WebInfoServer{
		Clients: map[*Client]int{},
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./assets/webinfo")))
	mux.HandleFunc("/ws/info", server.handleInfo)
	mux.HandleFunc("/api/info", server.getInfo)
	server.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return server
}

func (s *WebInfoServer) getInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d, _ := json.Marshal(s.Info)
	_, err := w.Write(d)
	if err != nil {
		lg.Warnf("api get info error: %s", err)
		return
	}
}

func (s *WebInfoServer) handleInfo(w http.ResponseWriter, r *http.Request) {
	lg.Debug("connection start")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		lg.Warnf("upgrade error: %s", err)
		return
	}
	client := &Client{
		conn:  conn,
		Data:  make(chan []byte, 16),
		Close: make(chan byte, 1),
	}
	s.addClient(client)
	defer s.removeClient(client)
	go func() {
		for {
			_, _, err := client.conn.ReadMessage()
			if err != nil {
				client.Close <- 1
			}
		}
	}()
	for {
		lg.Trace("waiting for message")
		select {
		case data := <-client.Data:
			writer, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				lg.Warn("get writer error", err)
				return
			}

			if _, err = writer.Write(data); err != nil {
				lg.Warn("send error:", err)
				return
			}
			if err = writer.Close(); err != nil {
				lg.Warnf("can't close writer: %s", err)
				return
			}
		case _ = <-client.Close:
			lg.Debug("client close")
			if err := client.conn.Close(); err != nil {
				lg.Warnf("close connection encouter an error: %s", err)
			}
			return
		}
	}
}

func (s *WebInfoServer) SendInfo(update string, info OutInfo) {
	for client := range s.Clients {
		d, _ := json.Marshal(WebsocketData{Update: update, Data: info})
		client.Data <- d
	}
}

func (s *WebInfoServer) addClient(c *Client) {
	s.lock.Lock()
	s.Clients[c] = 1
	s.lock.Unlock()
}

func (s *WebInfoServer) removeClient(c *Client) {
	s.lock.Lock()
	close(c.Data)
	delete(s.Clients, c)
	s.lock.Unlock()
}

func (s *WebInfoServer) Start() {
	go func() {
		err := s.Server.ListenAndServe()
		if err == http.ErrServerClosed {
			lg.Info("server closed")
			return
		}
		if err != nil {
			lg.Warnf("Failed to start webinfo server: %s", err)
			return
		}
	}()
}

func (s *WebInfoServer) Stop() error {
	return s.Server.Shutdown(context.Background())
}

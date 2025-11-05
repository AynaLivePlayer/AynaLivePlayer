package wshub

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/logger"
	"context"
	"encoding/json"
	"errors"
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

type wsClient struct {
	conn  *websocket.Conn
	Data  chan []byte
	Close chan byte
}

func (c *wsClient) start() {
	for {
		msgType, val, err := c.conn.ReadMessage()
		if err != nil {
			c.Close <- 1
			return
		}
		if msgType != websocket.TextMessage {
			return
		}
		var data EventDataReceived
		err = json.Unmarshal(val, &data)
		if err != nil {
			global.Logger.Warn("unmarshal event data failed", err)
			return
		}
		actualEventData, err := events.UnmarshalEventData(data.EventID, data.Data)
		if err != nil {
			global.Logger.Warn("unmarshal event data failed", err)
			return
		}
		if globalEnableWsHubControl {
			_ = global.EventBus.PublishToChannel(eventChannel, data.EventID, actualEventData)
		}
	}
}

type wsServer struct {
	Running       bool
	Server        *http.Server
	clients       map[*wsClient]bool
	mux           *http.ServeMux
	lock          sync.RWMutex
	port          *int
	localhostOnly *bool
	log           logger.ILogger
}

func newWsServer(port *int, localhostOnly *bool) *wsServer {
	mux := http.NewServeMux()
	s := &wsServer{
		Running:       false,
		clients:       make(map[*wsClient]bool),
		mux:           mux,
		port:          port,
		localhostOnly: localhostOnly,
		log:           global.Logger.WithPrefix("plugin.wshub.server"),
	}
	mux.HandleFunc("/wsinfo", s.handleWsInfo)
	return s
}

func (s *wsServer) broadcast(data []byte) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for client := range s.clients {
		client.Data <- data
	}
}

func (s *wsServer) register(client *wsClient) {
	s.lock.Lock()
	s.clients[client] = true
	s.lock.Unlock()
}

func (s *wsServer) unregister(client *wsClient) {
	s.lock.Lock()
	delete(s.clients, client)
	s.lock.Unlock()
}

func (s *wsServer) handleWsInfo(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("connection start")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Warnf("upgrade error: %s", err)
		return
	}
	client := &wsClient{
		conn:  conn,
		Data:  make(chan []byte, 16),
		Close: make(chan byte, 1),
	}
	s.register(client)
	defer s.unregister(client)
	go client.start()
	// send initial data
	for _, data := range eventCache {
		// ignore empty
		if data.EventID == "" {
			continue
		}
		eventCacheData, _ := toCapitalizedJSON(data)
		err := client.conn.WriteMessage(websocket.TextMessage, eventCacheData)
		if err != nil {
			s.log.Warn("write message failed", err)
			return
		}
	}
	// start message loop
	for {
		select {
		case data := <-client.Data:
			err := client.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				s.log.Warn("write message failed", err)
				return
			}
		case _ = <-client.Close:
			s.log.Infof("client %s close", client.conn.RemoteAddr().String())
			if err := client.conn.Close(); err != nil {
				s.log.Warnf("close connection encouter an error: %s", err)
			}
			return
		}
	}
}

func (s *wsServer) Start() {
	s.log.Debug("WebInfoServer starting...")
	s.Running = true
	go func() {
		var addr string
		if *s.localhostOnly {
			addr = fmt.Sprintf("localhost:%d", *s.port)
		} else {
			addr = fmt.Sprintf("0.0.0.0:%d", *s.port)
		}
		s.Server = &http.Server{
			Addr:    addr,
			Handler: s.mux,
		}
		err := s.Server.ListenAndServe()
		s.Running = false
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("WebInfoServer closed")
			return
		}
		if err != nil {
			s.log.Errorf("Failed to start webinfo server: %s", err)
			return
		}
	}()
}

func (s *wsServer) Stop() error {
	s.log.Debug("WebInfoServer stopping...")
	s.lock.Lock()
	s.clients = make(map[*wsClient]bool)
	s.lock.Unlock()
	if s.Server != nil {
		return s.Server.Shutdown(context.TODO())
	}
	return nil
}

func (s *wsServer) getWsUrl() string {
	if *s.localhostOnly {
		return fmt.Sprintf("ws://localhost:%d/wsinfo", *s.port)
	}
	return fmt.Sprintf("ws://0.0.0.0:%d/wsinfo", *s.port)
}

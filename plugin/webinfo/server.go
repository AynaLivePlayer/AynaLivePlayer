package webinfo

import (
	"AynaLivePlayer/common/config"
	"AynaLivePlayer/core/adapter"
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
	Info      OutInfo
	Port      int
	ServerMux *http.ServeMux
	Server    *http.Server
	Clients   map[*Client]int
	Running   bool
	Store     *TemplateStore
	lock      sync.Mutex
	log       adapter.ILogger
}

type Client struct {
	conn  *websocket.Conn
	Data  chan []byte
	Close chan byte
}

func NewWebInfoServer(port int, log adapter.ILogger) *WebInfoServer {
	server := &WebInfoServer{
		Store:   newTemplateStore(WebTemplateStorePath),
		Port:    port,
		Info:    OutInfo{Playlist: make([]MediaInfo, 0)},
		Clients: map[*Client]int{},
		log:     log,
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(config.GetAssetPath("webinfo"))))
	mux.HandleFunc("/ws/info", server.handleInfo)
	mux.HandleFunc("/api/info", server.getInfo)
	mux.HandleFunc("/api/template/list", server.tmplList)
	mux.HandleFunc("/api/template/get", server.tmplGet)
	mux.HandleFunc("/api/template/save", server.tmplSave)
	server.ServerMux = mux

	return server
}

func (s *WebInfoServer) tmplList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	d, _ := json.Marshal(s.Store.List())
	_, err := w.Write(d)
	if err != nil {
		s.log.Warnf("/api/template/list error: %s", err)
		return
	}
}

func (s *WebInfoServer) tmplGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "default"
	}
	d, _ := json.Marshal(s.Store.Get(name))
	_, err := w.Write(d)
	if err != nil {
		s.log.Warnf("/api/template/get error: %s", err)
		return
	}
}

func (s *WebInfoServer) tmplSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	s.log.Info(r.Method)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(1 << 16); err != nil {
		s.log.Warnf("ParseForm() err: %v", err)
		return
	}
	name := r.FormValue("name")
	tmpl := r.FormValue("template")
	if name == "" {
		name = "default"
	}
	s.log.Infof("change template %s", name)
	s.Store.Modify(name, tmpl)
	d, _ := json.Marshal(s.Store.Get(name))
	_, err := w.Write(d)
	if err != nil {
		s.log.Warnf("/api/template/save error: %s", err)
		return
	}
	s.Store.Save(WebTemplateStorePath)
}

func (s *WebInfoServer) getInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	d, _ := json.Marshal(s.Info)
	_, err := w.Write(d)
	if err != nil {
		s.log.Warnf("api get info error: %s", err)
		return
	}
}

func (s *WebInfoServer) handleInfo(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("connection start")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Warnf("upgrade error: %s", err)
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
		s.log.Debug("waiting for message")
		select {
		case data := <-client.Data:
			writer, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				s.log.Warn("get writer error", err)
				return
			}

			if _, err = writer.Write(data); err != nil {
				s.log.Warn("send error:", err)
				return
			}
			if err = writer.Close(); err != nil {
				s.log.Warnf("can't close writer: %s", err)
				return
			}
		case _ = <-client.Close:
			s.log.Debug("client close")
			if err := client.conn.Close(); err != nil {
				s.log.Warnf("close connection encouter an error: %s", err)
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
	s.log.Debug("WebInfoServer starting...")
	s.Running = true
	go func() {
		s.Server = &http.Server{
			Addr:    fmt.Sprintf("localhost:%d", s.Port),
			Handler: s.ServerMux,
		}
		err := s.Server.ListenAndServe()
		s.Running = false
		if err == http.ErrServerClosed {
			s.log.Info("WebInfoServer closed")
			return
		}
		if err != nil {
			s.log.Warnf("Failed to start webinfo server: %s", err)
			return
		}
	}()
}

func (s *WebInfoServer) Stop() error {
	s.log.Debug("WebInfoServer stopping...")
	s.lock.Lock()
	s.Clients = map[*Client]int{}
	s.lock.Unlock()
	if s.Server != nil {
		return s.Server.Shutdown(context.TODO())
	}
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebInfo struct {
	A string
	B string
}

type WebInfoServer struct {
	ServeMux http.ServeMux
	Clients  map[*Client]int
	lock     sync.Mutex
}

type Client struct {
	conn  *websocket.Conn
	Data  chan []byte
	Close chan byte
}

func NewWebInfoServer() *WebInfoServer {
	server := &WebInfoServer{
		Clients: map[*Client]int{},
	}
	server.ServeMux.Handle("/", http.FileServer(http.Dir("./assets/webinfo")))
	server.ServeMux.HandleFunc("/ws/info", server.handleInfo)
	return server
}

func (s *WebInfoServer) handleInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connection start")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade error", err)
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
		fmt.Println("waiting for message")
		select {
		case data := <-client.Data:
			writer, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println("get writer error", err)
				return
			}

			if _, err = writer.Write(data); err != nil {
				fmt.Println("send error:", err)
				return
			}
			if err = writer.Close(); err != nil {
				fmt.Println("can't close writer")
				return
			}
		case _ = <-client.Close:
			fmt.Println("client close", client.conn.Close())
			return
		}
	}
}

func (s *WebInfoServer) sendInfo(info *WebInfo) {
	for client := range s.Clients {
		d, _ := json.Marshal(info)
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

var info WebInfo = WebInfo{A: "asdf", B: "ffff"}

func main() {
	server := NewWebInfoServer()
	go func() {
		for {
			time.Sleep(time.Second * 5)
			server.sendInfo(&info)
		}
	}()
	http.ListenAndServe("localhost:8080", &server.ServeMux)
}

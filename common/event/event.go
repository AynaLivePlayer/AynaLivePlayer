package event

import (
	"sync"
)

type EventId string

type Event struct {
	Id        EventId
	Cancelled bool
	Data      interface{}
}

type HandlerFunc func(event *Event)

type Handler struct {
	EventId EventId
	Name    string
	Handler HandlerFunc
}

type Manager struct {
	handlers   map[string]*Handler
	queue      chan func()
	stopSig    chan int
	queueSize  int
	workerSize int
	lock       sync.RWMutex
}

func NewManger(queueSize int, workerSize int) *Manager {
	manager := &Manager{
		handlers:   make(map[string]*Handler),
		queue:      make(chan func(), queueSize),
		stopSig:    make(chan int, workerSize),
		queueSize:  queueSize,
		workerSize: workerSize,
	}
	for i := 0; i < workerSize; i++ {
		go func() {
			for {
				select {
				case <-manager.stopSig:
					return
				case f := <-manager.queue:
					f()
				}
			}
		}()
	}
	return manager
}

func (h *Manager) NewChildManager() *Manager {
	return &Manager{
		handlers:   make(map[string]*Handler),
		queue:      h.queue,
		stopSig:    h.stopSig,
		queueSize:  h.queueSize,
		workerSize: h.workerSize,
	}
}

func (h *Manager) Stop() {
	for i := 0; i < h.workerSize; i++ {
		h.stopSig <- 0
	}
}

func (h *Manager) Register(handler *Handler) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.handlers[handler.Name] = handler
}

func (h *Manager) RegisterA(id EventId, name string, handler HandlerFunc) {
	h.Register(&Handler{
		EventId: id,
		Name:    name,
		Handler: handler,
	})
}

func (h *Manager) UnregisterAll() {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.handlers = make(map[string]*Handler)
}

func (h *Manager) Unregister(name string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.handlers, name)
}

func (h *Manager) Call(event *Event) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	for _, eh := range h.handlers {
		if eh.EventId == event.Id {
			handler := eh.Handler
			h.queue <- func() {
				handler(event)
			}
		}
	}
}

func (h *Manager) CallA(id EventId, data interface{}) {
	h.Call(&Event{
		Id:   id,
		Data: data,
	})
}
